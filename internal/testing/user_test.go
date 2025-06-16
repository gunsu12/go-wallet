package testing

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUserHasingPassword(t *testing.T) {
	// Simulasi password dari user input (misal saat login)
	inputPassword := "royal2k19!"

	// Simulasi password hash dari database
	storedHash := "$2y$10$2hBIDjjNGQgosGCiqI4GjeLys9UDemKOlY86.V5RFOxoK3yCG8TjW"

	// Bandingkan hash dengan password input
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPassword))
	if err != nil {
		fmt.Println("Password TIDAK cocok:", err)
	} else {
		fmt.Println("Password cocok!")
	}
}
