package main

import (
	"bytes"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/dsl"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/ThoughtWorksStudios/bobcat/interpreter"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	bf "github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"path/filepath"
	re "regexp"
	"strings"
	"testing"
)

const (
	parseOnly     = "example-parse-only"
	ensureSuccess = "example-success"
	ensureFail    = "example-fail"
	DOCS          = "docs/*.md"
)

var NON_ALPHA *re.Regexp = re.MustCompile("[^a-z0-9]+")
var LEADING_HYPH *re.Regexp = re.MustCompile("^-+")
var TRAILING_HYPH *re.Regexp = re.MustCompile("-+$")
var URL *re.Regexp = re.MustCompile("^http[s]?:")

func TestCodeBlocksInDocumentationShouldBeValid(t *testing.T) {
	files := allDocs(t)

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		AssertNil(t, err, "Should not receive error when reading %q", file)

		AssertNotEqual(t, "", string(content))
		bf.Markdown(content, NewCodeBlockValidator(t, file, validateCodeBlocks), bf.EXTENSION_FENCED_CODE)
	}
}

func TestDocumentationShouldNotHaveBrokenLinks(t *testing.T) {
	files := allDocs(t)

	headers, links := make(SetMap), make(SetMap)
	collectHeaders := func(file, headerText string) { headers.add(file, toGithubHeaderLink(headerText)) }
	collectLinks := func(file, link string) {
		if !(URL.MatchString(link)) {
			links.add(file, link)
		}
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		AssertNil(t, err, "Should not receive error when reading %q", file)

		AssertNotEqual(t, "", string(content))

		bf.Markdown(content, NewLinkCollector(t, file, collectHeaders, collectLinks), bf.EXTENSION_FENCED_CODE)
	}

	for file, linksInFile := range links {
		for _, lk := range linksInFile {

			var linkedFile string
			var linkedHeader string

			if strings.Contains(lk, "#") {
				i := strings.Index(lk, "#")

				if i == 0 {
					linkedFile = file
				} else {
					linkedFile = lk[:i]
				}

				linkedHeader = lk[i:]
			} else {
				linkedFile = lk
				linkedHeader = ""
			}

			if file != linkedFile {
				linkedFile = filepath.Join(filepath.Dir(file), linkedFile)
			}

			Assert(t, contains(files, linkedFile), "BROKEN LINK: cannot resolve file %q from link %q", linkedFile, lk)

			if "" != linkedHeader {
				linkables, ok := headers[linkedFile]
				Assert(t, ok, "BROKEN LINK: File %q should have at least 1 header entry because it is referenced by %q", linkedFile, lk)
				Assert(t, contains(linkables, linkedHeader), "BROKEN LINK: file %q does not contain header reference %q for link %q", linkedFile, linkedHeader, lk)
			}
		}

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

func validateCodeBlocks(t *testing.T, file string, text []byte, lang string) {
	switch lang {
	case parseOnly, ensureSuccess, ensureFail:
		r, err := dsl.Parse("inline", text, dsl.Recover(false))
		AssertNil(t, err, "Should not receive error parsing code block.\n\nFile: %q\n\nCode:\n\n```%s\n%s```", file, lang, string(text))

		if parseOnly != lang {
			ast := r.(*Node)
			i := interpreter.New(NewDummyEmitter(), false)
			_, err = i.Visit(ast, NewRootScope(), false)

			if ensureSuccess == lang {
				AssertNil(t, err, "Should not receive error evaluating code block.\n\nFile: %q\n\nCode:\n\n```%s\n%s```", file, lang, string(text))
			} else {
				ExpectsError(t, "", err)
			}
		}
	case "":
		Assert(t, false, "You MUST tag your code blocks.\n\nFile: %q\n\nCode:\n\n```%s\n%s```", file, lang, string(text))
	case "bash", "dos": // Whitelist ignorable code blocks here
		return
	default:
		/**
		 * We use language tags to whether or not to evaluate the code block, and determine how far to validate.
		 * If you need to allow another language tag, add it to the list above. Otherwise, this `default` case should
		 * catch typos in tag names, and other unhandled tags.
		 */
		Assert(t, false, "Unexpected language tag %q.\n\nFile: %q\n\nCode:\n\n```%s\n%s```", lang, file, lang, string(text))
	}
}

type SetMap map[string][]string

func (s SetMap) add(key, value string) string {
	if set, present := s[key]; present {
		s[key] = append(set, value)
	} else {
		s[key] = []string{value}
	}

	return value
}

type TextCollector func(file, text string)
type CodeValidator func(t *testing.T, file string, text []byte, lang string)

func NewCodeBlockValidator(t *testing.T, file string, validator CodeValidator) bf.Renderer {
	return &DocTester{t: t, file: file, onCodeBlock: validator}
}

func NewLinkCollector(t *testing.T, file string, headerCollector TextCollector, linkCollector TextCollector) bf.Renderer {
	return &DocTester{t: t, file: file, onHeader: headerCollector, onLink: linkCollector}
}

type DocTester struct {
	t           *testing.T
	file        string
	onCodeBlock CodeValidator
	onHeader    TextCollector
	onLink      TextCollector
}

func (ar *DocTester) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if nil != ar.onCodeBlock {
		ar.onCodeBlock(ar.t, ar.file, text, lang)
	}
}

func (ar *DocTester) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	if nil != ar.onLink {
		ar.onLink(ar.file, string(link))
	}
}

