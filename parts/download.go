package parts

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"igo/utils"

	"github.com/Davincible/goinsta/v3"
	"github.com/PuerkitoBio/goquery"
)

func Download(insta *goinsta.Instagram, url string) {
	fmt.Println(utils.Yellow("Downloading..."))

	// check if url is valid
	if !strings.Contains(url, "instagram.com/") {
		fmt.Println(utils.Red("Invalid URL!"))
		fmt.Println("igo download https://www.instagram.com/...")
		return
	}

	code := ExtractCode(url)

	if len(code) > 12 {
		username, err := ExtractUsernameandCaption(url, "private")
		if err != nil {
			log.Fatal("Error:", err)
		}
		u := utils.GetUser(insta, username)

		feed := u.Feed()

		for feed.Next() {
			feed.Next()
		}

		finalDir := utils.GetCurrDir()

		for _, item := range feed.Items {
			if item.Code == code {
				item.DownloadTo(finalDir)
			}
		}

		fmt.Println(utils.Green("Downloaded Successfully!"))

	} else {
		username, err := ExtractUsernameandCaption(url, "public")
		if err != nil {
			log.Fatal("Error:", err)
		}
		u := utils.GetUser(insta, username)

		feed := u.Feed()

		for feed.Next() {
			feed.Next()
		}

		finalDir := utils.GetCurrDir()

		for _, item := range feed.Items {
			if item.Code == code {
				item.DownloadTo(finalDir)
			}
		}
		fmt.Println(utils.Green("Downloaded Successfully!"))

	}
}

func SanitizeURL(url string) string {
	parts := strings.Split(url, "?")
	return parts[0]
}

func ExtractCode(url string) string {
	sanitizedURL := SanitizeURL(url)
	parts := strings.Split(sanitizedURL, "/")
	code := parts[len(parts)-2]
	if strings.Contains(code, "?") {
		code = strings.Split(code, "?")[0]
	}

	return code
}

func ExtractUsernameandCaption(link string, accountType string) (string, error) {
	sanitizedURL := SanitizeURL(link)
	client := &http.Client{}
	req, err := http.NewRequest("GET", sanitizedURL, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	if accountType == "public" {
		metaTitle := doc.Find("meta[name='twitter:title']")
		if metaTitle.Length() == 0 {
			return "", fmt.Errorf(utils.Red("Something went wrong!"))
		}

		usernameContent := metaTitle.AttrOr("content", "")

		re := regexp.MustCompile(`@([^•]+) •`)
		match := re.FindStringSubmatch(usernameContent)
		if len(match) >= 2 {
			username := match[1]
			username = strings.Replace(username, ")", "", -1)

			return username, nil
		}

		return "", fmt.Errorf(utils.Red("Something went wrong!"))
	} else if accountType == "private" {

		titleElement := doc.Find("title")
		if titleElement.Length() == 0 {
			return "", fmt.Errorf(utils.Red("Something went wrong!"))
		}

		titleText := titleElement.Text()

		re := regexp.MustCompile(`@([^•]+) •`)
		match := re.FindStringSubmatch(titleText)
		if len(match) >= 2 {
			username := match[1]
			username = strings.Replace(username, ")", "", -1)

			return username, nil
		}

		return "", fmt.Errorf(utils.Red("Something went wrong!"))
	}

	return "", fmt.Errorf(utils.Red("Something went wrong!"))
}
