package deps

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma"
)

// nolint:noglobal
var javaScriptExtensionRegex = regexp.MustCompile(`\.\w{1,4}$`)

// StateJavaScript is a token parsing state.
type StateJavaScript int

const (
	// StateJavaScriptUnknown represents an unknown token parsing state.
	StateJavaScriptUnknown StateJavaScript = iota
	// StateJavaScriptImport means we are in import section during token parsing.
	StateJavaScriptImport
)

// ParserJavaScript is a dependency parser for the JavaScript programming language.
// It is not thread safe.
type ParserJavaScript struct {
	State  StateJavaScript
	Output []string
}

// Parse parses dependencies from JavaScript file content via ReadCloser using the chroma JavaScript lexer.
func (p *ParserJavaScript) Parse(reader io.ReadCloser, lexer chroma.Lexer) ([]string, error) {
	defer reader.Close()

	p.init()
	defer p.init()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from reader: %s", err)
	}

	iter, err := lexer.Tokenise(nil, string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize file content: %s", err)
	}

	for _, token := range iter.Tokens() {
		p.processToken(token)
	}

	return p.Output, nil
}

func (p *ParserJavaScript) append(dep string) {
	// trim whitespaces, single quotes and double quotes
	dep = strings.Trim(dep, `"' `)

	// if front slash path, select last element
	splitted := strings.Split(dep, `/`)
	dep = splitted[len(splitted)-1]

	// if back slash path, select last element
	splitted = strings.Split(dep, `\`)
	dep = splitted[len(splitted)-1]

	// remove extension
	dep = javaScriptExtensionRegex.ReplaceAllString(dep, "")

	p.Output = append(p.Output, dep)
}

func (p *ParserJavaScript) init() {
	p.State = StateJavaScriptUnknown
	p.Output = nil
}

func (p *ParserJavaScript) processToken(token chroma.Token) {
	switch {
	case token.Type == chroma.KeywordReserved:
		p.processKeywordReserved(token.Value)
	case token.Type == chroma.LiteralStringSingle:
		p.processLiteralStringSingle(token.Value)
	case token.Type == chroma.Punctuation:
		p.processPunctuation(token.Value)
	}
}

func (p *ParserJavaScript) processKeywordReserved(value string) {
	switch value {
	case "import":
		p.State = StateJavaScriptImport
	default:
		p.State = StateJavaScriptUnknown
	}
}

func (p *ParserJavaScript) processLiteralStringSingle(value string) {
	switch p.State {
	case StateJavaScriptImport:
		p.append(value)
	default:
		p.State = StateJavaScriptUnknown
	}
}

func (p *ParserJavaScript) processPunctuation(value string) {
	if p.State == StateJavaScriptImport && value == ";" {
		p.State = StateJavaScriptUnknown
	}
}
