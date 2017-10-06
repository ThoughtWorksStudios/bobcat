package main

import (
	"bytes"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/dsl"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/ThoughtWorksStudios/bobcat/interpreter"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	parseOnly     = "example-parse-only"
	ensureSuccess = "example-success"
	ensureFail    = "example-fail"
)

func TestCodeBlocksInDocumentationShouldBeValid(t *testing.T) {
	files := []string{"README.md"}
	matches, err := filepath.Glob("docs/*.md")
	AssertNil(t, err, "Should not receive error globbing `docs/` directory")
	files = append(files, matches...)

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		AssertNil(t, err, "Should not receive error when reading %q", file)

		AssertNotEqual(t, "", string(content))
		blackfriday.Markdown(content, &CodeBlockAssertionRenderer{t: t, file: file}, blackfriday.EXTENSION_FENCED_CODE)
	}
}

func TestLangFilesShouldBeValid(t *testing.T) {
	AssertNil(t, filepath.Walk("examples", func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".lang") {
			return nil
		}

		if err != nil {
			return err
		}

		i := interpreter.New(NewDummyEmitter(), false)

		if _, err = i.LoadFile(path, NewRootScope(), false); err != nil {
			t.Errorf("Should not receive error evaluating %q. error => %v", path, err)
			return err
		}

		return nil
	}), "Should not receive error processing examples/**.lang")
}

type CodeBlockAssertionRenderer struct {
	t    *testing.T
	file string
}

func (cb *CodeBlockAssertionRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	switch lang {
	case parseOnly, ensureSuccess, ensureFail:
		r, err := dsl.Parse("inline", text, dsl.Recover(false))
		AssertNil(cb.t, err, "Should not receive error parsing code block.\n\nFile: %q\n\nCode:\n\n```%s\n%s```", cb.file, lang, string(text))

		if parseOnly != lang {
			ast := r.(*Node)
			i := interpreter.New(NewDummyEmitter(), false)
			_, err = i.Visit(ast, NewRootScope(), false)

			if ensureSuccess == lang {
				AssertNil(cb.t, err, "Should not receive error evaluating code block.\n\nFile: %q\n\nCode:\n\n```%s\n%s```", cb.file, lang, string(text))
			} else {
				ExpectsError(cb.t, "", err)
			}
		}
	case "":
		Assert(cb.t, false, "You MUST tag your code blocks.\n\nFile: %q\n\nCode:\n\n```%s\n%s```", cb.file, lang, string(text))
	case "bash", "dos": // Whitelist ignorable code blocks here
		return
	default:
		/**
		 * We use language tags to whether or not to evaluate the code block, and determine how far to validate.
		 * If you need to allow another language tag, add it to the list above. Otherwise, this `default` case should
		 * catch typos in tag names, and other unhandled tags.
		 */
		Assert(cb.t, false, "Unexpected language tag %q.\n\nFile: %q\n\nCode:\n\n```%s\n%s```", lang, cb.file, lang, string(text))
	}

}

func (cb *CodeBlockAssertionRenderer) BlockQuote(out *bytes.Buffer, text []byte) {}
func (cb *CodeBlockAssertionRenderer) BlockHtml(out *bytes.Buffer, text []byte)  {}
func (cb *CodeBlockAssertionRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	if !text() {
		return
	}
}
func (cb *CodeBlockAssertionRenderer) HRule(out *bytes.Buffer) {}
func (cb *CodeBlockAssertionRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	if !text() {
		return
	}
}
func (cb *CodeBlockAssertionRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {}
func (cb *CodeBlockAssertionRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	if !text() {
		return
	}
}
func (cb *CodeBlockAssertionRenderer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
}
func (cb *CodeBlockAssertionRenderer) TableRow(out *bytes.Buffer, text []byte)                   {}
func (cb *CodeBlockAssertionRenderer) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {}
func (cb *CodeBlockAssertionRenderer) TableCell(out *bytes.Buffer, text []byte, flags int)       {}
func (cb *CodeBlockAssertionRenderer) Footnotes(out *bytes.Buffer, text func() bool) {
	if !text() {
		return
	}
}
func (cb *CodeBlockAssertionRenderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {}
func (cb *CodeBlockAssertionRenderer) TitleBlock(out *bytes.Buffer, text []byte)                    {}
func (cb *CodeBlockAssertionRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int)            {}
func (cb *CodeBlockAssertionRenderer) CodeSpan(out *bytes.Buffer, text []byte)                      {}
func (cb *CodeBlockAssertionRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte)                {}
func (cb *CodeBlockAssertionRenderer) Emphasis(out *bytes.Buffer, text []byte)                      {}
func (cb *CodeBlockAssertionRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
}
func (cb *CodeBlockAssertionRenderer) LineBreak(out *bytes.Buffer) {}
func (cb *CodeBlockAssertionRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
}
func (cb *CodeBlockAssertionRenderer) RawHtmlTag(out *bytes.Buffer, tag []byte)          {}
func (cb *CodeBlockAssertionRenderer) TripleEmphasis(out *bytes.Buffer, text []byte)     {}
func (cb *CodeBlockAssertionRenderer) StrikeThrough(out *bytes.Buffer, text []byte)      {}
func (cb *CodeBlockAssertionRenderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {}
func (cb *CodeBlockAssertionRenderer) Entity(out *bytes.Buffer, entity []byte)           {}
func (cb *CodeBlockAssertionRenderer) NormalText(out *bytes.Buffer, text []byte)         {}
func (cb *CodeBlockAssertionRenderer) DocumentHeader(out *bytes.Buffer)                  {}
func (cb *CodeBlockAssertionRenderer) DocumentFooter(out *bytes.Buffer)                  {}
func (cb *CodeBlockAssertionRenderer) GetFlags() int                                     { return 0 }
