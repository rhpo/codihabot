package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"instabot/ai"
	"instabot/insta"
	"instabot/postgen"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// main.go
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	// if ./posts directory doesnt exist, create
	if _, err := os.Stat("./posts"); os.IsNotExist(err) {
		if err := os.Mkdir("./posts", 0755); err != nil {
			fmt.Println("Error creating posts directory:", err)
			return
		}
	}

	// Parse command line arguments
	var autoUpload bool
	flag.BoolVar(&autoUpload, "y", false, "Automatically upload to Instagram without prompting")
	flag.Parse()

	// Initialize embedded templates filesystem
	postgen.SetTemplatesFS(GetTemplatesFS())

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	post, err := ai.GeneratePost(ctx)

	if err != nil {
		fmt.Println("Error generating post:", err)
		return
	}

	postPath, err := postgen.GeneratePost(post)

	if err != nil {
		fmt.Println("Error generating post:", err)
		return
	}

	images := make([]string, 0)
	images = append(images, postPath+"hero.jpg")
	for i := 0; i < len(post.Slides.Info); i++ {
		images = append(images, postPath+fmt.Sprintf("info_%d.jpg", i))
	}
	images = append(images, postPath+"finish.jpg")

	// run "open {postPath}"
	exec.Command("open", postPath).Run()

	// Check if user wants to upload to Instagram
	shouldUpload := autoUpload
	if !autoUpload {
		shouldUpload = PromptForUpload()
	}

	if shouldUpload {
		fmt.Println("Uploading to Instagram...")

		username := os.Getenv("INSTAGRAM_USERNAME")
		password := os.Getenv("INSTAGRAM_PASSWORD")

		if username == "" || password == "" {
			fmt.Println("Instagram credentials are not set.")

			for username == "" {
				fmt.Print("Please enter the username: ")
				fmt.Scanln(&username)
			}
			for password == "" {
				fmt.Printf("%s' password: ", username)
				fmt.Scanln(&password)
			}
		}

		account := insta.NewAccount(username, password)
		defer account.Dispose()
		account.UploadCarousel(images, post.Caption)
		fmt.Println("Upload completed!")
	} else {
		fmt.Println("Skipping Instagram upload.")
	}
}

func PromptForUpload() bool {
	fmt.Print("Do you want to upload this post to Instagram? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// func main() {
// 	client := images.NewClient()
// 	image, err := client.SearchPNGImage("ReactJS logo TRANSPARENT")
// 	if err != nil {
// 		fmt.Println("Error searching for image:", err)
// 		return
// 	}
//
// 	fmt.Println("Found image:", image)
// }
