package postgen

import (
	"fmt"
	"image/color"
	"instabot/postgen/markdown"
)

type HeroSlide struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type InfoSlide struct {
	Title   string `json:"title"`
	Content string `json:"content"` // It's a markdown string that should be converted before injection to the post
}

type FinishSlide struct {
	Description string `json:"description"`
	Image       string `json:"image"`
}

type Slides struct {
	Hero   HeroSlide   `json:"hero"`
	Info   []InfoSlide `json:"info"`
	Finish FinishSlide `json:"finish"`
}

type Post struct {
	Color string `json:"color"`

	Name    string `json:"name"`
	Caption string `json:"caption"`
	Slides  Slides `json:"slides"`
}

func ColorToCssColor(c color.Color, myName string) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("rgba(%d, %d, %d, %f)", r, g, b, float64(a)/255)
}

func NewPost(post Post) Post {

	post.Slides.Hero.Title = markdown.MarkdownToHTML(post.Slides.Hero.Title)
	post.Slides.Hero.Description = markdown.MarkdownToHTML(post.Slides.Hero.Description)

	for i := range post.Slides.Info {
		post.Slides.Info[i].Title = markdown.MarkdownToHTML(post.Slides.Info[i].Title)
		post.Slides.Info[i].Content = markdown.MarkdownToHTML(post.Slides.Info[i].Content)
	}

	post.Slides.Finish.Description = markdown.MarkdownToHTML(post.Slides.Finish.Description)
	return post
}
