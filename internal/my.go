package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fvrrvg/igo/utils"

	"github.com/Davincible/goinsta/v3"
	"github.com/eiannone/keyboard"
)

func MyStories(insta *goinsta.Instagram) {
	// // Stories
	stories, err := insta.Account.Stories()
	if err != nil {
		panic(err)
	}

	fmt.Println(utils.Yellow("Downloading your stories... "))
	finalDir := utils.GetCurrDir() + "/My/Stories/"

	if stories.Reel.Items == nil {
		fmt.Println(utils.Red("You have no stories posted now"))
		fmt.Println()
	} else {
		_, err = os.Stat(finalDir)
		if os.IsNotExist(err) {
			os.Mkdir(finalDir, 0o777)
		}
		for _, item := range stories.Reel.Items {

			item.DownloadTo(finalDir)
			counter += 1
			fmt.Fprintf(utils.Writer, utils.Yellow("%d / %d \n"), counter, len(stories.Reel.Items))
			utils.Writer.Flush()
		}
	}
}

func MyPosts(insta *goinsta.Instagram) {
	// Posts
	fmt.Println(utils.Yellow("Downloading your posts... "))
	finalDir := utils.GetCurrDir() + "/My/Posts/"

	posts := insta.Account.Feed()

	posts.Next()

	if posts.Items == nil {
		fmt.Println(utils.Red("You have no posts posted now"))
		fmt.Println()
	} else {
		_, err := os.Stat(finalDir)
		if os.IsNotExist(err) {
			os.Mkdir(finalDir, 0o777)
		}
		for _, item := range posts.Items {
			item.DownloadTo(finalDir)
			counter += 1
			fmt.Fprintf(utils.Writer, utils.Yellow("%d / %d \n"), counter, len(posts.Items))
			utils.Writer.Flush()
		}
		fmt.Fprint(utils.Writer, "\r")
		utils.Writer.Flush()

		fmt.Println(utils.Green("All your posts downloaded successfully!"))
		fmt.Println()
	}
}

func MyIGTV(insta *goinsta.Instagram) {
	// Igtv
	fmt.Println(utils.Yellow("Downloading your igtvs... "))
	finalDir := utils.GetCurrDir() + "/My/Igtv/"

	igtv := insta.IGTV

	if igtv.Items == nil {
		fmt.Println(utils.Red("You have no igtvs posted now"))
		fmt.Println()
	} else {
		_, err := os.Stat(finalDir)
		if os.IsNotExist(err) {
			os.Mkdir(finalDir, 0o777)
		}
		for _, item := range igtv.Items {
			item.DownloadTo(finalDir)
			counter += 1
			fmt.Fprintf(utils.Writer, utils.Yellow("%d / %d \n"), counter, len(igtv.Items))
			utils.Writer.Flush()
		}

		fmt.Println(utils.Green("All your IGTVS downloaded successfully!"))
	}
}

func MyHL(insta *goinsta.Instagram) {
	// Highlights
	fmt.Println(utils.Yellow("Downloading your highlights... "))
	finalDir := utils.GetCurrDir() + "/My/Highlights/"

	profile, err := insta.VisitProfile(insta.Account.Username)
	if err != nil {
		panic(err)
	}

	hl, err := profile.User.Highlights()
	if err != nil {
		panic(err)
	}

	if len(hl) == 0 {
		fmt.Println(utils.Red("You have no highlights posted now"))
		fmt.Println()
	} else {
		for _, v := range hl {
			v.Sync()
			counter = 0
			for _, i := range v.Items {
				i.DownloadTo(finalDir + "/" + v.Title + "/")
				counter += 1
				fmt.Fprintf(utils.Writer, utils.Yellow("%d / %d \n"), counter, len(v.Items))
				utils.Writer.Flush()
			}
		}
		fmt.Fprint(utils.Writer, "\r")
		utils.Writer.Flush()
		fmt.Println(utils.Green("All your highlights downloaded successfully!"))
		fmt.Println()
	}
}

func MyPP(insta *goinsta.Instagram) {
	// Profile Picture
	fmt.Println(utils.Yellow("Downloading your profile picture... "))
	finalDir := utils.GetCurrDir() + "/My/"

	url := insta.Account.ProfilePicURL

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching profile picture:", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading profile picture response:", err)
	}

	err = os.WriteFile(finalDir+"/Profile Picture.jpg", body, 0o664)
	if err != nil {
		fmt.Println("Error writing profile picture to file:", err)
	} else {
		fmt.Println(utils.Green("Your Profile picture downloaded successfully!"))
	}

	fmt.Println()
}

