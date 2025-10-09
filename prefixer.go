package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const mapStart = 536
const mapLength = 260

var mapSuffixes = [5]string{"rc", "f", "b", "a", "v"}

func main() {
	// Open current directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		EnterToContinue()
	}
	dirString := workingDir
	dir, err := os.Open(dirString)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		EnterToContinue()
	}
	defer dir.Close()

	// Get list of files in current directory
	files, err := dir.Readdirnames(0)
	if err != nil {
		fmt.Println("Error getting files in current directory:", err)
		EnterToContinue()
	}

	// For each file in directory...
	mapCount := 0
	var prevFileName string
	var mapName string
	for _, fileFullName := range files {

		// Get file name, extension, and path information
		extensionIndex := strings.LastIndex(fileFullName, ".")
		if extensionIndex < 1 {
			continue
		}
		extension := fileFullName[extensionIndex:]
		fileName := fileFullName[:extensionIndex]
		filePath := dirString + string(filepath.Separator) + fileFullName

		// If json or tga, rename if tied to a demo
		// If demo, then grab map name and all that jazz
		switch extension {
		case ".dem": // Demo file
			// Open file
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening %s: %s", filePath, err)
				EnterToContinue()
			}

			// Read file and extract map name
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			rawText := string(scanner.Text()[mapStart : mapStart+mapLength])
			for i := 0; i < len(rawText); i++ {
				if !unicode.IsPrint(rune(rawText[i])) {
					mapName = rawText[:i]
					break
				}
			}
			file.Close()

			// Check if map has as suffix and remove if it is a version number
			prefixIndex := strings.Index(mapName, "_")
			suffixIndex := strings.LastIndex(mapName, "_")
			if prefixIndex > 1 && prefixIndex != suffixIndex { // Skip when "suffix" is just the map prefix
				mapTail := mapName[suffixIndex+1:]
				for _, suffix := range mapSuffixes { // remove if suffix is a version number
					if strings.HasPrefix(mapTail, suffix) {
						mapName = mapName[:suffixIndex]
						break
					}
				}
			}

			// Update previous file name
			prevFileName = fileName

		case ".json": // json file
			if fileName != prevFileName {
				continue
			}

		case ".tga": // tga file
			if fileName != prevFileName {
				continue
			}

		default:
			continue
		}

		// If file name does not contain the map name, then replace existing prefix
		if !strings.HasPrefix(fileFullName, mapName) {
			dateIndex := strings.Index(fileFullName, "-") - 4 // prefixYYYY-MM-DD_HH-MM-SS
			oldFileName := fileFullName
			if dateIndex > 0 {
				fileFullName = fileFullName[dateIndex:] // Remove prefix
			}
			newPath := dirString + string(filepath.Separator) + mapName + "_" + fileFullName
			err := os.Rename(filePath, newPath)
			if err != nil {
				fmt.Printf("Error opening %s: %s", filePath, err)
				EnterToContinue()
			}
			fmt.Printf("Renamed: %s \tTo: %s\n", oldFileName, newPath[strings.LastIndex(newPath, string(filepath.Separator))+1:])
			mapCount++
		}
	}

	// Print total number of maps renamed
	fmt.Printf("Renamed %d files\n", mapCount)
	EnterToContinue()
}

func EnterToContinue() {
	fmt.Println("Press enter to exit...")
	fmt.Scanln()
	os.Exit(0)
}
