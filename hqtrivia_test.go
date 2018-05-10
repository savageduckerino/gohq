package go_hqtrivia

import (
	"testing"
	"fmt"
	"log"
)

func TestTrivia(t *testing.T) {
	ver, err := HQVerify("+12134191735")

	if err != nil {
		log.Fatal("Error verifying number:", err)
	}

	var code = "8981"
	/*fmt.Print("Enter the 4 digit code: ")
	fmt.Scanln(&code)*/

	auth, err := HQConfirm(ver, code)
	if err != nil {
		log.Fatal("Failed to confirm:", err)
	} else {
		if auth == nil {
			fmt.Println("This is an unregistered account.")
		} else {
			log.Fatal("This is an existing account:", auth.Auth.Username)
		}
	}

	info, err := HQCreate(ver, "randomuseme7191", "Discoli", "GB")
	if err != nil {
		log.Fatal("Failed to create account:", err)
	}

	fmt.Println("Signed up as ", info.Username, "with the id ", info.UserID)
	fmt.Print("Applying free weekly life!")
	if err = HQWeekly(info); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print("Done!")
	}
}