func MyCollections(insta *goinsta.Instagram) {
	fmt.Println(utils.Yellow("Getting your saved collections..."))

	collections := insta.Collections
	var colNames []string
	finalDir := utils.GetCurrDir() + "/My/Saved/"

	for collections.Next() {
		collections.Next()
	}

	if collections.Items == nil {
		fmt.Println(utils.Red("You have no saved collections now"))
		fmt.Println()
	} else {

		for _, collection := range collections.Items {
			colNames = append(colNames, collection.Name)
		}
		choice := SavedCollectionsMenu(insta, colNames)

		if choice == nil {
			fmt.Println(utils.Red("0 Collections Selected!"))
			fmt.Println()
			return
		}
		fmt.Println(utils.Yellow("Downloading your selected collections..."))

		for _, collection := range collections.Items {
			for _, choice := range choice {
				if collection.Name == choice {

					collection.Next()
					counter = 0
					for _, item := range collection.Items {
						item.DownloadTo(finalDir + "/" + collection.Name + "/")
						counter += 1
						fmt.Fprintf(utils.Writer, utils.Yellow("%d / %d \n"), counter, len(collection.Items))
						utils.Writer.Flush()
					}
					fmt.Fprint(utils.Writer, "\r")
					utils.Writer.Flush()

				}
			}
		}

		fmt.Println(utils.Green("All your saved collections downloaded successfully!"))
		fmt.Println()
	}
}

func SavedCollectionsMenu(insta *goinsta.Instagram, colNames []string) []string {
	fmt.Println(utils.Yellow("Choose collections to download (toggle selection with Space):"))
	fmt.Print("\033[H\033[2J")
	currentIndex := 0
	pageSize := 10
	selectedCollections := make(map[string]bool)
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
		if currentIndex >= len(colNames) {
			fmt.Println("No more collections!")
			break
		}

		startIndex := currentIndex
		endIndex := startIndex + pageSize
		if endIndex > len(colNames) {
			endIndex = len(colNames)
		}
		currentPageColl := colNames[startIndex:endIndex]

		fmt.Print(utils.Blue("Your Saved Collections: " + strconv.Itoa(startIndex+1) + "/" + strconv.Itoa(len(colNames)) + "\n"))

		for i, u := range currentPageColl {
			arrow := " "
			if i == currentIndex-startIndex {
				arrow = ">"
			}

			selectionIndicator := " "
			if selectedCollections[u] {
				selectionIndicator = " [✓]"
			}

			fmt.Printf("%s %d. %s%s\n", arrow, startIndex+i+1, u, utils.Green(selectionIndicator))
		}

		fmt.Println()
		fmt.Print(utils.Green("Navigation: Up/Down Arrow Keys\n"))
		fmt.Print(utils.Green("Space: Toggle collection selection\n"))
		fmt.Print(utils.Green("Enter: Confirm selected collections\n"))
		fmt.Print(utils.Green("Q: Quit\n"))
		fmt.Println()

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyArrowUp && currentIndex > 0 {
			currentIndex--
		} else if key == keyboard.KeyArrowDown && currentIndex < len(colNames)-1 {
			currentIndex++
		} else if key == keyboard.KeyArrowDown && currentIndex == len(colNames)-1 {
			currentIndex = 0
		} else if key == keyboard.KeyEnter {
			selected := []string{}
			for name, isSelected := range selectedCollections {
				if isSelected {
					selected = append(selected, name)
				}
			}
			return selected
		} else if key == keyboard.KeySpace {
			selectedCollections[colNames[currentIndex]] = !selectedCollections[colNames[currentIndex]]
		} else if char == 'q' || char == 'Q' {
			quit = true
		} else if key == keyboard.KeyCtrlC {
			quit = true
		}
	}

	return nil
}

func Whoami(insta *goinsta.Instagram) {
	fmt.Println(utils.Yellow("Getting your info..."))

	flr := getFollowers(insta)
	flg := getFollowings(insta, false)

	fmt.Println(utils.Green("Username: " + insta.Account.Username))
	fmt.Println(utils.Green("Full Name: " + insta.Account.FullName))
	fmt.Println(utils.Green("Bio: " + insta.Account.Biography))
	fmt.Println(utils.Green("Followers: " + strconv.Itoa(len(flr))))
	fmt.Println(utils.Green("Following: " + strconv.Itoa(len(flg))))

	fmt.Println()
}
