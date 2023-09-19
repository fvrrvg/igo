package parts

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"igo/utils"

	"github.com/Davincible/goinsta/v3"
)

var counter int = 0

func GetProfilePic(u *goinsta.User) error {
	finalDir := utils.Dir(u)
	profilePicURL := u.HdProfilePicURLInfo.URL

	fmt.Println(utils.Yellow("Downloading " + u.Username + "'s profile picture..."))

	res, err := http.Get(profilePicURL)
	if err != nil {
		fmt.Println("Error fetching profile picture:", err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading profile picture response:", err)
		return err
	}

	err = os.WriteFile(finalDir+"/Profile Picture.jpg", body, 0o664)
	if err != nil {
		fmt.Println("Error writing profile picture to file:", err)
		return err
	} else {
		fmt.Println(utils.Green(u.Username + "'s Profile picture downloaded successfully!"))
		fmt.Println()
	}
	return nil
}

func GetAllPosts(u *goinsta.User, flag bool) {
	if flag {
		fmt.Println(utils.Yellow("Downloading " + u.Username + "'s latest post..."))
	} else {
		fmt.Println(utils.Yellow("Downloading " + u.Username + "'s posts... "))
	}
	profile, err := u.VisitProfile()
	if err != nil {
		panic(err)
	}
	if !u.Friendship.Following && profile.User.IsPrivate {
		fmt.Println(utils.Red("This account is private and you are not following it!"))
		fmt.Println()

	} else {
		if len(profile.Feed.Latest()) == 0 {
			fmt.Println(utils.Red("This account has no posts"))
			fmt.Println()

		} else {

			finalDir := utils.Dir(u)
			if err != nil {
				panic(err)
			}
			if flag {
				profile.Feed.Items[0].DownloadTo(finalDir + "Posts/")
				fmt.Println(utils.Green("Latest " + u.Username + "'s post is downloaded successfully!"))
				fmt.Println()

			} else {
				utils.Writer.Start()
				for profile.Feed.Next() {
					profile.Feed.Next()
				}
				for _, item := range profile.Feed.Items {
					item.DownloadTo(finalDir + "Posts/")
					counter += 1
					fmt.Fprintf(utils.Writer, utils.Yellow("%d / %d \n"), counter, len(profile.Feed.Items))
					utils.Writer.Flush()

				}
				fmt.Fprint(utils.Writer, "\r")
				utils.Writer.Stop()

				fmt.Println(utils.Green("All " + u.Username + "'s posts downloaded successfully!"))
				fmt.Println()
			}
		}
	}
}

func GetAllIgtv(u *goinsta.User) {
	fmt.Println(utils.Yellow("Downloading " + u.Username + "'s IGTVs..."))

	profile, err := u.VisitProfile()
	if err != nil {
		panic(err)
	}
	if !u.Friendship.Following && profile.User.IsPrivate {
		fmt.Println(utils.Red("This account is private and you are not following it!"))
		fmt.Println()

	} else {
		if profile.User.IGTVCount == 0 {
			fmt.Println(utils.Red("This account has no IGTVs"))
			fmt.Println()

		} else {
			finalDir := utils.Dir(u)

			for profile.IGTV.Next() {
				profile.IGTV.Next()
			}

			for _, item := range profile.IGTV.Items {
				item.DownloadTo(finalDir + "IGTV/")
				counter += 1
				fmt.Fprintf(utils.Writer, utils.Yellow("%d / %d \n"), counter, len(profile.IGTV.Items))
				utils.Writer.Flush()

			}
			fmt.Fprint(utils.Writer, "\r")
			utils.Writer.Stop()

			fmt.Println(utils.Green("All " + u.Username + "'s IGTVs downloaded successfully!"))
			fmt.Println()

		}
	}
}
