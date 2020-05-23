package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Tree 命令
func Tree(out io.Writer, path string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	startFolder := myFolder{path: path}
	startFolder.setFiles(files)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if printFiles {
		err := showTreeWithFiles(out, startFolder)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	} else {
		err := showTreeNoFiles(out, startFolder)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}

// myFolder : myFolder describes folder and contains it's path, counter of folders being last, files list, and it's identation string
type myFolder struct {
	path        string
	lastCount   int
	files       map[string]bool
	levelSymbol string
}

// Get filename and IsDir from os.FileInfo and create map from those values. Then set this map to 'files' variable of myFolder type
func (f *myFolder) setFiles(filesInfo []os.FileInfo) {
	newFiles := make(map[string]bool, 0)
	for _, file := range filesInfo {
		newFiles[file.Name()] = file.IsDir()
	}
	f.files = newFiles
}

// Reverse string
func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

// List folders without files
func showTreeNoFiles(out io.Writer, Folder myFolder) error {
	// Essential variables
	fileNotLastSymbol := "├───"
	fileLastSymbol := "└───"
	files := Folder.files
	lastCount := Folder.lastCount

	// Get all folders and sort them
	folders := make([]string, 0)
	for file, isDir := range files {
		if isDir {
			folders = append(folders, file)
		}
	}
	sort.Strings(folders)

	for idx, folder := range folders {
		fullPath := Folder.path + string(os.PathSeparator) + folder
		folderIsLast := idx == len(folders)-1
		fileLevel := strings.Count(fullPath, string(os.PathSeparator))

		beforeFile := ""
		// So, if folder is last, then we need to delete 1 '|' symbol so that's why we store count of last folders
		if folderIsLast {
			beforeFile = fileLastSymbol
			lastCount++
		} else {
			beforeFile = fileNotLastSymbol
		}

		// If this is not starting foder, then we add '|' before then name
		levelSymbol := ""
		if Folder.levelSymbol == "" {
			for i := 1; i < fileLevel; i++ {
				levelSymbol += "│\t"
			}
		} else {
			levelSymbol = Folder.levelSymbol + "│\t"
		}

		// If files are in last folder, then we need to delete some of the '|'
		if Folder.lastCount > 0 {
			for i := 0; i < Folder.lastCount; i++ {
				parentCount := strings.Count(Folder.levelSymbol, "│")
				childCount := strings.Count(levelSymbol, "│")
				if parentCount == childCount {
					break
				}
				levelReversed := reverse(levelSymbol)
				levelReversed = strings.Replace(levelReversed, "│", "", 1)
				levelSymbol = reverse(levelReversed)
			}
		}

		if !folderIsLast {
			lastCount = 0
		}
		fmt.Fprintf(out, "%v%v%v\n", levelSymbol, beforeFile, folder)

		allFiles, _ := ioutil.ReadDir(fullPath)
		newFolder := myFolder{path: fullPath, lastCount: lastCount, levelSymbol: levelSymbol}
		newFolder.setFiles(allFiles)

		showTreeNoFiles(out, newFolder)
	}
	return nil
}

// List all folders with files and with their size
func showTreeWithFiles(out io.Writer, Folder myFolder) error {
	fileNotLastSymbol := "├───"
	fileLastSymbol := "└───"
	files := Folder.files
	lastCount := Folder.lastCount

	allFiles := make([]string, 0)
	for file := range files {
		allFiles = append(allFiles, file)
	}
	sort.Strings(allFiles)

	for idx, file := range allFiles {
		fullPath := Folder.path + string(os.PathSeparator) + file
		fileIsLast := idx == len(allFiles)-1
		fileLevel := strings.Count(fullPath, string(os.PathSeparator))

		beforeFile := ""
		if fileIsLast {
			beforeFile = fileLastSymbol
			lastCount++
		} else {
			beforeFile = fileNotLastSymbol
			lastCount = 0
		}
		levelSymbol := ""
		if Folder.levelSymbol == "" {
			for i := 1; i < fileLevel; i++ {
				levelSymbol += "│\t"
			}
		} else {
			levelSymbol = Folder.levelSymbol + "│\t"
		}

		if Folder.lastCount > 0 {
			for i := 0; i < Folder.lastCount; i++ {
				parentCount := strings.Count(Folder.levelSymbol, "│")
				childCount := strings.Count(levelSymbol, "│")
				if parentCount == childCount {
					break
				}
				levelReversed := reverse(levelSymbol)
				levelReversed = strings.Replace(levelReversed, "│", "", 1)
				levelSymbol = reverse(levelReversed)
			}
		}

		if isDir := files[file]; isDir {
			fmt.Fprintf(out, "%v%v%v\n", levelSymbol, beforeFile, file)
			allFiles, _ := ioutil.ReadDir(fullPath)
			newFolder := myFolder{path: fullPath, lastCount: lastCount, levelSymbol: levelSymbol}
			newFolder.setFiles(allFiles)
			showTreeWithFiles(out, newFolder)
		} else {
			fileStats, _ := os.Stat(fullPath)
			size := fileStats.Size()
			sizeStr := ""
			if size == 0 {
				sizeStr = "empty"
			} else {
				sizeStr = strconv.FormatInt(size, 10) + "b"
			}
			fmt.Fprintf(out, "%v%v%v (%v)\n", levelSymbol, beforeFile, file, sizeStr)
		}

	}
	return nil
}
