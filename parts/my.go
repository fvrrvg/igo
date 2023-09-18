package parts

import (
	"fmt"
	"os"
	"strconv"

	"igo/utils"

	"github.com/Davincible/goinsta/v3"
)

func My(insta *goinsta.Instagram) {
	// Stories
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

	// Posts
	fmt.Println(utils.Yellow("Downloading your posts... "))
	finalDir = utils.GetCurrDir() + "/My/Posts/"

	posts := insta.Account.Feed()

	posts.Next()

	if posts.Items == nil {
		fmt.Println(utils.Red("You have no posts posted now"))
		fmt.Println()
	} else {
		_, err = os.Stat(finalDir)
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

	// Igtv
	fmt.Println(utils.Yellow("Downloading your igtvs... "))
	finalDir = utils.GetCurrDir() + "/My/Igtv/"

	igtv := insta.IGTV

	if igtv.Items == nil {
		fmt.Println(utils.Red("You have no igtvs posted now"))
		fmt.Println()
	} else {
		_, err = os.Stat(finalDir)
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

	// Highlights
	fmt.Println(utils.Yellow("Downloading your highlights... "))
	finalDir = utils.GetCurrDir() + "/My/Highlights/"

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

	}
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
