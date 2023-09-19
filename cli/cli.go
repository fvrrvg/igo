package cli

import (
	"fmt"
	"os"

	"igo/parts"
	"igo/utils"

	"github.com/Davincible/goinsta/v3"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

var insta *goinsta.Instagram = parts.Login()

func downloadPosts(username string, latest bool) {
	if username == "" {
		fmt.Println(utils.Red("You have to enter a username!"))
		fmt.Println("igo posts -u username")
		os.Exit(1)
	}

	if !latest {
		parts.GetAllPosts(utils.GetUser(insta, username), true)
	} else {
		parts.GetAllPosts(utils.GetUser(insta, username), false)
	}
}

func StartCli() {
	var username string

	rootCmd := &cobra.Command{
		Use:   "igo",
		Short: "A CLI for Instagram built in Go.",
		Long: `A CLI for Instagram built in Go.
Helps you do some tasks faster and easier like downloading posts, stories, IGTV, and DMs.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to igo!")
			cmd.Help()
		},
	}

	followCmd := &cobra.Command{
		Use:   "follow",
		Short: "Follow a user.",
		Long:  `Follow a user. You can follow a user by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("username").Value.String() == "" {
				fmt.Println(utils.Red("You have to enter a username!"))
				fmt.Println("igo follow -u username")

				os.Exit(1)
			}
			parts.Follow(insta, username)
		},
	}

	unfollowCmd := &cobra.Command{
		Use:   "unfollow",
		Short: "Unfollow a user.",
		Long:  `Unfollow a user. You can unfollow a user by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("username").Value.String() == "" {
				fmt.Println(utils.Red("You have to enter a username!"))
				fmt.Println("igo unfollow -u username")

				os.Exit(1)
			}
			parts.Unfollow(insta, username)
		},
	}

	downloadstoriesCmd := &cobra.Command{
		Use:   "stories",
		Short: "Download the stories of a user.",
		Long:  `Download the stories of a user. You can download the stories of a user by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("username").Value.String() == "" {
				fmt.Println(utils.Red("You have to enter a username!"))
				fmt.Println("igo stories -u username")

				os.Exit(1)
			}
			parts.GetStories(utils.GetUser(insta, username))
		},
	}

	downloadigtvCmd := &cobra.Command{
		Use:   "igtv",
		Short: "Download the IGTVs of a user.",
		Long:  `Download the IGTVs of a user. You can download the IGTVs of a user by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("username").Value.String() == "" {
				fmt.Println(utils.Red("You have to enter a username!"))
				fmt.Println("igo igtv -u username")

				os.Exit(1)
			}

			parts.GetAllIgtv(utils.GetUser(insta, username))
		},
	}

	downloadhighlightsCmd := &cobra.Command{
		Use:   "highlights",
		Short: "Download the Highlights of a user.",
		Long:  `Download the Highlights of a user. You can download the Highlights of a user by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("username").Value.String() == "" {
				fmt.Println(utils.Red("You have to enter a username!"))
				fmt.Println("igo igtv -u username")

				os.Exit(1)
			}

			parts.GetAllHighlights(utils.GetUser(insta, username))
		},
	}

	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from your account.",
		Long:  `Logout from your account.`,
		Run: func(cmd *cobra.Command, args []string) {
			parts.Logout()
		},
	}

	dmsCmd := &cobra.Command{
		Use:   "dms",
		Short: "Download your direct messages with a user.",
		Long:  `Download your direct messages with a user. You can download your direct messages with a user by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("username").Value.String() == "" {
				fmt.Println(utils.Red("You have to enter a username!"))
				fmt.Println("igo dms -u username")

				os.Exit(1)
			}
			parts.GetDMS(insta, utils.GetUser(insta, username))
		},
	}

	ppCmd := &cobra.Command{
		Use:   "pp",
		Short: "Download the profile picture of a user.",
		Long:  `Download the profile picture of a user. You can download the profile picture of a user by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("username").Value.String() == "" {
				fmt.Println(utils.Red("You have to enter a username!"))
				fmt.Println("igo pp -u username")

				os.Exit(1)
			}
			parts.GetProfilePic(utils.GetUser(insta, username))
		},
	}

	downloadpostsCmd := &cobra.Command{
		Use:   "posts",
		Short: "Download the posts of a user.",
		Long:  `Download the posts of a user. You can download the posts of a user by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			username := cmd.Flag("username").Value.String()
			latest := cmd.Flag("latest").Value.String() == "false"
			downloadPosts(username, latest)
		},
	}

	everythingCmd := &cobra.Command{
		Use:   "everything",
		Short: "Download everything. (PP, Posts, Stories, IGTV, and DMS)",
		Long:  `Download everything (PP, Posts, Stories, IGTV, and DMS) by entering their username.`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("username").Value.String() == "" {
				fmt.Println(utils.Red("You have to enter a username!"))
				fmt.Println("igo everything -u username")
				os.Exit(1)
			}
			downloadstoriesCmd.Run(cmd, args)
			downloadPosts(username, true)
			downloadigtvCmd.Run(cmd, args)
			dmsCmd.Run(cmd, args)
			ppCmd.Run(cmd, args)
			downloadhighlightsCmd.Run(cmd, args)
		},
	}

	notfollowingbackCmd := &cobra.Command{
		Use:   "nfb",
		Short: "Get the users that don't follow you back.",
		Long:  `Get the users that don't follow you back.`,
		Run: func(cmd *cobra.Command, args []string) {
			jsonflag := cmd.Flag("json").Value.String()
			notverifiedflag := cmd.Flag("exclude-verified").Value.String()

			parts.NotFollowingBack(insta, jsonflag, notverifiedflag)
		},
	}

	myCmd := &cobra.Command{
		Use:   "my",
		Short: "Download your own posts, stories, IGTV, Highlights",
		Long:  `Download your own posts, stories, IGTV, Highlights`,
		Run: func(cmd *cobra.Command, args []string) {
			parts.My(insta)
		},
	}

	whoamiCmd := &cobra.Command{
		Use:   "whoami",
		Short: "Get your account's info.",
		Long:  `Get your account's info.`,
		Run: func(cmd *cobra.Command, args []string) {
			parts.Whoami(insta)
		},
	}

	downloadCmd := &cobra.Command{
		Use:   "download",
		Short: "Download post,iptv or reel by entering the url.",
		Long:  `Download post,iptv or reel by entering the url.`,
		Run: func(cmd *cobra.Command, args []string) {
			urlFlag := cmd.Flag("url").Value.String()
			parts.Download(insta, urlFlag)
		},
	}

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.DisableSuggestions = false
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if !rootCmd.HasHelpSubCommands() {
		rootCmd.SetHelpCommand(&cobra.Command{
			Use:    "no-help",
			Hidden: true,
		})
	}

	var jsonFlag bool
	var latestFlag bool
	var notverifiedFlag bool
	var urlFlag string

	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username of the user you want to perform tasks on.")
	notfollowingbackCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Export the users that don't follow you back in JSON format.")
	notfollowingbackCmd.Flags().BoolVarP(&notverifiedFlag, "exclude-verified", "e", false, "Exclude verified users from the list.")
	downloadpostsCmd.Flags().BoolVarP(&latestFlag, "latest", "l", false, "Download latest post only")
	downloadCmd.Flags().StringVarP(&urlFlag, "url", "", "", "Url of the post,igtv or reel you want to download.")

	rootCmd.AddCommand(followCmd)
	rootCmd.AddCommand(unfollowCmd)
	rootCmd.AddCommand(downloadstoriesCmd)
	rootCmd.AddCommand(downloadpostsCmd)
	rootCmd.AddCommand(downloadigtvCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(dmsCmd)
	rootCmd.AddCommand(everythingCmd)
	rootCmd.AddCommand(ppCmd)
	rootCmd.AddCommand(notfollowingbackCmd)
	rootCmd.AddCommand(myCmd)
	rootCmd.AddCommand(downloadhighlightsCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(downloadCmd)

	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.HiGreen + cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.HiRed + cc.Bold,
	})

	rootCmd.Root().Example = `  igo follow -u username 
  igo unfollow -u username
  igo everything -u username
  igo nfb`

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
