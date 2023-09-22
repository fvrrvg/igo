package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Davincible/goinsta/v3"
	"github.com/briandowns/spinner"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

var (
	Writer = uilive.New()
	Red    = color.New(color.FgRed).SprintFunc()
	Green  = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
	Purple = color.New(color.FgHiMagenta).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
	HiBlue = color.New(color.FgHiBlue).SprintFunc()
)

func GetCurrDir() string {
	currdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return currdir + "/"
}

func Dir(u *goinsta.User) string {
	dir := GetCurrDir() + u.Username + "/"

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 0o777)
		if err != nil {
			panic(err)
		}
	}

	return dir
}

func GetUser(insta *goinsta.Instagram, username string) *goinsta.User {
	user, err := insta.Profiles.ByName(username)
	if err != nil {
		fmt.Println(Yellow("User not found, Search by name: "))
		var query string
		Scanner := bufio.NewScanner(os.Stdin)
		Scanner.Scan()
		query = Scanner.Text()

		u := SearchMenu(query, insta)
		return u
	}

	return user
}

func SearchMenu(query string, insta *goinsta.Instagram) *goinsta.User {
	fmt.Print("\033[H\033[2J")

	search, err := insta.Searchbar.SearchUser(query, true)
	if err != nil {
		log.Fatal(err)
	}
	users := search.Users

	currentIndex := 0
	pageSize := 10
	quit := false

	err = keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	for !quit {
		fmt.Print("\033[H\033[2J")

		if currentIndex < 0 {
			currentIndex = 0
		}
		if currentIndex >= len(users) {
			fmt.Println("Not Found!")
			break
		}

		startIndex := currentIndex
		endIndex := startIndex + pageSize
		if endIndex > len(users) {
			endIndex = len(users)
		}
		currentPageUsers := users[startIndex:endIndex]

		fmt.Print(Blue("Select: " + (strconv.Itoa(startIndex + 1)) + "/" + (strconv.Itoa(len(users)) + "\n")))

		for i, u := range currentPageUsers {
			arrow := " "
			if i == currentIndex-startIndex {
				arrow = ">"
			}
			fmt.Printf("%s %d. %v\n", arrow, startIndex+i+1, u.Username)
		}
		fmt.Println()
		fmt.Print(Green("Navigation: Up/Down Arrow Keys\n"))
		fmt.Print(Green("Enter: to Select a user\n"))
		fmt.Println(Green("Q: Quit\n"))
		fmt.Println()

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyArrowUp && currentIndex > 0 {
			currentIndex--
		} else if key == keyboard.KeyArrowDown && currentIndex < len(users)-1 {
			currentIndex++
		} else if key == keyboard.KeyEnter {
			u, _ := insta.Profiles.ByName(users[currentIndex].Username)
			return u
		} else if char == 'q' || char == 'Q' {
			quit = true
			os.Exit(0)

		}

	}

	os.Exit(0)
	return nil
}

func GetSpinner() *spinner.Spinner {
	S := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	return S
}
