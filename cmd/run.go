/*
Copyright © 2023 altbdoor <lancersupraskyline@gmail.com>
MIT license, see LICENSE file
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs clean-links on a specified path",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: performRun,
}

func performRun(cmd *cobra.Command, args []string) {
	startOpTime := time.Now()
	paths := standardizePath(args)
	var htmlFiles []string

	for _, path := range paths {
		foundHtmlFiles, err := findHTMLFiles(path)

		if err != nil {
			fmt.Printf(`(!) Error: error while finding HTML files in "%s", path will be excluded\n`, path)
		} else {
			htmlFiles = append(htmlFiles, foundHtmlFiles...)
		}
	}

	htmlFilesCount := len(htmlFiles)
	fmt.Printf("(i) Found %d HTML files.\n", htmlFilesCount)

	for idx, file := range htmlFiles {
		startFileTime := time.Now()
		doc, err := parseHTMLFile(file)

		if err != nil {
			fmt.Printf(`(!) (%d/%d) Error: unable to parse "%s"\n`, idx+1, htmlFilesCount, file)
			continue
		}

		recursivePatchLinks(doc)

		var buffer bytes.Buffer
		html.Render(&buffer, doc)
		os.WriteFile(file, buffer.Bytes(), 0644)

		endFileTime := time.Since(startFileTime)
		fmt.Printf("(i) (%d/%d) Took %v\n", idx+1, htmlFilesCount, endFileTime)
	}

	endOpTime := time.Since(startOpTime)
	fmt.Printf("(i) Finished! Took %v\n", endOpTime)
}

func init() {
	// https://cobra.dev/#getting-started
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// runCmd.Flags().B
}
