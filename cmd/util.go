package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func standardizePath(args []string) []string {
	var paths []string

	for _, path := range args {
		fileInfo, err := os.Stat(path)

		if os.IsNotExist(err) {
			continue
		}
		if !fileInfo.IsDir() {
			continue
		}

		paths = append(paths, path)
	}

	return paths
}

func findHTMLFiles(rootPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func parseHTMLFile(filePath string) (*html.Node, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func recursivePatchLinks(node *html.Node) {
	if node.Type == html.ElementNode {
		if node.Data == "a" {
			foundRelAttr := false

			for idx, attr := range node.Attr {
				if strings.ToLower(attr.Key) == "rel" {
					foundRelAttr = true
					node.Attr[idx].Val = "noreferrer"
					break
				}
			}

			if !foundRelAttr {
				node.Attr = append(node.Attr, html.Attribute{Key: "rel", Val: "noreferrer"})
			}
		}
	}

	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		recursivePatchLinks(childNode)
	}
}
