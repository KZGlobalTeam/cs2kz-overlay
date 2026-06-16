//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andygrunwald/vdf"
	"golang.org/x/sys/windows/registry"
)

// Open a provided log file path, or auto-locate it when no override is given.
func getLogFile(overridePath string) (*os.File, error) {
	if overridePath != "" {
		if _, err := os.Stat(overridePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("log file not found at override path: %s", overridePath)
		}
		file, err := os.Open(overridePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open override log file: %v", err)
		}
		return file, nil
	}

	// Automatically locate the log file
	// Steam installation path
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("Steam is not installed or the registry key was not found: %v", err)
	}
	defer key.Close()

	steamPath, _, err := key.GetStringValue("SteamPath")
	if err != nil {
		return nil, fmt.Errorf("failed to read SteamPath: %v", err)
	}

	// libraryfolders.vdf path
	libraryfoldersPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	// fallback
	if _, err := os.Stat(libraryfoldersPath); os.IsNotExist(err) {
		libraryfoldersPath = filepath.Join(steamPath, "config", "libraryfolders.vdf")
	}

	if _, err := os.Stat(libraryfoldersPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("libraryfolders.vdf not found at expected locations")
	}

	library, err := os.Open(libraryfoldersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open libraryfolders.vdf: %v", err)
	}
	defer library.Close()

	// Parse the VDF file
	parser := vdf.NewParser(library)
	vdfData, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse VDF: %v", err)
	}

	// Navigate to libraryfolders
	/*
		Example structure of libraryfolders.vdf:
		"libraryfolders"
		{
			"0"
			{
				"path"		"C:\\Program Files (x86)\\Steam"
				"apps"
				{
					...
				}
			}
			"1"
			{
				"path"		"D:\\Games\\SteamLibrary"
				"apps"
				{
					"730"		"<size>"
				}
			}
		}
	*/
	libraryfolders, ok := vdfData["libraryfolders"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("libraryfolders not found in VDF")
	}

	// Search for CS2's 730; if found, traverse the path to steamapps and read appmanifest_730.acf to get the installdir, then construct the log path
	for _, folderValue := range libraryfolders {
		folder, ok := folderValue.(map[string]interface{})
		if !ok {
			continue
		}
		apps, ok := folder["apps"].(map[string]interface{})
		if !ok {
			continue
		}
		if _, hasCS2 := apps["730"]; hasCS2 {
			path, ok := folder["path"].(string)
			if !ok {
				continue
			}
			manifestPath := filepath.Join(path, "steamapps", "appmanifest_730.acf")
			f, err := os.Open(manifestPath)
			if err != nil {
				return nil, fmt.Errorf("failed to open appmanifest: %v", err)
			}
			defer f.Close()
			parser := vdf.NewParser(f)
			manifestData, err := parser.Parse()
			if err != nil {
				return nil, fmt.Errorf("failed to parse appmanifest: %v", err)
			}
			/*
				"AppState"
				{
					...
					"installdir"		"Counter-Strike Global Offensive"
					...
				}
			*/
			appState, ok := manifestData["AppState"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("AppState not found in manifest")
			}
			installdir, ok := appState["installdir"].(string)
			if !ok {
				return nil, fmt.Errorf("installdir not found in AppState")
			}
			// Construct full CS2 path
			filePath := filepath.Join(path, "steamapps", "common", installdir, "game", "csgo", "console.log")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return nil, fmt.Errorf("log file not found at expected path: %s", filePath)
			}
			return os.Open(filePath)
		}
	}
	return nil, fmt.Errorf("failed to find CS2 library path in any library folder")
}
