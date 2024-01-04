package cmd

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func findAll(node *html.Node, nodeName string) []*html.Node {
	var gathered []*html.Node

	if node.Type == html.ElementNode {
		if node.Data == nodeName {
			gathered = append(gathered, node)
		}
	}

	for childNode := node.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		nestedGather := findAll(childNode, nodeName)
		gathered = append(gathered, nestedGather...)
	}

	return gathered
}

func getAttrFromNode(node *html.Node, attrName string) string {
	for _, attr := range node.Attr {
		if strings.ToLower(attr.Key) == attrName {
			return attr.Val
		}
	}

	return ""
}

func renderNodeToString(node *html.Node) string {
	var b bytes.Buffer
	html.Render(&b, node)
	return b.String()
}

func TestFixLinksWithExclude(t *testing.T) {
	htmlContent := `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<div>
	<p>Hello world</p>
	<a href="https://example.com">Example link</a>
	<a href="#">Example link</a>
	<a href="https://example.com/2" class="foo-exclude">Example link</a>
	<a href="https://example.com/3" referrerpolicy="noopener">Example link</a>
	<div class="foo-exclude">
		<a href="https://example.com">Example link</a>
		<a href="#">Example link</a>
		<a href="https://example.com/2">Example link</a>
		<a href="https://example.com/3" referrerpolicy="noopener">Example link</a>
	</div>
</div>
	`

	doc, _ := html.Parse(strings.NewReader(htmlContent))
	recursivePatchNode(doc, []string{"a"}, "noopener noreferer", "foo-exclude")

	checks := []string{"noopener noreferer", "noopener noreferer", "", "noopener noreferer", "", "", "", "noopener"}

	for idx, link := range findAll(doc, "a") {
		relValue := getAttrFromNode(link, "referrerpolicy")
		if relValue != checks[idx] {
			t.Log(renderNodeToString(doc))
			t.Logf(`EXPECTED: referrerpolicy="%s"`, checks[idx])
			t.Logf("IS: %s", renderNodeToString(link))
			t.FailNow()
			break
		}
	}
}
