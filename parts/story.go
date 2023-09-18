package parts

import (
	"fmt"

	"igo/utils"

	"github.com/Davincible/goinsta/v3"
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
