package main

import "core:fmt"
import "core:os"
import "core:path/filepath"
import "core:slice"
import "core:strings"

WISH_DIR :: "tf/custom"
MAP_DIRS :: [3]string{"maps", "download/maps", "workshop/content/maps"}
CFG_DIR :: "custom/demo-prefixer/cfg"
COMP_GAMEMODES :: []string{"cp", "koth", "pl"}
COMPPLUS_GAMEMODES :: []string{"bball", "cp", "koth", "pass", "pl", "ultiduo", "ultitrio"}
MAP_SUFFIXES :: [5]string{"_a", "_b", "_f", "_rc", "_v"}

main :: proc() {
	// Check program location
	currentDir := check_dir()
	tfDir, _ := strings.substring_to(currentDir, strings.last_index(currentDir, "/"))
	os.set_current_directory(tfDir)

	if len(os.args) > 1 {
		fmt.println("Commandline arguments: ", os.args[1:])
	}

	hasCompOnly: bool
	hasNoSuffix: bool
	compOnly: int
	noSuffix: bool
	madeDir: bool

	// Iterate through different map directories
	for dir in MAP_DIRS {
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
			msg := fmt.tprintfln("Failed to get contents of tf/%v", dir)
			panic(msg)
		}

		// Filter map names
		for m in maps {
			// Remove file extension
			mapName, mnEr := strings.substring_to(m.name, strings.index(m.name, "."))
			if !mnEr {
				fmt.printfln("Could not get map name for %v", m.name)
				continue
			}

			// Only add maps
			if strings.ends_with(m.name, ".bsp") {
				append(&filteredMaps, mapName)
			} else {
				fmt.printfln("Skipping: %v, not a map.", m.name)
			}
		}

		// Customize map names based on user input / commandline args
		if len(os.args) > 1 { 	// Has commandline args
			// CompOnly arg
			if slice.contains(os.args[1:], "componly") || slice.contains(os.args[1:], "CompOnly") {
				hasCompOnly = true
				compOnly = 1
			}
			// CompPlus arg
			if slice.contains(os.args[1:], "compplus") || slice.contains(os.args[1:], "CompPlus") {
				hasCompOnly = true
				compOnly = 2
			}
			// NoSuffix arg
			if slice.contains(os.args[1:], "nosuffix") || slice.contains(os.args[1:], "NoSuffix") {
				hasNoSuffix = true
				noSuffix = true
			}
			// Silent arg
			if slice.contains(os.args[1:], "silent") || slice.contains(os.args[1:], "Silent") {
				hasCompOnly = true
				hasNoSuffix = true
			}
		}
		for !hasCompOnly { 	// User input for CompOnly
			buf: [256]byte
			fmt.println(
				"\nWould you like to generate prefixes only for competitive gamemodes (AD, CP, PL, KOTH)?\nYou may also generate for extended competitive gamemodes (bball, passtime, ultiduo, ultitrio): [Y] / [E] / [N]",
			)
			r, coEr := os.read(os.stdin, buf[:])
			if coEr != nil {
				panic("Error getting user input for CompOnly")
			}
			response := strings.trim_right_space(string(buf[:r]))
			switch {
			case strings.equal_fold(response, "y") || strings.equal_fold(response, "yes"):
				hasCompOnly = true
				compOnly = 1
				fmt.println("Generating map prefixes for competitive gamemodes only")
			case strings.equal_fold(response, "e") || strings.equal_fold(response, "extended"):
				hasCompOnly = true
				compOnly = 2
				fmt.println("Generating map prefixes for extended competitive gamemodes")
			case strings.equal_fold(response, "n") || strings.equal_fold(response, "no"):
				hasCompOnly = true
				compOnly = 0
				fmt.println("Generating map prefixes for all gamemodes")
			}
		}
		for !hasNoSuffix { 	// User input for NoSuffix
			buf: [256]byte
			fmt.println(
				"\nWould you like to remove suffixes from map name prefixes? [Y] / [N]\nExample: cp_process_f12 -> cp_process",
			)
			r, nsEr := os.read(os.stdin, buf[:])
			if nsEr != nil {
				panic("Error getting user input for NoSuffix")
			}
			response := strings.trim_right_space(string(buf[:r]))
			fmt.printfln("Responded: %v", response)
			switch {
			case strings.equal_fold(response, "y") || strings.equal_fold(response, "yes"):
				hasNoSuffix = true
				noSuffix = true
				fmt.println("Generating map prefixes without suffixes")
			case strings.equal_fold(response, "n") || strings.equal_fold(response, "no"):
				hasNoSuffix = true
				noSuffix = false
				fmt.println("Generating map prefixes with suffixes")
			}
		}

		// Only apply prefix customizations if relevant
		mapNames: [dynamic]string
		if compOnly != 0 || noSuffix {
			filteredMaps, mapNames = customize_maps(filteredMaps, compOnly, noSuffix)
		} else {
			mapNames = filteredMaps
		}

		// Generate tf/custom/demo-prefixer/cfg if not done yet
		if !madeDir {
			make_cfg_dir()
			madeDir = true
		}

		write_cfgs(filteredMaps, mapNames)

		free_all(context.temp_allocator)

		fmt.println()
	}

	// Don't wait for input if the user has opted for silent
	if slice.contains(os.args[1:], "silent") || slice.contains(os.args[1:], "Silent") {
		os.exit(0)
	}
	fmt.println("Program complete.")
	press_to_exit()
}

