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

var mapSuffixes = [5]string{"_a", "_b", "_f", "_rc", "_v"}

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
	for _, fileName := range files {

		// Only go through .dem files
		if !strings.HasSuffix(fileName, ".dem") {
			continue
		}

		// Open file
		filePath := dirString + string(filepath.Separator) + fileName
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error opening %s: %s", filePath, err)
			EnterToContinue()
		}

		// Read file and extract map name
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		rawText := string(scanner.Text()[mapStart : mapStart+mapLength])
		var mapName string
		for i := 0; i < len(rawText); i++ {
			if !unicode.IsPrint(rune(rawText[i])) {
				mapName = rawText[:i]
				break
			}
		}
		file.Close()

		// Remove map version suffixes
		suffixIndex := strings.LastIndex(mapName, "_")
		if suffixIndex > 3 { // We skip when "suffix" is just the map prefix
			mapTail := mapName[suffixIndex:]
			for _, suffix := range mapSuffixes {
				if strings.HasPrefix(mapTail, suffix) {
					mapName = mapName[:suffixIndex]
				}
			}
		}

		// If file name does not contain the map name, then replace existing index
		if !strings.HasPrefix(fileName, mapName) {
			dateIndex := strings.Index(fileName, "-") - 4 // prefixYYYY-MM-DD_HH-MM-SS
			oldFileName := fileName
			if dateIndex > 0 {
				fileName = fileName[dateIndex:] // Remove prefix
			}
			newPath := dirString + string(filepath.Separator) + mapName + "_" + fileName
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
	fmt.Printf("Renamed %d demos\n", mapCount)
	EnterToContinue()
}

func EnterToContinue() {
	fmt.Println("Press enter to exit...")
	//fmt.Scanln()
	os.Exit(0)
}
