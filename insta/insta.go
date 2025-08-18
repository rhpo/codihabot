package insta

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Davincible/goinsta/v3"
)

type Account struct {
	insta *goinsta.Instagram
}

func NewAccount(username, password string) *Account {
	insta := goinsta.New(username, password)
	insta.Login(password)
	return &Account{insta: insta}
}

func (a *Account) Dispose() {
	a.insta.Export("~/.goinsta")
}

func (a *Account) UploadMedia(mediaPath string, caption string) {
	// Upload media item
	file, err := os.Open(mediaPath)
	if err != nil {
		log.Fatalf("Failed to open media file %s: %v", mediaPath, err)
	}

	defer file.Close()

	item, err := a.insta.Upload(&goinsta.UploadOptions{
		File:    file,
		Caption: caption,
	})

	if err != nil {
		log.Fatalf("Failed to upload media %s: %v", mediaPath, err)
	}

	fmt.Printf("Media uploaded successfully: %s\n", item.ID)
}

func (a *Account) UploadCarousel(mediaPaths []string, caption string) {
	// Upload media items as a carousel
	var mediaFiles []*os.File
	for _, mediaPath := range mediaPaths {
		file, err := os.Open(mediaPath)
		if err != nil {
			log.Fatalf("Failed to open media file %s: %v", mediaPath, err)
		}

		mediaFiles = append(mediaFiles, file)
	}
	// Convert []*os.File to []io.Reader
	var readers []io.Reader
	for _, f := range mediaFiles {
		readers = append(readers, f)
	}
	// Defer closing files after upload
	defer func() {
		for _, f := range mediaFiles {
			f.Close()
		}
	}()

	item, err := a.insta.Upload(&goinsta.UploadOptions{
		Album:   readers,
		Caption: caption,
	})

	if err != nil {
		log.Fatalf("Failed to upload media %v: %v", mediaPaths, err)
	}

	fmt.Printf("Media uploaded successfully: %s\n", item.ID)
}

func (a *Account) DeleteMedia(postID string) {

	// Fetch the media item
	item, err := a.insta.GetMedia(postID)
	if err != nil {
		log.Fatalf("Failed to get media %s: %v", postID, err)
	}

	// Delete it
	if err := item.Delete(); err != nil {
		log.Fatalf("Failed to delete post with ID %s: %v", postID, err)
	} else {
		fmt.Printf("Post with ID %s deleted successfully.\n", postID)
	}

}
