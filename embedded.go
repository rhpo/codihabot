package main

import (
	"embed"
	"io/fs"
)

// Embed all templates and resources into the binary
//
//go:embed templates/*
var templatesFS embed.FS

//go:embed resources/*
var resourcesFS embed.FS

// GetTemplatesFS returns the embedded templates filesystem
func GetTemplatesFS() fs.FS {
	return templatesFS
}

// GetResourcesFS returns the embedded resources filesystem
func GetResourcesFS() fs.FS {
	return resourcesFS
}
