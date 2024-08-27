package parser

import (
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
)

func VaultParser() parser.Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.New(meta.WithStoresInDocument()),
			mathjax.MathJax,
			&wikilink.Extender{},
			&hashtag.Extender{
				Resolver: nil,
				Variant:  hashtag.ObsidianVariant,
			},
			emoji.Emoji,
		),
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