check_dir :: proc() -> string {
	currentDir, cdPass := filepath.to_slash(os.get_current_directory(context.temp_allocator))
	if !cdPass {
		fmt.println("Error checking working directory")
	}
	fmt.printfln("Current directory: %v", currentDir)

	if !strings.ends_with(currentDir, WISH_DIR) {
		fmt.printfln(
			"Program is currently in: \"%v\"\nPlease put program into your \"tf/custom\" folder.",
			currentDir,
		)
		press_to_exit()
	}
	return currentDir
}

customize_maps :: proc(
	maps: [dynamic]string,
	compOnly: int,
	noSuffix: bool,
) -> (
	[dynamic]string,
	[dynamic]string,
) {
	filteredMaps: [dynamic]string
	switch compOnly {
	case 1:
		for m in maps {
			prefix, _ := strings.substring_to(m, strings.index(m, "_"))
			if slice.contains(COMP_GAMEMODES, prefix) { 	// Add if comp gamemode
				append(&filteredMaps, m)
			}
		}
	case 2:
		for m in maps {
			prefix, _ := strings.substring_to(m, strings.index(m, "_"))
			if slice.contains(COMPPLUS_GAMEMODES, prefix) { 	// Add if comp extended gamemode
				append(&filteredMaps, m)
			}
		}
	case:
		filteredMaps = maps
	}

	mapNames: [dynamic]string
	if noSuffix {
		for m, i in filteredMaps {
			suffixFound: bool
			if strings.count(m, "_") > 1 { 	// Check if prefix exists
				suffix, _ := strings.substring_from(m, strings.last_index(m, "_"))
				for s in MAP_SUFFIXES { 	// If suffix is one that should be removed
					if strings.starts_with(suffix, s) {
						name, _ := strings.substring_to(m, strings.last_index(m, "_"))
						suffixFound = true
						append(&mapNames, name)
						break
					}
				}
			}
			if !suffixFound { 	// Suffix not removed
				append(&mapNames, m)
			}
		}
	} else { 	// Set map names to their file name
		mapNames = filteredMaps
	}
	return filteredMaps, mapNames
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

write_cfgs :: proc(filteredMaps: [dynamic]string, mapNames: [dynamic]string) {
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
		cfgContents := strings.concatenate({"ds_prefix ", mapNames[i]})

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
	fmt.print("Press Enter to exit")
	buf: [256]byte
	n, _ := os.read(os.stdin, buf[:])

	// Wait for input to proceed
	for len(string(buf[:n])) < 1 {
	}
	os.exit(0)
}
