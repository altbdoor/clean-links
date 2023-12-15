package cmd

import (
	"os"
	"path/filepath"
	"slices"
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

func recursivePatchNode(node *html.Node, nodeName string, relValue string, excludeClass string) {
	if node.Type == html.ElementNode {
		var attrs []string
		for _, attr := range node.Attr {
			attrs = append(attrs, strings.ToLower(attr.Key))
		}

		classNameIdx := slices.Index(attrs, "class")
		if classNameIdx != -1 && strings.Contains(node.Attr[classNameIdx].Val, excludeClass) {
			return
		}

		if node.Data == nodeName {
			relIdx := slices.Index(attrs, "rel")
			if relIdx == -1 {
				node.Attr = append(node.Attr, html.Attribute{Key: "rel", Val: relValue})
			} else {
				node.Attr[relIdx].Val = relValue
			}
		}
	}

	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		recursivePatchNode(childNode, nodeName, relValue, excludeClass)
	}
}
