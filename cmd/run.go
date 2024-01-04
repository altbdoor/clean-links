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
	Args:  cobra.MinimumNArgs(1),
	Short: "Runs clean-links on a specified path",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: performRun,
}

var choiceElems = []string{"area", "img", "iframe", "script", "link"}

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

	excludeClass, _ := cmd.Flags().GetString("exclude-class")
	value, _ := cmd.Flags().GetString("value")

	validElems := []string{"a"}
	isFixAll, _ := cmd.Flags().GetBool("fix-all")

	if isFixAll {
		validElems = append(validElems, choiceElems...)
	} else {
		for _, elemName := range choiceElems {
			elemFlag, err := cmd.Flags().GetBool("fix-" + elemName)

			if err == nil && elemFlag {
				validElems = append(validElems, elemName)
			}
		}
	}

	fmt.Printf("(i) Fixing %v HTML elements.\n", validElems)

	for idx, file := range htmlFiles {
		startFileTime := time.Now()
		doc, err := parseHTMLFile(file)

		if err != nil {
			fmt.Printf(`(!) (%d/%d) Error: unable to parse "%s"\n`, idx+1, htmlFilesCount, file)
			continue
		}

		recursivePatchNode(doc, validElems, value, excludeClass)

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
	runCmd.Flags().StringP(
		"exclude-class",
		"e",
		"clean-links-exclude",
		`HTML elements with this class will not be processed. Subsequent child elements will not be processed as well.`,
	)
	runCmd.Flags().String("value", "noreferrer", `The value for the "referrerpolicy" attribute in the elements.`)

	for _, elemName := range choiceElems {
		helpMsg := fmt.Sprintf("Includes <%s> elements to be fixed. Not included by default.", elemName)
		runCmd.Flags().Bool("fix-"+elemName, false, helpMsg)
	}

	runCmd.Flags().Bool("fix-all", false, "Fixes all HTML elements listed below.")
}
