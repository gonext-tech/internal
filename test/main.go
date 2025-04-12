package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$08$WTZTs7HrNAkMs2URrogzoOvDcP.IvpPow6Ru6Zu2.4JzGrRCK8rGS"
	password := "123456"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("❌ Password mismatch:", err)
	} else {
		fmt.Println("✅ Password match!")
	}
}
