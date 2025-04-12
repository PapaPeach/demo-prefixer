package main

import "core:fmt"
import "core:os"
import "core:slice"
import "core:strings"

mapsDir := "maps"
downloadMapsDir := "download/maps"
cfgDir := "custom/demo-prefixer/cfg"

main :: proc() {
	mapsHandle, mhEr := os.open(mapsDir)
	if mhEr != nil {
		panic("Failed to convert maps path to handle")
	}

	// Get contents of maps directory
	maps, mEr := os.read_dir(mapsHandle, 0)
	if mEr != nil {
		panic("Failed to get contents of tf/maps")
	}

	filteredMaps := make([dynamic]string)
	mapName: string
	prevMap := ""
	for i in maps {
		// Remove file extension
		mapName, _ = strings.substring_to(i.name, strings.index(i.name, "."))
		defer prevMap = mapName

		// Remove maps with existing config files or non-map files
		switch {
		case strings.contains(i.name, ".cfg") && strings.equal_fold(mapName, prevMap):
			fmt.printf(
				"Found config file: %v for map: %v\nSkipping this map...\n",
				i.name,
				prevMap,
			)
		case !strings.ends_with(i.name, ".bsp"):
		case:
			append(&filteredMaps, mapName)
		}
	}
	os.close(mapsHandle)

	// Make tf/custom/demo-prefixer/cfg
	mdEr := os.make_directory("custom")
	if mdEr != nil && strings.equal_fold(os.error_string(mdEr), "Exist") {
		panic("Failed to make tf/custom")
	}
	mdEr = os.make_directory("custom/demo-prefixer")
	if mdEr != nil && strings.equal_fold(os.error_string(mdEr), "Exist") {
		panic("Failed to make tf/custom/demo-prefixer")
	}
	mdEr = os.make_directory(cfgDir, os.O_RDWR)
	if mdEr != nil && strings.equal_fold(os.error_string(mdEr), "Exist") {
		panic("Failed to make tf/custom/demo-prefixer/cfg")
	}

	// Make cfg files
	for i in filteredMaps {
		// Create file
		cfgFile := strings.concatenate({cfgDir, "/", i, ".cfg"})
		_ = os.write_entire_file_or_err(cfgFile, nil)
		// Open file
		cfgHandle, cfghEr := os.open(cfgFile, os.O_RDWR)
		if cfghEr != nil {
			panic("Failed to get cfgHandle")
		}
		defer os.close(cfgHandle)

		// Generate contents
		cfgContents := strings.concatenate({"ds_prefix ", i})

		// Write file
		_, wsEr := os.write_string(cfgHandle, cfgContents)
		if wsEr != nil {
			fmt.println(os.get_last_error())
			fmt.panicf("Error writing cfg contents to file: ", cfgFile)
		}
	}
}
