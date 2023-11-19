package internal

import (
	"fmt"
	"os"

	"github.com/fvrrvg/igo/utils"

	"github.com/Davincible/goinsta/v3"
)

func GetAllHighlights(u *goinsta.User) {
	fmt.Println(utils.Yellow("Downloading " + u.Username + "'s highlights..."))

	hl, err := u.Highlights()
	if err != nil {
		panic(err)
	}

	profile, err := u.VisitProfile()
	if err != nil {
		panic(err)
	}

	finalDir := utils.GetCurrDir() + "/" + u.Username + "/Highlights/"

	if !profile.Friendship.Following && u.IsPrivate {
		fmt.Println(utils.Red("This account is private and you are not following it!"))
		fmt.Println()

	} else if len(hl) == 0 {
		fmt.Println(utils.Red(u.Username + " has no highlights!"))
	} else {

		_, err := os.Stat(finalDir)
		if err != nil {
			os.Mkdir(finalDir, 0o755)
		}

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

		fmt.Println(utils.Green(u.Username + "'s Highlights downloaded successfully!"))
	}
}
