package spec_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"markee/internal/parser"
	"markee/internal/renderer"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type CommonmarkExample struct {
	Markdown  string `json:"markdown"`
	HTML      string `json:"html"`
	Example   int    `json:"example"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Section   string `json:"section"`
}

func loadSpec(t *testing.T) []CommonmarkExample {
	path := filepath.Join("data", "spec.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read spec.json: %v", err)
	}

	var examples []CommonmarkExample
	if err := json.Unmarshal(data, &examples); err != nil {
		t.Fatalf("cannot parse spec.json: %v", err)
	}

	return examples
}

func showWhitespace(s string) string {
	s = strings.ReplaceAll(s, "\t", "→")
	s = strings.ReplaceAll(s, " ", "‿")
    s = strings.ReplaceAll(s, "\n", "↓\n")
    return s
}

var categoryFilter string

func TestMain(m *testing.M) {
	flag.StringVar(&categoryFilter, "category", "", "Run only examples from the specified section")
	flag.Parse()
	os.Exit(m.Run())
}

func TestCommonmarkCompliance(t *testing.T) {
	examples := loadSpec(t)
	passed, failed := 0, 0

	for _, ex := range examples {
		if categoryFilter != "" && ex.Section != categoryFilter {
			continue
		}

		ex := ex
		t.Run(fmt.Sprintf("%s#%d", ex.Section, ex.Example), func(t *testing.T) {
			doc := parser.Parse(ex.Markdown)
			got := renderer.RenderHTML(doc)

			if got != ex.HTML {
				t.Errorf("\x1b[31mexample %d\x1b[0m (%s) mismatch\nMarkdown:\n%s\nExpected:\n%s\nGot:\n%s",
					ex.Example, ex.Section, showWhitespace(ex.Markdown), showWhitespace(ex.HTML), showWhitespace(got),
				)
				failed++
			} else {
				passed++
			}
		})
	}

	t.Logf("Passed %d, Failed %d, Total %d", passed, failed, passed+failed)
}
