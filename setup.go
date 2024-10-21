package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	inertia "github.com/romsar/gonertia"
)

func InitInertia() *inertia.Inertia {
	viteHotFile := "./public/hot"
	rootViewFile := "resources/views/root.html"

	_, err := os.Stat(viteHotFile)
	if err == nil {
		i, err := inertia.NewFromFile(
			rootViewFile,
			// inertia.WithSSR(),
		)
		if err != nil {
			log.Fatal(err)
		}
		i.ShareTemplateFunc("vite", func(entry string) (string, error) {
			content, err := os.ReadFile(viteHotFile)
			if err != nil {
				return "", err
			}
			url := strings.TrimSpace(string(content))
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				url = url[strings.Index(url, ":")+1:]
			} else {
				url = "//localhost:8080"
			}
			if entry != "" && !strings.HasPrefix(entry, "/") {
				entry = "/" + entry
			}
			return url + entry, nil
		})
		return i
	}

	manifestPath := "./public/build/manifest.json"
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		err := os.Rename("./public/build/.vite/manifest.json", "./public/build/manifest.json")
		if err != nil {
			return nil
		}
	}

	i, err := inertia.NewFromFile(
		rootViewFile,
		inertia.WithVersionFromFile(manifestPath),
		// inertia.WithSSR(),
	)
	if err != nil {
		log.Fatal(err)
	}

	i.ShareTemplateFunc("vite", Vite(manifestPath, "/build/"))

	return i
}

func Vite(manifestPath, buildDir string) func(path string) (string, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		log.Fatalf("cannot open provided vite manifest file: %s", err)
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})
	err = json.NewDecoder(f).Decode(&viteAssets)

	for k, v := range viteAssets {
		log.Printf("%s: %s\n", k, v.File)
	}

	if err != nil {
		log.Fatalf("cannot unmarshal vite manifest file to json: %s", err)
	}

	return func(p string) (string, error) {
		if val, ok := viteAssets[p]; ok {
			return path.Join("/", buildDir, val.File), nil
		}
		return "", fmt.Errorf("asset %q not found", p)
	}
}
