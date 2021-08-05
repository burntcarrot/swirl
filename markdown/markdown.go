package markdown

import (
	"os"
	"path/filepath"
	"text/template"
	"time"

	bfc "github.com/Depado/bfchroma"
	bf "github.com/russross/blackfriday/v2"
)

var bfFlags = bf.UseXHTML | bf.Smartypants | bf.SmartypantsFractions |
	bf.SmartypantsDashes | bf.NofollowLinks
var bfExts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
	bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
	bf.HeadingIDs | bf.Footnotes | bf.NoEmptyLineBeforeBlock

type Output struct {
	HTML []byte
	Meta Matter
}

// Renders markdown to html, and fetches metadata.
func (out *Output) RenderMarkdown(source []byte) {
	md := MarkdownDoc{}
	md.Extract(source)

	out.HTML = bf.Run(
		md.Body,
		bf.WithRenderer(
			bfc.NewRenderer(
				bfc.ChromaStyle(Icy),
				bfc.Extend(
					bf.NewHTMLRenderer(bf.HTMLRendererParameters{
						Flags: bfFlags,
					}),
				),
			),
		),
		bf.WithExtensions(bfExts),
	)
	out.Meta = md.Frontmatter
}

// Renders out.HTML into dst html file, using the template specified
// in the frontmatter. data is the template struct.
func (out *Output) RenderHTML(dst, tmplDir string, data interface{}) error {
	metaTemplate := out.Meta["template"]
	if metaTemplate == "" {
		metaTemplate = "text.html"
	}

	t, err := template.New("").Funcs(template.FuncMap{
		"parsedate": func(s string) time.Time {
			date, _ := time.Parse("2006-01-02", s)
			return date
		},
	}).ParseGlob(filepath.Join(tmplDir, "*.html"))
	if err != nil {
		return err
	}

	w, err := os.Create(dst)
	if err != nil {
		return err
	}

	if err = t.ExecuteTemplate(w, metaTemplate, data); err != nil {
		return err
	}
	return nil
}