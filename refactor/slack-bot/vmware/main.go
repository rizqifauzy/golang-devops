package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func main() {
	// Memuat file .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s\n", err.Error())
		return
	}

	// Membaca username, password, dan hostname dari .env
	hostname := os.Getenv("VCENTER_HOSTNAME")
	username := os.Getenv("VCENTER_USERNAME")
	password := os.Getenv("VCENTER_PASSWORD")
	if hostname == "" || username == "" || password == "" {
		fmt.Println("Error: VCENTER_HOSTNAME, VCENTER_USERNAME, or VCENTER_PASSWORD is not set in .env")
		return
	}

	// Encode username dan password
	encodedUsername := url.QueryEscape(username)
	encodedPassword := url.QueryEscape(password)

	// Membuat URL untuk koneksi ke vCenter
	fmt.Println("Creating a VIM/SOAP session.")
	vcURL := "https://" + encodedUsername + ":" + encodedPassword + "@" + hostname + "/sdk"
	u, err := url.Parse(vcURL)
	if err != nil {
		fmt.Printf("Error parsing url %s\n", vcURL)
		return
	}

	// Membuat context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Login ke vCenter
	vimClient, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Printf("Error logging in: %s\n", err.Error())
		return
	}
	defer vimClient.Logout(ctx)

	fmt.Println("Login successful!")

	// List Virtual Machines
	err = listVMs(ctx, vimClient)
	if err != nil {
		fmt.Printf("Error listing VMs: %s\n", err.Error())
		return
	}

	// Stop VM
	vmNameToStop := "servertesting2" // Ganti dengan nama VM yang ingin dihentikan
	err = stopVM(ctx, vimClient, vmNameToStop)
	if err != nil {
		fmt.Printf("Error stopping VM: %s\n", err.Error())
		return
	}

	// Start VM
	vmNameToStart := "servertesting2" // Ganti dengan nama VM yang ingin dijalankan
	err = startVM(ctx, vimClient, vmNameToStart)
	if err != nil {
		fmt.Printf("Error starting VM: %s\n", err.Error())
		return
	}
}

// listVMs mencantumkan semua VM yang ada di vCenter
func listVMs(ctx context.Context, client *govmomi.Client) error {
	// Membuat view manager untuk objek VirtualMachine
	m := view.NewManager(client.Client)

	// Membuat container view untuk VM
	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return fmt.Errorf("failed to create container view: %w", err)
	}
	defer v.Destroy(ctx)

	// Mengambil daftar objek VirtualMachine
	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"name"}, &vms)
	if err != nil {
		return fmt.Errorf("failed to retrieve VMs: %w", err)
	}

	// Menampilkan nama VM
	fmt.Println("Virtual Machines:")
	for _, vm := range vms {
		fmt.Println("- " + vm.Name)
	}

	return nil
}

// stopVM menghentikan VM dengan nama yang diberikan
func stopVM(ctx context.Context, client *govmomi.Client, vmName string) error {
	// Cari VM berdasarkan nama
	vm, err := findVM(ctx, client, vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM: %w", err)
	}

	// Hentikan VM
	task, err := vm.PowerOff(ctx)
	if err != nil {
		return fmt.Errorf("failed to power off VM: %w", err)
	}

	// Tunggu task selesai
	err = task.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for power off task: %w", err)
	}

	fmt.Printf("VM %s stopped successfully.\n", vmName)
	return nil
}

// startVM menjalankan VM dengan nama yang diberikan
func startVM(ctx context.Context, client *govmomi.Client, vmName string) error {
	// Cari VM berdasarkan nama
	vm, err := findVM(ctx, client, vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM: %w", err)
	}

	// Jalankan VM
	task, err := vm.PowerOn(ctx)
	if err != nil {
		return fmt.Errorf("failed to power on VM: %w", err)
	}

	// Tunggu task selesai
	err = task.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for power on task: %w", err)
	}

	fmt.Printf("VM %s started successfully.\n", vmName)
	return nil
}

// findVM mencari VM berdasarkan nama
func findVM(ctx context.Context, client *govmomi.Client, vmName string) (*object.VirtualMachine, error) {
	// Membuat view manager untuk objek VirtualMachine
	m := view.NewManager(client.Client)

	// Membuat container view untuk VM
	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create container view: %w", err)
	}
	defer v.Destroy(ctx)

	// Mengambil daftar objek VirtualMachine
	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"name"}, &vms)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve VMs: %w", err)
	}

	// Cari VM berdasarkan nama
	for _, vm := range vms {
		if vm.Name == vmName {
			return object.NewVirtualMachine(client.Client, vm.Reference()), nil
		}
	}

	return nil, fmt.Errorf("VM %s not found", vmName)
}
