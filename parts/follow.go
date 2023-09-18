package parts

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"igo/utils"

	"github.com/Davincible/goinsta/v3"
	"github.com/briandowns/spinner"
	"github.com/eiannone/keyboard"
)

func Follow(goinsta *goinsta.Instagram, username string) {
	user, err := getUserByUsername(goinsta, username)
	if err != nil {
		fmt.Println(utils.Red("User not found."))
		os.Exit(1)
	}

	user.Follow()
	fmt.Println(utils.Green("Followed successfully!"))
}

func Unfollow(goinsta *goinsta.Instagram, username string) {
	user, err := getUserByUsername(goinsta, username)
	if err != nil {
		fmt.Println(utils.Red("User not found."))
		os.Exit(1)
	}

	user.Unfollow()
	fmt.Println(utils.Green("Unfollowed " + username + " successfully!"))
}

func getUserByUsername(goinsta *goinsta.Instagram, username string) (*goinsta.User, error) {
	user, err := goinsta.Profiles.ByName(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NotFollowingBack(goinsta *goinsta.Instagram, json string, excludeverified string) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.Color("yellow")
	s.Start()
	followers := getFollowers(goinsta)
	notFollowingBack := []string{}
	counter := 1
	var following []string

	if excludeverified == "true" {
		following = getFollowings(goinsta, true)
	} else {
		following = getFollowings(goinsta, false)
	}

	for _, user := range following {
		if !containsUser(followers, user) {

			notFollowingBack = append(notFollowingBack, user)
			counter++

		}
	}

	if json == "true" {
		s.Stop()

		WriteToJSON(notFollowingBack)
	} else {
		s.Stop()
		createUnfollowTable(goinsta, notFollowingBack)

	}
}

func getFollowings(goinsta *goinsta.Instagram, verified bool) []string {
	following := goinsta.Account.Following("", "default")
	followingUsers := []string{}
	if verified {
		following.Next()
		if !following.BigList {
			for _, user := range following.Users {
				if !user.IsVerified {
					followingUsers = append(followingUsers, user.Username)
				}
			}
		} else {
			for following.BigList {
				for _, user := range following.Users {
					if !user.IsVerified {
						followingUsers = append(followingUsers, user.Username)
					}
				}
				following.Next()

			}
			for _, user := range following.Users {
				if !user.IsVerified {
					followingUsers = append(followingUsers, user.Username)
				}
			}
		}
		return followingUsers

	} else {
		following.Next()
		if !following.BigList {
			for _, user := range following.Users {
				followingUsers = append(followingUsers, user.Username)
			}
		} else {
			for following.BigList {
				for _, user := range following.Users {
					followingUsers = append(followingUsers, user.Username)
				}
				following.Next()

			}
			for _, user := range following.Users {
				followingUsers = append(followingUsers, user.Username)
			}
		}
		return followingUsers

	}
}

func getFollowers(goinsta *goinsta.Instagram) []string {
	followers := goinsta.Account.Followers("")
	followersUsers := []string{}

	followers.Next()
	if !followers.BigList {
		for _, user := range followers.Users {
			followersUsers = append(followersUsers, user.Username)
		}
	} else {
		for followers.BigList {
			for _, user := range followers.Users {
				followersUsers = append(followersUsers, user.Username)
			}
			followers.Next()
		}
		for _, user := range followers.Users {
			followersUsers = append(followersUsers, user.Username)
		}
	}

	return followersUsers
}

func containsUser(users []string, target string) bool {
	for _, user := range users {
		if user == target {
			return true
		}
	}
	return false
}

type User struct {
	Username string `json:"username"`
}

func WriteToJSON(notFollowingBackUsers []string) error {
	users := make([]User, len(notFollowingBackUsers))
	for i, username := range notFollowingBackUsers {
		users[i] = User{Username: username}
	}

	file, err := os.Create("NFB.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(users); err != nil {
		return err
	}
	fmt.Println()
	fmt.Println(utils.Green("Data exported to " + file.Name() + " successfully!"))
	return nil
}

func createUnfollowTable(insta *goinsta.Instagram, users []string) {
	fmt.Print("\033[H\033[2J")
	currentIndex := 0
	pageSize := 10
	quit := false

	err := keyboard.Open()
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
			fmt.Println("No more users to unfollow!")
			break
		}

		startIndex := currentIndex
		endIndex := startIndex + pageSize
		if endIndex > len(users) {
			endIndex = len(users)
		}
		currentPageUsers := users[startIndex:endIndex]

		fmt.Print(utils.Blue("Users not following you back: " + (strconv.Itoa(startIndex + 1)) + "/" + (strconv.Itoa(len(users)) + "\n")))

		for i, u := range currentPageUsers {
			arrow := " "
			if i == currentIndex-startIndex {
				arrow = ">"
			}
			fmt.Printf("%s %d. %s\n", arrow, startIndex+i+1, u)
		}
		fmt.Println()
		fmt.Print(utils.Green("Navigation: Up/Down Arrow Keys\n"))
		fmt.Print(utils.Green("Enter: Unfollow the selected user\n"))
		fmt.Print(utils.Green("A: Unfollow all the users\n"))
		fmt.Print(utils.Green("Q: Quit\n"))

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyArrowUp && currentIndex > 0 {
			currentIndex--
		} else if key == keyboard.KeyArrowDown && currentIndex < len(users)-1 {
			currentIndex++
		} else if key == keyboard.KeyArrowDown && currentIndex == len(users)-1 {
			currentIndex = 0
		} else if key == keyboard.KeyEnter {
			Unfollow(insta, users[currentIndex])
			users[currentIndex] = utils.Green(users[currentIndex]) + utils.Green(" (Unfollowed Successfully)")
		} else if char == 'a' || char == 'A' {
			fmt.Print(utils.Red("Press 'a' again to confirm and unfollow all users . "))
			char2, _, _ := keyboard.GetSingleKey()
			if char2 == 'a' || char2 == 'A' {
				for _, user := range users {
					if currentIndex%10 == 0 {
						time.Sleep(5 * time.Second)
					}
					Unfollow(insta, user)
					users[currentIndex] = utils.Green(users[currentIndex]) + utils.Green(" (Unfollowed Successfully)")
				}
			}
		} else if char == 'q' || char == 'Q' {
			quit = true
		}

	}
}
