package main

import (
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
)

func Parser() parser.Parser {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, meta.Meta, mathjax.MathJax, &wikilink.Extender{}, &hashtag.Extender{}),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return md.Parser()
}
