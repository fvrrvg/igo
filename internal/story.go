package internal

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/fvrrvg/igo/utils"

	"github.com/Davincible/goinsta/v3"
	"github.com/briandowns/spinner"
	"github.com/eiannone/keyboard"
)

func GetStories(u *goinsta.User) {
	fmt.Println(utils.Yellow("Downloading " + u.Username + "'s stories... "))

	stories, err := u.Stories()
	if err != nil {
		panic(err)
	}

	profile, err := u.VisitProfile()
	if err != nil {
		panic(err)
	}

	if stories.Reel.Items == nil {
		if !u.Friendship.Following && profile.User.IsPrivate {
			fmt.Println(utils.Red("This account is private and you are not following it!"))
			fmt.Println()

		} else {
			fmt.Println(utils.Red("This account has no stories posted now"))
			fmt.Println()

		}
	} else {
		counter := 0
		for _, item := range stories.Reel.Items {
			finalDir := utils.Dir(u)
			item.DownloadTo(finalDir + "Stories/")
			counter += 1
			fmt.Fprintf(utils.Writer, utils.Yellow("%d / %d \n"), counter, len(stories.Reel.Items))
			utils.Writer.Flush()

		}

		fmt.Fprint(utils.Writer, "\r")
		utils.Writer.Flush()

		fmt.Println(utils.Green("All " + u.Username + "'s Stories downloaded successfully!"))
		fmt.Println()

	}
}

func GetStoriesNow(insta *goinsta.Instagram) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.Color("yellow")
	s.Start()

	tl := insta.Timeline
	tl.Next()
	stories := tl.Stories()
	fmt.Println(stories)
	var usernames []string
	for _, s := range stories {
		usernames = append(usernames, s.User.Username)
	}
	s.Stop()
	final := StoriesMenu(usernames)

	for _, u := range final {
		user := utils.GetUser(insta, u)
		GetStories(user)
	}
}

func StoriesMenu(list []string) []string {
	fmt.Print("\033[H\033[2J")
	currentIndex := 0
	pageSize := 10
	selectedStories := make(map[string]bool)
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
		if currentIndex >= len(list) {
			fmt.Println("No more users!")
			break
		}

		startIndex := currentIndex
		endIndex := startIndex + pageSize
		if endIndex > len(list) {
			endIndex = len(list)
		}
		currentPageColl := list[startIndex:endIndex]

		fmt.Print(utils.Blue("Users that have story posted now: " + strconv.Itoa(startIndex+1) + "/" + strconv.Itoa(len(list)) + "\n"))

		for i, u := range currentPageColl {
			arrow := " "
			if i == currentIndex-startIndex {
				arrow = ">"
			}

			selectionIndicator := " "
			if selectedStories[u] {
				selectionIndicator = " [✓]"
			}

			fmt.Printf("%s %d. %s%s\n", arrow, startIndex+i+1, u, utils.Green(selectionIndicator))
		}

		fmt.Println()
		fmt.Print(utils.Green("Navigation: Up/Down Arrow Keys\n"))
		fmt.Print(utils.Green("Space: Toggle story selection\n"))
		fmt.Print(utils.Green("Enter: Confirm selected stories\n"))
		fmt.Print(utils.Green("Q: Quit\n"))
		fmt.Println()

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyArrowUp && currentIndex > 0 {
			currentIndex--
		} else if key == keyboard.KeyArrowDown && currentIndex < len(list)-1 {
			currentIndex++
		} else if key == keyboard.KeyArrowDown && currentIndex == len(list)-1 {
			currentIndex = 0
		} else if key == keyboard.KeyEnter {
			selected := []string{}
			for name, isSelected := range selectedStories {
				if isSelected {
					selected = append(selected, name)
				}
			}
			return selected
		} else if key == keyboard.KeySpace {
			selectedStories[list[currentIndex]] = !selectedStories[list[currentIndex]]
		} else if char == 'q' || char == 'Q' {
			quit = true
		} else if key == keyboard.KeyCtrlC {
			quit = true
		}
	}

	return nil
}
