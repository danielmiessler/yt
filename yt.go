package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func getVideoID(url string) string {
	pattern := `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func getTranscript(videoID string) (string, error) {
	url := "https://www.youtube.com/watch?v=" + videoID
	resp, err := soup.Get(url)
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(resp)
	scriptTags := doc.FindAll("script")
	for _, scriptTag := range scriptTags {
		if strings.Contains(scriptTag.Text(), "captionTracks") {
			regex := regexp.MustCompile(`"captionTracks":(\[.*?\])`)
			match := regex.FindStringSubmatch(scriptTag.Text())
			if len(match) > 1 {
				var captionTracks []struct {
					BaseURL string `json:"baseUrl"`
				}
				json.Unmarshal([]byte(match[1]), &captionTracks)
				if len(captionTracks) > 0 {
					transcriptURL := captionTracks[0].BaseURL
					transcriptResp, err := soup.Get(transcriptURL)
					if err != nil {
						return "", err
					}
					return transcriptResp, nil
				}
			}
		}
	}
	return "", fmt.Errorf("transcript not found")
}

func getComments(service *youtube.Service, videoID string) []string {
	var comments []string

	call := service.CommentThreads.List([]string{"snippet", "replies"}).VideoId(videoID).TextFormat("plainText").MaxResults(100)
	response, err := call.Do()
	if err != nil {
		log.Printf("Failed to fetch comments: %v", err)
		return comments
	}

	for _, item := range response.Items {
		topLevelComment := item.Snippet.TopLevelComment.Snippet.TextDisplay
		comments = append(comments, topLevelComment)

		if item.Replies != nil {
			for _, reply := range item.Replies.Comments {
				replyText := reply.Snippet.TextDisplay
				comments = append(comments, "    - "+replyText)
			}
		}
	}

	return comments
}

func parseDuration(durationStr string) (int, error) {
	matches := regexp.MustCompile(`(?i)PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`).FindStringSubmatch(durationStr)
	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration string: %s", durationStr)
	}

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])

	return hours*60 + minutes + seconds/60, nil
}

func mainFunction(url string, options *Options) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting user home directory")
	}
	env_file := home_dir + "/.config/fabric/.env"
	err = godotenv.Load(env_file)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: YOUTUBE_API_KEY not found in ~/.config/fabric/.env. To add please run \"echo YOUTUBE_API_KEY=\"[Your API Key]\" >> ~/.config/fabric/.env\".")
	}

	videoID := getVideoID(url)
	if videoID == "" {
		log.Fatal("Invalid YouTube URL")
	}

	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	videoResponse, err := service.Videos.List([]string{"contentDetails"}).Id(videoID).Do()
	if err != nil {
		log.Fatalf("Error getting video details: %v", err)
	}

	durationStr := videoResponse.Items[0].ContentDetails.Duration
	durationMinutes, err := parseDuration(durationStr)
	if err != nil {
		log.Fatalf("Error parsing video duration: %v", err)
	}

	var comments []string
	if options.Comments {
		comments = getComments(service, videoID)
	}
	var transcriptText string
	transcript, err := getTranscript(videoID)
	if err != nil {
		transcriptText = fmt.Sprintf("Transcript not available. (%v)", err)
	} else {
		// Parse the XML transcript
		doc := soup.HTMLParse(transcript)
		// Extract the text content from the <text> tags
		textTags := doc.FindAll("text")
		var textBuilder strings.Builder
		for _, textTag := range textTags {
			textBuilder.WriteString(textTag.Text())
			textBuilder.WriteString(" ")
			transcriptText = textBuilder.String()
		}
	}

	if options.Duration {
		fmt.Println(durationMinutes)
	} else if options.Transcript {
		fmt.Println(transcriptText)
	} else if options.Comments {
		jsonComments, _ := json.MarshalIndent(comments, "", "  ")
		fmt.Println(string(jsonComments))
	} else {
		output := map[string]interface{}{
			"transcript": transcriptText,
			"duration":   durationMinutes,
			"comments":   comments,
		}
		jsonOutput, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonOutput))
	}
}

type Options struct {
	Duration   bool
	Transcript bool
	Comments   bool
	Lang       string
}

func main() {
	options := &Options{}
	flag.BoolVar(&options.Duration, "duration", false, "Output only the duration")
	flag.BoolVar(&options.Transcript, "transcript", false, "Output only the transcript")
	flag.BoolVar(&options.Comments, "comments", false, "Output the comments on the video")
	flag.StringVar(&options.Lang, "lang", "en", "Language for the transcript (default: English)")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Error: No URL provided.")
	}

	url := flag.Arg(0)
	mainFunction(url, options)
}
