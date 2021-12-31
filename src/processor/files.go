package processor

import (
	"os"
	"time"
)

func getGoFilesAsync(folder string, out chan<- string, signal chan<- bool) {
	defer func() {
		// Signal that routine is done
		signal <- true
	}()

	// Get files and folders
	fs, err := os.ReadDir(folder)
	if err != nil {
		return
	}
	for _, f := range fs {
		// Ignore hidden stuff
		if f.Name()[0] == '.' {
			continue
		}

		// Search recursively
		if f.IsDir() {
			// Signal that routine is starting
			signal <- false
			go getGoFilesAsync(folder+string(os.PathSeparator)+f.Name(), out, signal)
			continue
		}

		// Check if is a .go file
		name := f.Name()
		if len(name) > 3 && name[len(name)-3:] == ".go" {
			out <- folder + string(os.PathSeparator) + name
		}
	}
}

func getGoFiles(rootFolder string) (fileList []string) {
	outChannel := make(chan string, 100)
	sigChannel := make(chan bool, 100)
	fileList = []string{}

	// Get files from directory
	go getGoFilesAsync(rootFolder, outChannel, sigChannel)

	// Get files from outchannel
	routines := 1
	for {
		select {
		case sig := <-sigChannel:
			{
				if sig {
					routines -= 1
				} else {
					routines += 1
				}
			}
		case out := <-outChannel:
			{
				fileList = append(fileList, out)
			}
		default:
			{
				if routines == 0 && len(outChannel) == 0 {
					return fileList
				}
				time.Sleep(250 * time.Millisecond)
			}
		}
	}
}
