package utils

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
)

func Lighten(color string) string {
	// Remove # if present
	color = strings.TrimPrefix(color, "#")

	// Parse RGB values
	if len(color) != 6 {
		return color // Return original if invalid format
	}

	r, err1 := strconv.ParseUint(color[0:2], 16, 8)
	g, err2 := strconv.ParseUint(color[2:4], 16, 8)
	b, err3 := strconv.ParseUint(color[4:6], 16, 8)

	if err1 != nil || err2 != nil || err3 != nil {
		return color // Return original if parsing fails
	}

	var num uint64 = 60 // Amount to lighten each component

	// Lighten by adding 30 to each component (capped at 255)
	r = min(255, r+num)
	g = min(255, g+num)
	b = min(255, b+num)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func RenderHTML(content string) template.HTML {
	return template.HTML(content)
}

func GetAbsolutePath(relativePath string) string {
	absPath, _ := filepath.Abs(relativePath)

	// Convert to file URL
	return "file://" + filepath.ToSlash(absPath)
}
