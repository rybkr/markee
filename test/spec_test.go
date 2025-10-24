package spec_test

import (
	"encoding/json"
    "fmt"
	"markee/internal/parser"
	"markee/internal/renderer"
	"os"
	"path/filepath"
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

func TestCommonmarkCompliance(t *testing.T) {
	examples := loadSpec(t)
    var passed, failed int

	for _, ex := range examples {
		ex := ex
		t.Run(fmt.Sprintf("%s#%d", ex.Section, ex.Example), func(t *testing.T) {
			doc := parser.Parse(ex.Markdown)
			got := renderer.RenderHTML(doc)

			if got != ex.HTML {
				t.Errorf("\x1b[31mexample %d\x1b[0m (%s) mismatch\nMarkdown:\n%s\nExpected:\n%s\nGot:\n%s",
					ex.Example, ex.Section, ex.Markdown, ex.HTML, got,
				)
				failed++
			} else {
				passed++
			}
		})
	}

    if failed > 0 {
        t.Logf("Passed %d, Failed %d, Total %d", passed, failed, passed + failed)
    }
}
