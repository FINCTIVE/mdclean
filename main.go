package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var unusedImages = make([]string, 0)

	path, _ := os.Getwd()
	fmt.Println("working directory: ", path)
	imgRe := regexp.MustCompile(`.*\.(?i:jpg|gif|png|bmp|svg)`)
	// check each image files
	err := filepath.Walk(path, func(imagePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if imgRe.Match([]byte(imagePath)) {
			fmt.Println("\n\n===> image: ", imagePath)
			// find all markdown files in the same folder
			dirPath := filepath.Dir(imagePath)
			files, err := ioutil.ReadDir(dirPath)
			if err != nil {
				log.Fatal(err)
			}

			usedByMarkdown := false
			for _, f := range files {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
					fmt.Print("markdown file: ", f.Name())
					bytes, err := ioutil.ReadFile(filepath.Join(dirPath, f.Name()))
					if err != nil {
						log.Fatal(err)
					}
					if findImage(imagePath, dirPath, string(bytes)) {
						fmt.Print("  (y)\n")
						usedByMarkdown = true
						break
					} else {
						fmt.Print("  (n)\n")
					}
				}
			}

			if !usedByMarkdown {
				unusedImages = append(unusedImages, imagePath)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\n\nResults(images which are not used by markdown files):")
	for _, i := range unusedImages {
		fmt.Println(i)
	}
}

// findImage checks whether the markdown file uses the image.
func findImage(imagePath string, dirPath string, mdContent string) (exist bool) {
	mdImageRe := regexp.MustCompile(`!\[[^]]*]\((?P<image>.*)\)`)
	results := mdImageRe.FindAllStringSubmatch(mdContent, -1)

	if len(results) == 0 {
		return false
	}

	for _, subMatch := range results {
		// TODO: check whether the method is correct
		mdImage := filepath.Join(dirPath, subMatch[1])
		imagePath, _ := filepath.Abs(imagePath)
		//fmt.Println(mdImage, imagePath)
		if mdImage == imagePath {
			return true
		}
	}
	return false
}
