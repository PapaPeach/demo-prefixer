package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const WISH_DIR = "tf" + string(filepath.Separator) + "custom"
const CFG_DIR = "custom" + string(filepath.Separator) + "demo-prefixer" + string(filepath.Separator) + "cfg"

var mapDirs = [3]string{"maps", "download" + string(filepath.Separator) + "maps", "workshop" + string(filepath.Separator) + "content" + string(filepath.Separator) + "maps"}
var compGamemodes = []string{"ad", "bball", "cp", "koth", "pl", "ultiduo", "ultitrio"}
var mapSuffixes = [5]string{"_a", "_b", "_f", "_rc", "_v"}

func main() {
	// Check program location
	currentDir := CheckDir()
	tfDir := currentDir[:strings.LastIndex(currentDir, string(filepath.Separator))]
	os.Chdir(tfDir)

	if len(os.Args) > 1 {
		fmt.Println("Commandline arguments:", os.Args[1:])
	}

	var hasCompOnly bool
	var hasNoSuffix bool
	var compOnly bool
	var noSuffix bool
	var madeDir bool

	// Iterate through different map directories
	for _, dir := range mapDirs {
		fmt.Printf("Searching tf/%v directory for maps...", dir)

		var filteredMaps []string
		mapsHandle, err := os.Open(dir)
		switch {
		// If file doesn't exist, skip it
		case errors.Is(err, os.ErrNotExist):
			continue
		// Unexpected error
		case err != nil:
			fmt.Println(err)
			panic("Failed to convert maps path to handle")
		}
		defer mapsHandle.Close()

		// Get contents of maps directory
		maps, err := mapsHandle.Readdirnames(0)
		if err != nil {
			panic("Failed to get contents of tf/" + dir)
		}

		// Filter map names
		for _, m := range maps {
			// Remove file extension
			if !strings.Contains(m, ".") {
				continue
			}
			mapName := m[:strings.Index(m, ".")]

			// Only add maps
			if strings.HasSuffix(m, ".bsp") {
				filteredMaps = append(filteredMaps, mapName)
			} else {
				fmt.Printf("Skipping: %v, not a map.\n", m)
			}
		}

		// Customize map names based on user input / commandline args
		if len(os.Args) > 1 { // Has commandline args
			// CompOnly arg
			if slices.Contains(os.Args[1:], "componly") || slices.Contains(os.Args[1:], "CompOnly") {
				hasCompOnly = true
				compOnly = true
			}
			// NoSuffix arg
			if slices.Contains(os.Args[1:], "nosuffix") || slices.Contains(os.Args[1:], "NoSuffix") {
				hasNoSuffix = true
				noSuffix = true
			}
			// Silent arg
			if slices.Contains(os.Args[1:], "silent") || slices.Contains(os.Args[1:], "Silent") {
				hasCompOnly = true
				hasNoSuffix = true
			}
		}
		for !hasCompOnly { // User input for CompOnly
			fmt.Println("\nWould you like to generate prefixes only for competitive gamemodes (AD, CP, PL, KOTH)? [Y] / [N]")

			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				panic("Error getting user input for CompOnly")
			}
			response = strings.TrimSpace(response)
			switch {
			case strings.EqualFold(response, "y") || strings.EqualFold(response, "yes"):
				hasCompOnly = true
				compOnly = true
				fmt.Println("Generating map prefixes for competitive gamemodes only")
			case strings.EqualFold(response, "n") || strings.EqualFold(response, "no"):
				hasCompOnly = true
				compOnly = false
				fmt.Println("Generating map prefixes for all gamemodes")
			}
		}
		for !hasNoSuffix { // User input for NoSuffix
			fmt.Println("\nWould you like to remove suffixes from map name prefixes? [Y] / [N]\nExample: cp_process_f12 -> cp_process")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				panic("Error getting user input for NoSuffix")
			}
			response = strings.TrimSpace(response)
			switch {
			case strings.EqualFold(response, "y") || strings.EqualFold(response, "yes"):
				hasNoSuffix = true
				noSuffix = true
				fmt.Println("Generating map prefixes without suffixes")
			case strings.EqualFold(response, "n") || strings.EqualFold(response, "no"):
				hasNoSuffix = true
				noSuffix = false
				fmt.Println("Generating map prefixes with suffixes")
			}
		}

		// Only apply prefix customizations if relevant
		var mapNames []string
		if compOnly || noSuffix {
			filteredMaps, mapNames = CustomizeMaps(filteredMaps, compOnly, noSuffix)
		} else {
			mapNames = filteredMaps
		}

		// Generate tf/custom/demo-prefixer/cfg if not done yet
		if !madeDir {
			MakeCfgDir()
			madeDir = true
		}

		WriteCfgs(filteredMaps, mapNames)

		fmt.Println()
	}
}

func CheckDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error checking working directory")
	}

	if !strings.HasSuffix(currentDir, WISH_DIR) {
		fmt.Printf("Program is currently in: \"%v\"\nPlease put program into your \"tf/custom\" folder.\n", currentDir)
		PressToExit()
	}
	return currentDir
}

func CustomizeMaps(maps []string, compOnly bool, noSuffix bool) ([]string, []string) {
	var filteredMaps []string
	if compOnly {
		for _, m := range maps {
			if !strings.Contains(m, "_") {
				continue
			}
			prefix := m[:strings.Index(m, "_")]
			if slices.Contains(compGamemodes, prefix) { // Add if comp gamemode
				filteredMaps = append(filteredMaps, m)
			}
		}
	} else {
		filteredMaps = maps
	}

	var mapNames []string
	if noSuffix {
		for _, m := range filteredMaps {
			var suffixFound bool
			if strings.Count(m, "_") > 1 { // Check if prefix exists
				suffix := m[strings.LastIndex(m, "_"):]
				for _, s := range mapSuffixes { // If suffix is one that should be removed
					if strings.HasPrefix(suffix, s) {
						name := m[:strings.LastIndex(m, "_")]
						suffixFound = true
						mapNames = append(mapNames, name)
						break
					}
				}
			}
			if !suffixFound { // Suffix not removed
				mapNames = append(mapNames, m)
			}
		}
	} else { // Set map names to their file name
		mapNames = filteredMaps
	}
	return filteredMaps, mapNames
}

func MakeCfgDir() {
	// Make tf/custom/demo-prefixer/cfg
	err := os.MkdirAll(CFG_DIR, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic("Failed to make tf/custom/demo-prefixer/cfg")
	}
}

func WriteCfgs(filteredMaps []string, mapNames []string) {
	fmt.Printf(
		"Starting generation with \"%v\", ending with \"%v\"\n",
		filteredMaps[0],
		filteredMaps[len(filteredMaps)-1],
	)
	for i, m := range filteredMaps {
		// Create file
		cfgFile := fmt.Sprintf("%v%v%v.cfg", CFG_DIR, string(filepath.Separator), m)
		cfgHandle, err := os.Create(cfgFile)
		if err != nil {
			panic("Failed to get cfgHandle")
		}
		defer cfgHandle.Close()

		// Generate contents
		cfgContents := "ds_prefix " + mapNames[i]

		// Write file
		_, err = cfgHandle.WriteString(cfgContents)
		if err != nil {
			panic("Error writing cfg contents to file: " + cfgFile)
		}
	}

	fmt.Printf("Generated %v config files", len(filteredMaps))
}

func PressToExit() {
	fmt.Println("\nPress Enter to exit")
	_, err := fmt.Scanln()
	if err != nil {
		// Don't care
		os.Exit(0)
	}
	os.Exit(0)
}
