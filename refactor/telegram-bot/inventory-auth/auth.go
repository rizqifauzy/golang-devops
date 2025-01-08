package main

func isUserAuthorized(userID int64) bool {
	// Daftar ID pengguna yang diizinkan
	allowedUsers := []int64{
		243211339, // Ganti dengan ID pengguna yang diizinkan
	}

	for _, id := range allowedUsers {
		if userID == id {
			return true
		}
	}
	return false
}
