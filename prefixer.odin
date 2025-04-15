package main

import "core:fmt"
import "core:os"
import "core:slice"
import "core:strings"

MAP_DIRS :: [3]string{"maps", "download/maps", "workshop/content/maps"}
CFG_DIR :: "custom/demo-prefixer/cfg"

main :: proc() {
	madeDir: bool
	// Iterate through different map directories
	for dir, i in MAP_DIRS {
		fmt.printfln("Searching tf/%v directory for maps...", dir)

		filteredMaps := make([dynamic]string, context.temp_allocator)
		mapsHandle, mhEr := os.open(dir)
		switch {
		// If file doesn't exits, skip it
		case strings.equal_fold(os.error_string(mhEr), "file does not exist"):
			continue
		// Unexpected error
		case mhEr != nil:
			fmt.println(os.error_string(mhEr))
			panic("Failed to convert maps path to handle")
		}
		defer os.close(mapsHandle)

		// Get contents of maps directory
		maps, mEr := os.read_dir(mapsHandle, 0, context.temp_allocator)
		if mEr != nil {
			panic("Failed to get contents of tf/maps")
		}

		// Filter map names
		// TODO: maps folders shouldn't contain .cfgs and only .bsps should be loaded
		for m, i in maps {
			// Remove file extension
			mapName, _ := strings.substring_to(m.name, strings.index(m.name, "."))

			// Remove maps with existing config files or non-map files
			switch {
			// Skip
			case i >= len(maps) - 1:
			// Remove map from filtered if wierd naming lets it slip in
			case i > 0 &&
			     strings.ends_with(maps[i + 1].name, ".cfg") &&
			     strings.contains(filteredMaps[len(filteredMaps) - 1], mapName):
				fmt.printfln(
					"Found config file: %v for map: %v\nSkipping this map...",
					maps[i + 1].name,
					m.name,
				)
				unordered_remove(&filteredMaps, len(filteredMaps) - 1)
			// file is a config for an existing map
			case strings.ends_with(maps[i + 1].name, ".cfg") &&
			     strings.contains(maps[i + 1].name, mapName):
				fmt.printfln(
					"Found config file: %v for map: %v\nSkipping this map...",
					maps[i + 1].name,
					m.name,
				)
			// Map should be added
			case !strings.ends_with(m.name, ".bsp"):
				fmt.printfln("Skipping: %v, not a map.", m.name)
			case:
				append(&filteredMaps, mapName)
			}
		}

		// Generate tf/custom/demo-prefixer/cfg if not done yet
		if !madeDir {
			make_cfg_dir()
			madeDir = true
		}

		write_cfgs(filteredMaps)

		free_all(context.temp_allocator)

		fmt.println()
	}

	press_to_exit()
}

make_cfg_dir :: proc() {
	// Make tf/custom
	mdEr := os.make_directory("custom")
	if mdEr != nil && strings.equal_fold(os.error_string(mdEr), "Exist") {
		panic("Failed to make tf/custom")
	}

	// Make tf/custom/demo-prefixer
	mdEr = os.make_directory("custom/demo-prefixer")
	if mdEr != nil && strings.equal_fold(os.error_string(mdEr), "Exist") {
		panic("Failed to make tf/custom/demo-prefixer")
	}

	// Make tf/custom/demo-prefixer/cfg
	mdEr = os.make_directory(CFG_DIR, os.O_RDWR)
	if mdEr != nil && strings.equal_fold(os.error_string(mdEr), "Exist") {
		panic("Failed to make tf/custom/demo-prefixer/cfg")
	}
}

write_cfgs :: proc(filteredMaps: [dynamic]string) {
	// Make cfg files
	fmt.printfln(
		"Starting generation with \"%v\", ending with \"%v\"",
		filteredMaps[0],
		filteredMaps[len(filteredMaps) - 1],
	)
	for m, i in filteredMaps {
		// Create file
		cfgFile := fmt.tprintf("%v/%v.cfg", CFG_DIR, m)
		_ = os.write_entire_file_or_err(cfgFile, nil)
		// Open file
		cfgHandle, cfghEr := os.open(cfgFile, os.O_RDWR)
		if cfghEr != nil {
			panic("Failed to get cfgHandle")
		}
		defer os.close(cfgHandle)

		// Generate contents
		cfgContents := strings.concatenate({"ds_prefix ", m})

		// Write file
		_, wsEr := os.write_string(cfgHandle, cfgContents)
		if wsEr != nil {
			fmt.println(os.get_last_error())
			fmt.panicf("Error writing cfg contents to file: ", cfgFile)
		}
	}

	fmt.printfln("Generated %v config files", len(filteredMaps))
}

press_to_exit :: proc() {
	buf: [256]byte
	fmt.print("Program complete.\nPress Enter to exit: ")
	n, _ := os.read(os.stdin, buf[:])

	// Wait for input to proceed
	for len(string(buf[:n])) < 1 {
	}
	os.exit(0)
}
