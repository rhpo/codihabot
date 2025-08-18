package screenshot

import (
	"bytes"
	"context"
	"image/jpeg"
	"image/png"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func HTMLToJPEG(html string, output string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	dataURL := "data:text/html," + url.PathEscape(html)

	var pngBuf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.FullScreenshot(&pngBuf, 100),
	)
	if err != nil {
		return err
	}

	// Decode PNG
	img, err := png.Decode(bytes.NewReader(pngBuf))
	if err != nil {
		return err
	}

	// Save as JPEG
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
}
