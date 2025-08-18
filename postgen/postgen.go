package postgen

import (
	"fmt"
	"instabot/postgen/screenshot"
	"instabot/utils"
	"io/fs"
	"os"
	"strings"

	"html/template"
)

var templatesFS fs.FS
var postsPath = "./posts/"

// SetTemplatesFS sets the embedded templates filesystem
func SetTemplatesFS(fsys fs.FS) {
	templatesFS = fsys
}

func renderTemplate(templateFileName string, data any, funcMap template.FuncMap) (string, error) {
	// Read template from embedded filesystem
	templateData, err := fs.ReadFile(templatesFS, "templates/"+templateFileName)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templateFileName, err)
	}

	tmpl, err := template.New(templateFileName).Funcs(funcMap).Parse(string(templateData))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateFileName, err)
	}

	var rendered strings.Builder
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return rendered.String(), nil
}

func GenerateHero(post Post, output string) (string, error) {
	html, err := renderTemplate("hero.html", post, template.FuncMap{
		"lighten": utils.Lighten,
		"html":    utils.RenderHTML,
		"safeURL": func(url string) template.URL { return template.URL(url) },
	})

	if err != nil {
		fmt.Println("Error generating hero slide:", err)
		return "", err
	}

	os.WriteFile("debug.html", []byte(html), 0644)

	err = screenshot.HTMLToJPEG(html, output)
	if err != nil {
		return "", err
	}
	println("Generated hero slide PNG:", output)

	return output, nil
}

func GenerateFinish(post Post, output string) (string, error) {
	html, err := renderTemplate("finish.html", post, template.FuncMap{
		"lighten": utils.Lighten,
		"html":    utils.RenderHTML,
		"safeURL": func(url string) template.URL { return template.URL(url) },
	})
	if err != nil {
		fmt.Println("Error generating finish slide:", err)
		return "", err
	}

	err = screenshot.HTMLToJPEG(html, output)
	if err != nil {
		return "", err
	}
	println("Generated finish slide PNG:", output)

	return output, nil
}

func GenerateInfo(info InfoSlide, color string, output string) (string, error) {
	data := struct {
		Slide InfoSlide
		Color string
	}{
		Slide: info,
		Color: color,
	}

	html, err := renderTemplate("info.html", data, template.FuncMap{
		"lighten": utils.Lighten,
		"html":    utils.RenderHTML,
		"safeURL": func(url string) template.URL { return template.URL(url) },
	})
	if err != nil {
		fmt.Println("Error generating info slide:", err)
		return "", err
	}

	err = screenshot.HTMLToJPEG(html, output)
	if err != nil {
		return "", err
	}
	println("Generated info slide PNG:", output)

	return output, nil
}

func GeneratePost(post Post) (string, error) {
	postPath := postsPath + post.Name + "/"

	// check if post directory exists
	if err := os.MkdirAll(postPath, 0755); err != nil {
		return "", err
	}

	GenerateHero(post, postPath+"hero.jpg")

	for i, info := range post.Slides.Info {
		if _, err := GenerateInfo(info, post.Color, postPath+"info_"+fmt.Sprint(i)+".jpg"); err != nil {
			return "", err
		}
	}

	GenerateFinish(post, postPath+"finish.jpg")

	return postPath, nil
}
