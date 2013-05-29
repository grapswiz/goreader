// from goquery https://github.com/PuerkitoBio/goquery

package goreader

import (
	"bytes"
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"net/url"
)

type Document struct {
	*Selection
	url      *url.URL
	rootNode *html.Node
}

type Selection struct {
	Nodes    []*html.Node
	document *Document
	prevSel  *Selection
}

func (this *Selection) Each(f func(int, *Selection)) *Selection {
	for i, n := range this.Nodes {
		f(i, newSingleSelection(n, this.document))
	}
	return this
}

func (this *Selection) Text() string {
	var buf bytes.Buffer

	for _, n := range this.Nodes {
		buf.WriteString(getNodeText(n))
	}
	return buf.String()
}

func getNodeText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	} else if node.FirstChild != nil {
		var buf bytes.Buffer
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			buf.WriteString(getNodeText(c))
		}
		return buf.String()
	}
	return ""
}

func (this *Selection) Attr(attrName string) (val string, exists bool) {
	if len(this.Nodes) == 0 {
		return
	}
	return getAttributeValue(attrName, this.Nodes[0])
}

func getAttributeValue(attrName string, n *html.Node) (val string, exists bool) {
	if n == nil {
		return
	}

	for _, a := range n.Attr {
		if a.Key == attrName {
			val = a.Val
			exists = true
			return
		}
	}
	return
}

func pushStack(fromSel *Selection, nodes []*html.Node) (result *Selection) {
	result = &Selection{nodes, fromSel.document, fromSel}
	return
}

func (this *Selection) Find(selector string) *Selection {
	return pushStack(this, findWithSelector(this.Nodes, selector))
}

func findWithSelector(nodes []*html.Node, selector string) []*html.Node {
	sel := cascadia.MustCompile(selector)

	return mapNodes(nodes, func(i int, n *html.Node) (result []*html.Node) {
		for f := n.FirstChild; f != nil; f = f.NextSibling {
			if f.Type == html.ElementNode {
				result = append(result, sel.MatchAll(f)...)
			}
		}
		return
	})
}

func mapNodes(nodes []*html.Node, f func(int, *html.Node) []*html.Node) (result []*html.Node) {
	for i, n := range nodes {
		if vals := f(i, n); len(vals) > 0 {
			result = appendWithoutDuplicates(result, vals)
		}
	}
	return
}

func appendWithoutDuplicates(target []*html.Node, nodes []*html.Node) []*html.Node {
	for _, n := range nodes {
		if !isInSlice(target, n) {
			target = append(target, n)
		}
	}

	return target
}

func isInSlice(slice []*html.Node, node *html.Node) bool {
	return indexInSlice(slice, node) > -1

}

func indexInSlice(slice []*html.Node, node *html.Node) int {
	if node != nil {
		for i, n := range slice {
			if n == node {
				return i
			}
		}
	}
	return -1
}

func NewDocument(root *html.Node, url *url.URL) (d *Document) {
	d = &Document{nil, url, root}
	d.Selection = newSingleSelection(root, d)
	return
}

func newSingleSelection(node *html.Node, doc *Document) *Selection {
	return &Selection{[]*html.Node{node}, doc, nil}
}