func (ar *DocTester) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	start := out.Len()
	if text() && nil != ar.onHeader {
		headerText := string(out.Bytes()[start:])
		ar.onHeader(ar.file, headerText)
	}
}

func (ar *DocTester) HRule(out *bytes.Buffer)                                               {}
func (ar *DocTester) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (ar *DocTester) AutoLink(out *bytes.Buffer, link []byte, kind int)                     {}
func (ar *DocTester) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte)        {}
func (ar *DocTester) LineBreak(out *bytes.Buffer)                                           {}
func (ar *DocTester) RawHtmlTag(out *bytes.Buffer, tag []byte)                              {}
func (ar *DocTester) FootnoteRef(out *bytes.Buffer, ref []byte, id int)                     {}
func (ar *DocTester) Entity(out *bytes.Buffer, entity []byte)                               {}
func (ar *DocTester) DocumentHeader(out *bytes.Buffer)                                      {}
func (ar *DocTester) DocumentFooter(out *bytes.Buffer)                                      {}
func (ar *DocTester) List(out *bytes.Buffer, text func() bool, flags int)                   { text() }
func (ar *DocTester) Paragraph(out *bytes.Buffer, text func() bool)                         { text() }
func (ar *DocTester) Footnotes(out *bytes.Buffer, text func() bool)                         { text() }
func (ar *DocTester) BlockQuote(out *bytes.Buffer, text []byte)                             { out.Write(text) }
func (ar *DocTester) BlockHtml(out *bytes.Buffer, text []byte)                              { out.Write(text) }
func (ar *DocTester) TableRow(out *bytes.Buffer, text []byte)                               { out.Write(text) }
func (ar *DocTester) TableHeaderCell(out *bytes.Buffer, text []byte, flags int)             { out.Write(text) }
func (ar *DocTester) TableCell(out *bytes.Buffer, text []byte, flags int)                   { out.Write(text) }
func (ar *DocTester) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int)          { out.Write(text) }
func (ar *DocTester) TitleBlock(out *bytes.Buffer, text []byte)                             { out.Write(text) }
func (ar *DocTester) Emphasis(out *bytes.Buffer, text []byte)                               { out.Write(text) }
func (ar *DocTester) ListItem(out *bytes.Buffer, text []byte, flags int)                    { out.Write(text) }
func (ar *DocTester) CodeSpan(out *bytes.Buffer, text []byte)                               { out.Write(text) }
func (ar *DocTester) DoubleEmphasis(out *bytes.Buffer, text []byte)                         { out.Write(text) }
func (ar *DocTester) TripleEmphasis(out *bytes.Buffer, text []byte)                         { out.Write(text) }
func (ar *DocTester) StrikeThrough(out *bytes.Buffer, text []byte)                          { out.Write(text) }
func (ar *DocTester) NormalText(out *bytes.Buffer, text []byte)                             { out.Write(text) }
func (ar *DocTester) GetFlags() int                                                         { return 0 }

func allDocs(t *testing.T) (files []string) {
	var err error

	if files, err = filepath.Glob(DOCS); err != nil {
		t.Fatalf("Should not receive error while retrieving documentation: %v", err)
		return
	}

	if len(files) == 0 {
		t.Fatalf("Found no documentation with pattern %q", DOCS)
		return
	}

	files = append(files, "README.md")
	return
}

func toGithubHeaderLink(headerText string) string {
	return "#" + string(
		LEADING_HYPH.ReplaceAll(
			TRAILING_HYPH.ReplaceAll(
				NON_ALPHA.ReplaceAll(
					[]byte(strings.ToLower(headerText)), []byte("-"),
				), []byte{},
			), []byte{},
		),
	)
}

func contains(set []string, subject string) bool {
	for _, expected := range set {
		if expected == subject {
			return true
		}
	}
	return false
}
