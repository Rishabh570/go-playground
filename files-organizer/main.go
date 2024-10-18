package main

import (
	"organizer/cmd"
)

// Goal
// 1. provide a dir to observe
// 2. provide a destination dir
// 3. provide file extensions to organize
// 4. run CLI with above inputs
// 5. CLI will organize files in the provided dir: .png/.jpg will go to destDir/images, .pdf will go to destDir/documents, etc.
// 6. CLI will create the destination directory if it doesn't exist
// 7. CLI will move the files to the appropriate directory
// 8. CLI will print out the files that were moved and the directories that were created
func main() {
	cmd.RootCmd.Execute()
}
