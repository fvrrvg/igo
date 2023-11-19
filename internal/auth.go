package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fvrrvg/igo/utils"

	"github.com/Davincible/goinsta/v3"
	"github.com/briandowns/spinner"
	"golang.org/x/term"
)

func GetCreds() (string, string) {
	var username, password string

	fmt.Println("You have to login first!")
	fmt.Print("Enter your username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter your password: ")

	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	for {
		char := make([]byte, 1)
		_, err := os.Stdin.Read(char)
		if err != nil {
			fmt.Println("Error:", err)
		}

		// Enter key (13) terminates the input
		if char[0] == 13 {
			break
		}

		// Backspace (127) handling
		if char[0] == 127 && len(password) > 0 {
			password = password[:len(password)-1]
			fmt.Print("\b \b") // Clear the last character
		} else {
			password += string(char)
			fmt.Print("*")
		}
	}

	fmt.Println()
	fmt.Println()

	username = strings.ToLower(username)
	return username, password
}

func ConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(utils.Red("Could not get user home directory!"))
	}

	if os.PathSeparator == '\\' {
		return filepath.Join(home, "AppData", "Local", "igo")
	} else {
		return filepath.Join(home, ".config", "igo")
	}
}

func Login() *goinsta.Instagram {
	CF := ConfigFile()
	if _, err := os.Stat(CF); err == nil {
		insta, err := goinsta.Import(CF)
		if err != nil {
			username, password := GetCreds()
			insta := goinsta.New(username, password)
			loginWithSpinner(insta)
			defer insta.Export(CF)
			return insta
		}
		return insta
	} else {
		username, password := GetCreds()
		insta := goinsta.New(username, password)
		loginWithSpinner(insta)
		defer insta.Export(CF)
		return insta
	}
}

func loginWithSpinner(insta *goinsta.Instagram) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.Color("yellow")
	s.Start()
	err := insta.Login()

	if err == goinsta.Err2FARequired {
		s.Stop()
		fmt.Println("2FA is required!")
		fmt.Print("Enter the code: ")
		var code string
		fmt.Scanln(&code)
		s.Start()
		s.Color("yellow")
		err := insta.TwoFactorInfo.Login2FA(code)
		if err != nil {
			s.Stop()
			fmt.Println(utils.Red(err.Error()))
		}
	} else if err != nil {
		s.Stop()
		fmt.Println(utils.Red(err.Error()))
	} else if err == goinsta.ErrChallengeRequired {
		s.Stop()
		fmt.Println(utils.Red("Challenge Required, please open Instagram app to allow the login request or trust it"))
	} else {
		s.Stop()
		fmt.Println(utils.Green("Logged in successfully!"))
		fmt.Println()
	}
}

func Logout() {
	CF := ConfigFile()
	fmt.Println("Logging out...")
	if _, err := os.Stat(CF); err == nil {
		err := os.Remove(CF)
		if err != nil {
			panic(err)
		}
	}
	fmt.Print(utils.Green("Logged out successfully!"))
}
