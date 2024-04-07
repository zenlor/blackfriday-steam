package steam

import (
	"bytes"
	"io"

	bf "github.com/russross/blackfriday/v2"
)

// Renderer is the rendering interface for confluence wiki output.
type Renderer struct {
	w bytes.Buffer

	// Flags allow customizing this renderer's behavior.
	Flags Flag

	lastOutputLen int
}

// Flag control optional behavior of this renderer.
type Flag int

const (
	// FlagsNone does not allow customizing this renderer's behavior.
	FlagsNone Flag = 0 << iota

	// InformationMacros allow using info, tip, note, and warning macros.
	InformationMacros

	// IgnoreMacroEscaping will not escape any text that contains starts with `{`
	// in a block of text.
	IgnoreMacroEscaping
)

var (
	quoteTag         = []byte("quote")
	codeTag          = []byte("code")
	imageTag         = []byte("img")
	strongTag        = []byte("b")
	strikethroughTag = []byte("strike")
	emTag            = []byte("i")
	linkTag          = []byte("url")
	imgTag           = []byte("img")
	liTag            = []byte("list")
	olTag            = []byte("olist")
	litemTag         = []byte("*")
	hrTag            = []byte("[hr][/hr]")
	h1Tag            = []byte("h1")
	h2Tag            = []byte("h2")
	h3Tag            = []byte("h3")
	spoilerTag       = []byte("spoiler")
	noparseTag       = []byte("noparse")
	urlTag           = []byte("url")
	tableTag         = []byte("table")
	tableRowTag      = []byte("tr")
	tableCellHeadTag = []byte("th")
	tableCellTag     = []byte("td")
)

var (
	begin      = []byte{'['}
	end        = []byte{']'}
	slash      = []byte{'/'}
	nline      = []byte{'\n'}
	spaceBytes = []byte{' '}
	equals     = []byte{'='}
)

var itemLevel = 0

func (r *Renderer) esc(w io.Writer, text []byte) {
	// r.out(w, begin, noparseTag, end)
	r.out(w, text)
	// r.out(w, begin, slash, noparseTag, end)
}

func (r *Renderer) cr(w io.Writer) {
	if r.lastOutputLen > 0 {
		r.out(w, nline)
	}
}

func (r *Renderer) out(w io.Writer, text ...[]byte) {
	r.lastOutputLen = 0
	for _, b := range text {
		w.Write(b)
		r.lastOutputLen += len(b)
	}
}

func headingTagFromLevel(level int) []byte {
	switch level {
	case 1:
		return h1Tag
	case 2:
		return h2Tag
	default:
		return h3Tag
	}
}

// RenderNode is a confluence renderer of a single node of a syntax tree.
func (r *Renderer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	case bf.Text:
		r.esc(w, node.Literal)

	case bf.Softbreak:
		break

	case bf.Hardbreak:
		w.Write(nline)

	case bf.BlockQuote:
		if entering {
			r.out(w, begin, quoteTag, end)
			r.cr(w)
		} else {
			r.out(w, begin, slash, quoteTag, end)
			r.cr(w)
		}

	case bf.CodeBlock:
		r.out(w, begin, codeTag, end)
		w.Write(node.Literal)
		r.out(w, begin, slash, codeTag, end)
		r.cr(w)

	case bf.Code:
		r.out(w, begin, codeTag, end, node.Literal, begin, slash, codeTag, end)

	case bf.Emph:
		if entering {
			r.out(w, begin, emTag, end)
		} else {
			r.out(w, begin, slash, emTag, end)
		}

	case bf.Heading:
		headingTag := headingTagFromLevel(node.Level)
		if entering {
			r.out(w, begin, headingTag, end)
		} else {
			r.out(w, begin, slash, headingTag, end)
			r.cr(w)
		}
	case bf.Image:
		if entering {
			dest := node.LinkData.Destination
			title := node.LinkData.Title
			r.out(w, begin, imageTag, end, dest, begin, slash, imageTag, end)
			if len(title) > 0 {
				r.out(w, spaceBytes, begin, emTag, end, title, begin, slash, emTag, end)
			}
		}
	case bf.Link:
		if entering {
			r.out(w, begin, linkTag)
			if dest := node.LinkData.Destination; dest != nil {
				r.out(w, equals, dest)
			}
			r.out(w, end)
		} else {
			r.out(w, begin, slash, urlTag, end)
		}
	case bf.HorizontalRule:
		r.cr(w)
		r.out(w, hrTag)
		r.cr(w)

	case bf.Item:
		if entering {
			r.out(w, begin, litemTag, end)
		}

	case bf.List:
		t := liTag
		if node.ListFlags&bf.ListTypeOrdered != 0 {
			t = olTag
		}

		if entering {
			r.out(w, begin, t, end)
			r.cr(w)
		} else {
			r.out(w, begin, slash, t, end)
			r.cr(w)
		}

	case bf.Document:
		break
	case bf.HTMLBlock:
		break
	case bf.HTMLSpan:
		break
	case bf.Paragraph:
		if !entering {
			if node.Next != nil && node.Next.Type == bf.Paragraph {
				w.Write(nline)
				w.Write(nline)
			} else {
				if node.Parent.Type != bf.Item {
					r.cr(w)
				}
				r.cr(w)
			}
		}
	case bf.Strong:
		if entering {
			r.out(w, begin, strongTag, end)
		} else {
			r.out(w, begin, slash, strongTag, end)
		}
	case bf.Del:
		if entering {
			r.out(w, begin, strikethroughTag, end)
		} else {
			r.out(w, begin, slash, strikethroughTag, end)
		}
	case bf.Table:
		if entering {
			r.out(w, begin, tableTag, end)
			r.cr(w)
		} else {
			r.out(w, begin, slash, tableTag, end)
			r.cr(w)
		}
	case bf.TableCell:
		t := tableCellTag
		if node.IsHeader {
			t = tableCellHeadTag
		}

		if entering {
			r.out(w, begin, t, end)
		} else {
			r.out(w, begin, slash, t, end)
		}
	case bf.TableHead:
		break
	case bf.TableBody:
		break
	case bf.TableRow:
		if entering {
			r.out(w, begin, tableRowTag, end)
			r.cr(w)
		} else {
			r.out(w, begin, slash, tableRowTag, end)
			r.cr(w)
		}
	default:
		panic("Unknown node type " + node.Type.String())
	}
	return bf.GoToNext
}

// Render prints out the whole document from the ast.
func (r *Renderer) Render(ast *bf.Node) []byte {
	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		return r.RenderNode(&r.w, node, entering)
	})

	return r.w.Bytes()
}

// RenderHeader writes document header (unused).
func (r *Renderer) RenderHeader(w io.Writer, ast *bf.Node) {
}

// RenderFooter writes document footer (unused).
func (r *Renderer) RenderFooter(w io.Writer, ast *bf.Node) {
}

// Run prints out the confluence document.
func Run(input []byte, opts ...bf.Option) []byte {
	r := &Renderer{Flags: InformationMacros}
	optList := []bf.Option{bf.WithRenderer(r), bf.WithExtensions(bf.CommonExtensions)}
	optList = append(optList, opts...)
	parser := bf.New(optList...)
	ast := parser.Parse([]byte(input))
	return r.Render(ast)
}
