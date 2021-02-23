package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"errors"
	"log"
	"encoding/base32"
	// "path/filepath"

	"github.com/BurntSushi/toml"
)

const colorReset string = "\033[0m"
const colorRed string = "\033[31m"
const colorGreen string = "\033[32m"
const colorYellow string = "\033[33m"
const colorBlue string = "\033[34m"
const colorPurple string = "\033[35m"
const colorCyan string = "\033[36m"
const colorWhite string = "\033[37m"

type Options struct {
	Verbosity    int
	Color        bool
	WorkDirName  string
	RepoName     string
	RepoPath     string
	Branch       string
	DestRepoName string
	DestRepoPath string
}

type workDirConfig struct {
	WorkDirName string
	Branch      string
	DefaultRepo string
	Repos map[string]string
}

type AmforaBookmarks struct {
	Bookmarks map[string]string
}

// Halt if error parameter is not nil
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Calculate the verbosity level given parameters
func CalculateVerbosity(debug bool, verbose bool, quiet bool) int {
	if quiet {
		return 0
	}

	if debug {
		return 3
	}

	if verbose {
		return 2
	}

	return 1
}

// Set the verbosity level given command-line flags
func SetVerbosity(opts Options, debug bool, verbose bool, quiet bool) Options {
	opts.Verbosity = CalculateVerbosity(debug, verbose, quiet)
	return opts
}

// Create a folder path with appropriate permissions
func CreateFolder(folderName string, verbosity int) {
	if verbosity >= 1 {
		fmt.Printf("Creating folder %s\n", folderName)
	}
	os.MkdirAll(folderName, 0777)
}

// Create a subfolder given a parent folder
func CreateSubFolder(parentFolder string, subFolder string, verbosity int) {
	folderPath := path.Join(parentFolder, subFolder)
	CreateFolder(folderPath, verbosity)
}

// Get the user's home folder path
func GetHome() string {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		// fmt.Println(pair[0])

		if pair[0] == "HOME" || pair[0] == "USERPROFILE" {
			return pair[1]
		}
	}

	fmt.Println("Warning! No home variable defined")
	return ""
}

// Convert a Date/Time string into a format that is a valid file path
func TimeToPath(timeStr string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(timeStr, ":", "-"), "/", "-"), " ", "T")
}

// Copy the source file to a destination file. Any existing file 
// will be overwritten and will not copy file attributes.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// Load a project working directory configuration given
// the working directory path
func ReadBookMarks() (AmforaBookmarks, error) {
	// var configPath string

	// if len(workDir) == 0 {
	// 	configPath = filepath.Join(".dupver", "config.toml")
	// } else {
	// 	configPath = filepath.Join(workDir, ".dupver", "config.toml")
	// }
	bookmarksPath := "bookmarks.toml"

	return ReadBookMarksFile(bookmarksPath)
}

// Load a project working directory configuration given
// the project working directory configuration file path
func ReadBookMarksFile(filePath string) (AmforaBookmarks, error) {
	var bookmarks AmforaBookmarks

	f, err := os.Open(filePath)

	if err != nil {
		return AmforaBookmarks{}, errors.New("config file missing")
	}

	if _, err = toml.DecodeReader(f, &bookmarks); err != nil {
		log.Fatal("Invalid configuration file " + filePath)
	}

	f.Close()

	return bookmarks, nil
}

func main() {
	bookmarks, _ := ReadBookMarks()
	fmt.Println("# Amfora Bookmarks\n")

	for encUrl, name := range bookmarks.Bookmarks {
		decUrl, _ := base32.StdEncoding.DecodeString(strings.ToUpper(encUrl))
		url := string(decUrl)
		fmt.Printf("=> %s %s\n", url, name)
	}
}
