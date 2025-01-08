package ssh

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func StartServer(privateKey []byte, authorizedKeys []byte) error {
	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeys) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeys)
		if err != nil {
			return fmt.Errorf("Parse Authorized keys error: %s", err)
		}

		authorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKeys = rest
	}

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			if authorizedKeysMap[string(pubKey.Marshal())] {
				return &ssh.Permissions{
					// Record the public key used for authentication.
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}

	private, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("ParsePrivateKey error: %s", err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		return fmt.Errorf("Listen error: %s", err)
	}

	for {
		nConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Listener accept error: %s\n", err)
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		conn, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			fmt.Printf("NewServerConn error: %s\n", err)
		}
		if conn != nil && conn.Permissions != nil {
			log.Printf("logged in with key %s", conn.Permissions.Extensions["pubkey-fp"])
		}

		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		go handleConnection(conn, chans)
	}
}

func handleConnection(conn *ssh.ServerConn, chans <-chan ssh.NewChannel) {
	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			fmt.Printf("Could not accept channel: %v\n", err)
			continue // not shown in demo lecture, but we can skip this loop iteration when there's an error
		}

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "shell" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Printf("Request Type made by client: %s\n", req.Type)
				switch req.Type {
				case "exec":
					payload := req.Payload
					if len(payload) < 4 {
						fmt.Printf("Payload too short: %v\n", payload)
						channel.Write([]byte("Payload too short\n"))
						channel.SendRequest("exit-status", false, []byte{0, 0, 0, 1})
						req.Reply(false, nil)
						channel.Close()
						continue
					}
					commandLength := int(payload[0])<<24 | int(payload[1])<<16 | int(payload[2])<<8 | int(payload[3])
					if len(payload) < 4+commandLength {
						fmt.Printf("Payload length mismatch: %v\n", payload)
						channel.Write([]byte("Payload length mismatch\n"))
						channel.SendRequest("exit-status", false, []byte{0, 0, 0, 1})
						req.Reply(false, nil)
						channel.Close()
						continue
					}
					command := string(bytes.TrimSpace(payload[4 : 4+commandLength]))
					fmt.Printf("Received payload: %s\n", command)
					output := execSomething(conn, command)                        // Execute the command
					channel.Write([]byte(output))                                 // Write output to the channel
					channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0}) // Send exit status
					req.Reply(true, nil)                                          // Acknowledge the request
					channel.Close()                                               // Close the channel after execution                                            // Close the channel after execution
				case "shell":
					req.Reply(true, nil)
				case "pty-req":
					createTerminal(conn, channel)
				default:
					req.Reply(false, nil)
				}

			}
		}(requests)
	}
}

func createTerminal(conn *ssh.ServerConn, channel ssh.Channel) {
	termInstance := term.NewTerminal(channel, "> ")
	go func() {
		defer channel.Close()
		for {
			line, err := termInstance.ReadLine()
			if err != nil {
				fmt.Printf("ReadLine error: %s", err)
				break
			}
			switch line {
			// case "whoami":
			// 	termInstance.Write([]byte(execSomething(conn, []byte("whoami"))))
			case "":
			case "quit":
				termInstance.Write([]byte("Goodbye!\n"))
				channel.Close()
			default:
				output := execSomething(conn, line)
				termInstance.Write([]byte(output + "\n"))
			}
		}
	}()
}

// func execSomething(conn *ssh.ServerConn, payload []byte) string {
// 	switch string(payload) {
// 	case "whoami":
// 		return fmt.Sprintf("You are: %s\n", conn.Conn.User())
// 	default:
// 		return fmt.Sprintf("Command Not Found: %s\n", string(payload))
// 	}
// }

// func execSomething(conn *ssh.ServerConn, payload []byte) string {
// 	fmt.Printf("Received payload: %s\n", string(payload))
// 	// Trim payload to remove SSH protocol-specific prefix
// 	command := string(bytes.TrimSpace(payload)) // Convert payload to string and trim spaces
// 	fmt.Printf("Executing command: %s\n", command)

// 	// Use sh -c to allow shell commands execution
// 	cmd := exec.Command("/bin/bash", "-c", command)
// 	cmd.Env = append(cmd.Env, "PATH=/bin:/usr/bin:/usr/local/bin")
// 	// Debug PATH
// 	fmt.Printf("Using PATH: %s\n", cmd.Env)

// 	// Capture combined output (stdout and stderr)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("Command failed: %s\n", err)
// 		return fmt.Sprintf("Error executing command: %s\nOutput: %s\n", err, string(output))
// 	}
// 	return string(output)
// }

func execSomething(conn *ssh.ServerConn, command string) string {
	fmt.Printf("Executing command: %s\n", command)

	// Split the command into arguments
	args := strings.Fields(command)
	if len(args) == 0 {
		return "No command provided\n"
	}

	// Use the first argument as the command and the rest as arguments
	cmd := exec.Command(args[0], args[1:]...)
	//cmd.Env = append(cmd.Env, "PATH=/bin:/usr/bin:/usr/local/bin")
	// Debug PATH
	//fmt.Printf("Using PATH: %s\n", cmd.Env)

	// Capture combined output (stdout and stderr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed: %s\n", err)
		return fmt.Sprintf("Error executing command: %s\nOutput: %s\n", err, string(output))
	}
	return string(output)
}
