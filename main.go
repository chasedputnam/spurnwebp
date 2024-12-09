package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/image/webp"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Removed path from a filepath. Leaves dir path as-is.
// Returns "" for relative path values or for root dir
func getPathFromFilePath(filePath string) string {
	path := filepath.Dir(filePath)
	if strings.HasPrefix(path, ".") || path == "/" {
		return ""
	}
	return path
}

// Takes input webp file and converts it to JPG or PNG
func convertWebp(input string, outputType string) {
	path := getPathFromFilePath(input)
	output := filepath.Base(input)
	if !strings.HasSuffix(output, ".webp") {
		fmt.Println("Error: Input file is not a webp file type.")
		os.Exit(1)
	}
	output = output[0 : len(output)-len(filepath.Ext(output))] // Remove file extension
	f, err := os.Open(input)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer f.Close()

	img, err := webp.Decode(f)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if outputType == "PNG" {
		// Convert to PNG
		pngFile, err := os.Create(path + "/" + output + ".png")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer func(pngFile *os.File) {
			err := pngFile.Close()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}(pngFile)

		err = png.Encode(pngFile, img)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	} else {
		// Convert to JPEG
		jpegFile, err := os.Create(path + "/" + output + ".jpg")
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer func(jpegFile *os.File) {
			err := jpegFile.Close()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}(jpegFile)

		err = jpeg.Encode(jpegFile, img, &jpeg.Options{Quality: 90})
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
}

func main() {
	var input string
	flag.StringVar(&input, "input", "", "path to input WebP file")
	var watch string
	flag.StringVar(&watch, "watch", "", "path to watch for WebP files")
	var outputType string
	flag.StringVar(&outputType, "outputType", "", "image type to output: PNG (Default) or JPG")
	flag.Parse()

	if input != "" && watch != "" { // If both input parameters are used, exit
		fmt.Println("Error: You cannot specify both -watch and -input")
		os.Exit(1)
	} else if input == "" && watch == "" { // If neither input parameters are used, exit
		fmt.Println("Error: You need to specify one of -watch or -input")
		os.Exit(1)
	}

	if outputType == "" { // Set to PNG by default if no outputType is provided
		outputType = "PNG"
	}

	if input != "" && watch == "" { // Direct conversion of a file
		convertWebp(input, outputType)
		fmt.Println("Done: Conversion completed successfully")
		os.Exit(0)
	} else if watch != "" && input == "" { // Watching a directory for webp files to convert
		// creates a new file watcher
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Println("ERROR", err)
		}
		// Ensure watch path is a directory
		path := getPathFromFilePath(watch)
		if path != "" {
			err = watcher.Add(watch)
			if err != nil {
				fmt.Println("ERROR", err)
				os.Exit(1)
			}
		}
		defer watcher.Close()

		// Set done to a channel to run goroutine through
		done := make(chan bool)

		// goroutine for watching file system location
		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					//fmt.Printf("EVENT! %#v\n", event)
					if event.Op.Has(fsnotify.Create) && strings.HasSuffix(event.Name, ".webp") {
						convertWebp(event.Name, outputType)
					}
					// watch for errors
				case err := <-watcher.Errors:
					fmt.Println("ERROR", err)
				}
			}
		}()

		// out of the box fsnotify can watch a single file, or a single directory
		if err := watcher.Add(watch); err != nil {
			fmt.Println("ERROR", err)
		}
		<-done

	} else {
		fmt.Println("Error: No valid input parameters and values specified")
		os.Exit(1)
	}
}
