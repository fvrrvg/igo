package parts

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"igo/utils"

	"github.com/Davincible/goinsta/v3"
)

var (
	counterM int = 1
	counterV int = 1
	counterI int = 1
)

func GetDMS(insta *goinsta.Instagram, u *goinsta.User) {
	fmt.Print(utils.Yellow("Downloading " + u.Username + "'s DMs...\n"))
	inbox := insta.Inbox
	err := inbox.Sync()
	if err != nil {
		panic(err)
	}

	hasDms := false

	for inbox.Next() {
		if !inbox.HasOlder {
			continue
		}
	}

	for _, item := range inbox.Conversations {
		for _, user := range item.Users {
			if user.Username == u.Username && !item.IsGroup {
				hasDms = true
				for item.Next() {
					if !item.HasOlder {
						break
					}
				}

				finalDir := filepath.Join(utils.Dir(u), "Dms/")
				if err := os.MkdirAll(finalDir, 0o755); err != nil {
					panic(err)
				}

				dmsFilePath := filepath.Join(finalDir, "dms.txt")
				file, err := os.Create(dmsFilePath)
				if err != nil {
					panic(err)
				}

				defer file.Close()

				for _, msg := range item.Items {
					switch msg.Type {
					case "text":
						if msg.UserID == insta.Account.ID {
							if _, err := file.WriteString("You: " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						} else {
							if _, err := file.WriteString(u.Username + ": " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						}
					case "media":
						if msg.UserID == insta.Account.ID {
							msg.Media.DownloadTo(finalDir + "/")
							if msg.Media.MediaType == 1 {

								parsedURL, err := url.Parse(msg.Media.Images.GetBest())
								if err != nil {
									panic(err)
								}

								pathSegments := strings.Split(parsedURL.Path, "/")

								filename := pathSegments[len(pathSegments)-1]
								err = os.Rename(finalDir+"/"+filename, finalDir+"/[MEDIA FILE "+strconv.Itoa(counterM)+"].jpg")
								if err != nil {
									panic(err)
								}
							} else {
								parsedURL, err := url.Parse(msg.Media.Videos[0].URL)
								if err != nil {
									panic(err)
								}

								pathSegments := strings.Split(parsedURL.Path, "/")

								filename := pathSegments[len(pathSegments)-1]

								err = os.Rename(finalDir+"/"+filename, finalDir+"/[MEDIA FILE "+strconv.Itoa(counterM)+"].mp4")
								if err != nil {
									panic(err)
								}

							}

							if _, err := file.WriteString("You: [SENT MEDIA FILE " + strconv.Itoa(counterM) + "] \n"); err != nil {
								panic(err)
							}

							counterM++
						} else {
							msg.Media.DownloadTo(finalDir + "/")
							if msg.Media.MediaType == 1 {

								parsedURL, err := url.Parse(msg.Media.Images.GetBest())
								if err != nil {
									panic(err)
								}

								pathSegments := strings.Split(parsedURL.Path, "/")

								filename := pathSegments[len(pathSegments)-1]
								err = os.Rename(finalDir+"/"+filename, finalDir+"/[MEDIA FILE "+strconv.Itoa(counterM)+"].jpg")
								if err != nil {
									panic(err)
								}
							} else {
								parsedURL, err := url.Parse(msg.Media.Videos[0].URL)
								if err != nil {
									panic(err)
								}

								pathSegments := strings.Split(parsedURL.Path, "/")

								filename := pathSegments[len(pathSegments)-1]

								err = os.Rename(finalDir+"/"+filename, finalDir+"/[MEDIA FILE "+strconv.Itoa(counterM)+"].mp4")
								if err != nil {
									panic(err)
								}

							}

							if _, err := file.WriteString(u.Username + ":  [SENT MEDIA FILE " + strconv.Itoa(counterM) + "] \n"); err != nil {
								panic(err)
							}
							counterM++

						}

					case "voice_media":
						s := DownloadVN(msg.VoiceMedia.Media.Audio.AudioSrc, finalDir)
						if msg.UserID == insta.Account.ID {
							if _, err := file.WriteString("You: [SENT " + s + "]\n"); err != nil {
								panic(err)
							}
						} else {
							if _, err := file.WriteString(u.Username + ": [SENT " + s + "]\n"); err != nil {
								panic(err)
							}
						}
					case "link":
						if msg.UserID == insta.Account.ID {
							if _, err := file.WriteString("You: " + msg.Link.Text + "\n" + msg.Link.Context.URL + "\n"); err != nil {
								panic(err)
							}
						} else {
							if _, err := file.WriteString(u.Username + ": " + msg.Link.Text + "\n" + msg.Link.Context.URL + "\n"); err != nil {
								panic(err)
							}
						}
					case "expiutils.Red_placeholder":
						if msg.UserID == insta.Account.ID {
							if len(msg.Text) == 0 {
								if _, err := file.WriteString("You: [MENTIONED HIM/HER IN A STORY] " + "\n"); err != nil {
									panic(err)
								}
							} else {
								if _, err := file.WriteString("You: [REPLIED TO THEIR STORY] " + msg.Text + "\n"); err != nil {
									panic(err)
								}
							}
						} else {
							if len(msg.Text) == 0 {
								if _, err := file.WriteString(u.Username + ": [MENTIONED YOU IN THEIR STORY] " + msg.Text + "\n"); err != nil {
									panic(err)
								}
							} else {
								if _, err := file.WriteString(u.Username + ": [REPLIED TO YOUR STORY] " + msg.Text + "\n"); err != nil {
									panic(err)
								}
							}
						}
					case "placeholder":
						if msg.UserID == insta.Account.ID {
							if _, err := file.WriteString("You: [HIDDEN POST] " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						} else {
							if _, err := file.WriteString(u.Username + ": [HIDDEN POST]" + msg.Text + "\n"); err != nil {
								panic(err)
							}
						}

					case "media_share":
						if msg.UserID == insta.Account.ID {
							if _, err := file.WriteString("You: " + "[SHAutils.Red A POST] " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						} else {
							if _, err := file.WriteString(u.Username + ": " + "[SHAutils.Red A POST] " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						}

					case "clip":
						if msg.UserID == insta.Account.ID {
							if _, err := file.WriteString("You: [SHAutils.Red A REEL] " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						} else {
							if _, err := file.WriteString(u.Username + ": [SHAutils.Red A REEL] " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						}

					case "xma_reel_share":
						if msg.UserID == insta.Account.ID {
							if _, err := file.WriteString("You: " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						} else {
							if _, err := file.WriteString(u.Username + ": " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						}

					case "video_call_event":
						if msg.UserID == insta.Account.ID {
							if _, err := file.WriteString("You: [AUDIO/VIDEO CALL] " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						} else {
							if _, err := file.WriteString(u.Username + ": [AUDIO/VIDEO CALL] " + msg.Text + "\n"); err != nil {
								panic(err)
							}
						}

					}
				}
				ReverseFile(dmsFilePath)
				fmt.Print(utils.Green("All " + u.Username + "'s DMs downloaded successfully!\n"))
				fmt.Println()
			}
		}
	}
	if !hasDms {
		fmt.Println(utils.Red(u.Username + " has no DMs with you!"))
		fmt.Println()
	}
}

func DownloadVN(url string, dir string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	s := "VOICE MESSAGE " + strconv.Itoa(counterV)

	err = os.WriteFile(filepath.Join(dir, s+".mp3"), body, 0o644)
	if err != nil {
		panic(err)
	}

	counterV++
	return s
}

func DownloadVIDEO(url string, dir string, id string) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(dir, id+".mp4"), body, 0o644)
	if err != nil {
		panic(err)
	}
}

func DownloadIMG(url string, dir string) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(dir, strconv.Itoa(counterI)+".jpg"), body, 0o644)
	if err != nil {
		panic(err)
	}

	counterI++
}

func ReverseFile(filePath string) {
	lines, err := readLines(filePath)
	if err != nil {
		panic(err)
	}

	reverseLines(lines)

	if err := writeLines(filePath, lines); err != nil {
		panic(err)
	}
}

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func reverseLines(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}

func writeLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	return writer.Flush()
}
