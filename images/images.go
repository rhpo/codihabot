package images

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	apiURL = "https://www.googleapis.com/customsearch/v1"
)

type Client struct {
	APIKey         string
	SearchEngineID string
}

type searchResult struct {
	Items []struct {
		Link string `json:"link"` // direct image URL
	} `json:"items"`
}

var apiKey = os.Getenv("SEARCH_ENGINE_API_KEY")
var searchEngineID = os.Getenv("SEARCH_ENGINE_ID")

// NewClient creates a new Client with the specified API key and search engine ID.
func NewClient() *Client {
	return &Client{
		APIKey:         apiKey,
		SearchEngineID: searchEngineID,
	}
}

// SearchPNGImage searches for PNG images matching the query,
// downloads the first PNG image found, makes white background transparent if needed,
// and returns a base64 data URI string.
func (c *Client) SearchPNGImage(query string) (string, error) {
	params := url.Values{}
	params.Set("key", c.APIKey)
	params.Set("cx", c.SearchEngineID)
	params.Set("q", query)
	params.Set("searchType", "image")
	params.Set("fileType", "png")
	params.Set("num", "10")

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("failed to fetch images: status code " + resp.Status)
	}

	var result searchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, item := range result.Items {
		// Parse the link and check extension robustly
		u, err := url.Parse(item.Link)
		if err != nil {
			continue
		}
		if strings.HasSuffix(strings.ToLower(u.Path), ".png") {
			// Download and fix transparency
			dataURI, err := downloadFixAndEncodeBase64(item.Link)
			if err != nil {
				continue // try next image
			}
			return dataURI, nil
		}
	}

	return "", errors.New("no valid PNG image found")
}

func downloadFixAndEncodeBase64(imgURL string) (string, error) {
	fmt.Println("Downloading image:", imgURL)

	resp, err := http.Get(imgURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", err
	}

	if !hasTransparency(img) {
		img = makeWhiteTransparent(img)
	}

	// Automatically crop transparent areas to fit the icon perfectly
	img = autoCropTransparent(img)

	// Re-encode image to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	imgData = buf.Bytes()

	// base64 encode
	b64 := base64.StdEncoding.EncodeToString(imgData)
	dataURI := "data:image/png;base64," + b64
	return dataURI, nil
}

// hasTransparency checks if the given image has any transparent pixels.
// It returns true if at least one pixel is transparent, otherwise false.
func hasTransparency(img image.Image) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a < 65535 {
				return true
			}
		}
	}
	return false
}

// makeWhiteTransparent removes the white background from an image using the remove.bg API.
//
// The function first attempts to encode the provided image to PNG format. If encoding fails, it falls back to a manual method for converting white to transparent. It then constructs a multipart HTTP request to the remove.bg API, including the image data and necessary headers. If the API call is successful, it reads the response and decodes the processed image. In case of any errors during the process, it defaults to the manual conversion method.
func makeWhiteTransparent(img image.Image) image.Image {
	// First encode the image to bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		// If encoding fails, fallback to manual white->transparent conversion
		return makeWhiteTransparentManual(img)
	}

	// Call remove.bg API to remove background
	client := &http.Client{}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add API key header
	part, err := writer.CreateFormFile("image_file", "image.png")
	if err != nil {
		return makeWhiteTransparentManual(img)
	}

	if _, err := part.Write(buf.Bytes()); err != nil {
		return makeWhiteTransparentManual(img)
	}

	// Add size parameter
	if err := writer.WriteField("size", "auto"); err != nil {
		return makeWhiteTransparentManual(img)
	}

	writer.Close()

	req, err := http.NewRequest("POST", "https://api.remove.bg/v1.0/removebg", &body)
	if err != nil {
		return makeWhiteTransparentManual(img)
	}

	req.Header.Set("X-Api-Key", "ufqrYoMAwS3hXtcnx4db2wPX")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return makeWhiteTransparentManual(img)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// If API call fails, fallback to manual method
		return makeWhiteTransparentManual(img)
	}

	// Read the response image
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return makeWhiteTransparentManual(img)
	}

	// Decode the processed image
	processedImg, err := png.Decode(bytes.NewReader(responseData))
	if err != nil {
		return makeWhiteTransparentManual(img)
	}

	return processedImg
}

// Fallback function for manual white->transparent conversion
func makeWhiteTransparentManual(img image.Image) *image.NRGBA {
	b := img.Bounds()
	newImg := image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b_, a := img.At(x, y).RGBA()
			rr := uint8(r >> 8)
			gg := uint8(g >> 8)
			bb := uint8(b_ >> 8)
			aa := uint8(a >> 8)

			// Make white pixels (and near-white pixels) transparent
			if rr > 240 && gg > 240 && bb > 240 {
				newImg.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 0})
			} else {
				newImg.SetNRGBA(x, y, color.NRGBA{rr, gg, bb, aa})
			}
		}
	}
	return newImg
}

// autoCropTransparent automatically crops transparent areas around the image to fit the icon content perfectly into frame.
//
// It scans the image to determine the bounds of non-transparent content by checking the alpha value of each pixel. If non-transparent content is found, it calculates the minimum and maximum coordinates and applies a padding based on the content size. Finally, it creates and returns a new cropped image that retains the relevant content while removing excess transparent areas. If no content is found, the original image is returned.
func autoCropTransparent(img image.Image) image.Image {
	bounds := img.Bounds()

	// Find the bounds of non-transparent content
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	hasContent := false
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			// Consider pixel as non-transparent if alpha > 25% (16383 out of 65535)
			if a > 16383 {
				hasContent = true
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// If no content found, return original image
	if !hasContent {
		return img
	}

	// Add a small padding (5% of the content size or minimum 2 pixels)
	contentWidth := maxX - minX + 1
	contentHeight := maxY - minY + 1
	paddingX := max(2, contentWidth/20)
	paddingY := max(2, contentHeight/20)

	// Apply padding but keep within original bounds
	minX = max(bounds.Min.X, minX-paddingX)
	minY = max(bounds.Min.Y, minY-paddingY)
	maxX = min(bounds.Max.X-1, maxX+paddingX)
	maxY = min(bounds.Max.Y-1, maxY+paddingY)

	// Create cropped image
	croppedWidth := maxX - minX + 1
	croppedHeight := maxY - minY + 1
	cropped := image.NewNRGBA(image.Rect(0, 0, croppedWidth, croppedHeight))

	for y := 0; y < croppedHeight; y++ {
		for x := 0; x < croppedWidth; x++ {
			srcX := minX + x
			srcY := minY + y
			cropped.Set(x, y, img.At(srcX, srcY))
		}
	}

	return cropped
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
