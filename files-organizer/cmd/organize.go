package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var extensions []string
var srcDir string
var destDir string

func init() {
	OrganizeCmd.PersistentFlags().StringSliceVar(&extensions, "ext", []string{}, "File extensions to organize")
	OrganizeCmd.PersistentFlags().StringVarP(&destDir, "output", "o", "", "Destination directory")
	OrganizeCmd.PersistentFlags().StringVarP(&srcDir, "input", "i", "", "Source directory")
	RootCmd.AddCommand(OrganizeCmd)
}

var OrganizeCmd = &cobra.Command{
	Use:   "organize",
	Short: "Organizes files in a directory",
	Run: func(cmd *cobra.Command, args []string) {

		for _, ext := range extensions {
			err := walkDir(srcDir, ext)
			if err != nil {
				fmt.Printf("Error walking the path %q: %v\n", srcDir, err)
			}
			fmt.Printf("Successfully organized %q\n", ext)
		}
	},
}

// Walk the source directory
func walkDir(dir string, ext string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("err:", err)
			return err
		}

		// Check if file has the desired extension
		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			fmt.Println("name:", info.Name())

			// Construct the new path in the destination directory
			destPath := filepath.Join(destDir, info.Name())
			fmt.Printf("Moving: %s to %s\n", path, destPath)

			// Move the file
			if err := moveFile(path, destPath); err != nil {
				return fmt.Errorf("failed to move file %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", srcDir, err)
	}
	return nil
}

func moveFile(src, dst string) error {
	// Create the destination directory for the file
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	// Move the file
	return os.Rename(src, dst)
}
