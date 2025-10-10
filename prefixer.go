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
	// Get working directory
	workingDir, wderr := os.Getwd()
	// workingDir += "/test1 copy"
	if wderr != nil {
		fmt.Println("Error getting working directory:", wderr)
		EnterToContinue()
	}

	prefix(workingDir)
	EnterToContinue()
}

func prefix(dirString string) {
	// Open current directory
	dir, oerr := os.Open(dirString)
	if oerr != nil {
		fmt.Println("Error opening directory:", oerr)
		EnterToContinue()
	}
	defer dir.Close()

	// Get list of files in current directory
	files, rerr := dir.Readdirnames(0)
	if rerr != nil {
		fmt.Println("Error getting files in current directory:", rerr)
		EnterToContinue()
	}

	// Open _events.txt
	eventsFile, eferr := os.Open(dirString + string(filepath.Separator) + "_events.txt")
	if eferr != nil {
		fmt.Println("Error reading _events.txt")
	}
	eventScanner := bufio.NewScanner(eventsFile)

	// Get all valid event lines
	hasEvents := false
	var events []string
	for eventScanner.Scan() {
		line := eventScanner.Text()
		events = append(events, line)

		// Only update for events with valid formatting
		if !hasEvents && strings.Contains(line, "(\"") {
			hasEvents = true
		}
	}

	eventsFile.Close()

	// For each file in directory...
	mapCount := 0
	eventCount := 0
	eventsIndex := 0 //TODO: Add "thorough" argument to skip this
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
		isDemo := false
		switch extension {
		case ".dem": // Demo file
			// Open file
			file, oferr := os.Open(filePath)
			if oferr != nil {
				fmt.Printf("Error opening %s: %s\n", filePath, oferr)
				EnterToContinue()
			}

			// Read file and extract map name
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			rawText := scanner.Text()[mapStart : mapStart+mapLength]
			for i := range mapLength {
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
			isDemo = true

		case ".json": // json file
			if fileName != prevFileName {
				continue
			}

		case ".tga": // tga file
			if fileName != prevFileName {
				continue
			}

		default: // other file we don't care about
			continue
		}

		// If file name does not contain the map name, then replace existing prefix
		if !strings.HasPrefix(fileFullName, mapName) {
			dateIndex := strings.Index(fileFullName, "-") - 4 // prefixYYYY-MM-DD_HH-MM-SS
			oldFileName := fileFullName
			if dateIndex > 0 && fileFullName[dateIndex] == '2' { // Validate date
				fileFullName = fileFullName[dateIndex:] // Remove prefix
			}
			newPath := dirString + string(filepath.Separator) + mapName + "_" + fileFullName

			// Rename
			rnerr := os.Rename(filePath, newPath)
			if rnerr != nil {
				fmt.Printf("Error opening %s: %s\n", filePath, rnerr)
				EnterToContinue()
			}

			// Update any associated event entries
			if isDemo && hasEvents {
				isMatch := false
				for i, event := range events {
					if i >= eventsIndex && event[0] != '>' { // Skip already checked and marker lines
						eventStart := strings.Index(event, "(\"") + 2
						if eventStart > -1 && strings.HasPrefix(event[eventStart:], fileName) { // If event matches demo title
							newEvent := (event[:eventStart] +
								strings.Replace(event[eventStart:], fileName, (mapName+"_"+fileFullName[:len(fileFullName)-4]), 1))
							events[i] = newEvent
							eventsIndex = i + 1
							isMatch = true
							eventCount++
							continue // Continue loop as long as we are matching
						}

						// Stop loop once we are no longer matching (if for some reason reason there's no marker line)
						if isMatch {
							fmt.Println("Marker missing at _events.txt line:", i+1)
							eventsIndex = i + 1
							break
						}
					} else if isMatch { // Stop loop once we are no longer matching
						eventsIndex = i + 1
						break
					}
				}

				// If no match was found after our starting point, circle back to start
				if !isMatch {
					// From start to where the above loop started
					for i, event := range events[:eventsIndex] {
						if event[0] != '>' { // Skip marker lines
							eventStart := strings.Index(event, "(\"") + 2
							if eventStart > -1 && strings.HasPrefix(event[eventStart:], fileName) { // If event matches demo title
								newEvent := (event[:eventStart] +
									strings.Replace(event[eventStart:], fileName, (mapName+"_"+fileFullName[:len(fileFullName)-4]), 1))
								events[i] = newEvent
								eventsIndex = i + 1
								isMatch = true
								eventCount++
								continue // Continue loop as long as we are matching
							}

							// Stop loop once we are no longer matching (if for some reason reason there's no marker line)
							if isMatch {
								fmt.Println("Marker missing at _events.txt line:", i+1)
								eventsIndex = i + 1
								break
							}
						} else if isMatch { // Stop loop once we are no longer matching
							eventsIndex = i + 1
							break
						}
					}
				}
			}

			// Print rename
			fmt.Printf("Renamed: %s \tTo: %s\n", oldFileName, newPath[strings.LastIndex(newPath, string(filepath.Separator))+1:])
			mapCount++
		}
	}

	// Update _events.txt with our changes
	if eventCount > 0 {
		// Open _events.txt for writing
		eventsWriter, ewerr := os.OpenFile((dirString + string(filepath.Separator) + "_events.txt"), os.O_CREATE, 0644)
		if ewerr != nil {
			fmt.Printf("Error writing _events.txt\n")
			EnterToContinue()
		}

		// Write our updated events
		for _, event := range events {
			_, wserr := eventsWriter.WriteString(event + "\n")
			if wserr != nil {
				fmt.Println("Error writing to _events.txt")
				EnterToContinue()
			}
		}

		eventsWriter.Close()
	}

	// Print total number of maps renamed
	fmt.Printf("Renamed %d files.\n", mapCount)
	fmt.Printf("Updated %d events.\n", eventCount)
}

func EnterToContinue() {
	fmt.Println("Press enter to exit...")
	_, _ = fmt.Scanln()
	os.Exit(0)
}
