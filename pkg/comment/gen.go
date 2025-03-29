package comment

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Pineapple217/TFAnnotate/pkg/state"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const CommentSymbol = "///"

type GenFiles map[string]*GenFile

type GenFile struct {
	Path        string
	FileUpdates []FileUpdate
}

type FileUpdate struct {
	Line    int
	Comment string
}

func Gen(s state.State, blocks []*hclsyntax.Block, c Config) {
	files := GenFiles{}
	for _, b := range blocks {
		for _, a := range c.Annotations {
			if a.Module != nil {
				if b.Type == "module" {
					v, _ := b.Body.Attributes["source"].Expr.Value(&hcl.EvalContext{})
					source := v.AsString()
					if a.Module.Source != source {
						continue
					}
					values := map[string]string{}
					for _, value := range a.Values {
						tt := strings.Split(value.Target, ".")
						r, err := s.GetResource(state.Query{
							Module: b.Labels[0],
							Type:   tt[0],
							Name:   tt[1],
						})
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						extractedValue := r.Instances[0]["attributes"].(map[string]any)[tt[2]].(string)
						values[value.Name] = extractedValue
					}
					tp, err := template.New(a.Name).Parse(a.Comment)
					if err != nil {
						panic(err)
					}
					buf := &bytes.Buffer{}
					err = tp.Execute(buf, values)
					if err != nil {
						panic(err)
					}
					fileName := b.Body.SrcRange.Filename
					files.AddFileUpdate(fileName, FileUpdate{
						Line:    b.Body.SrcRange.Start.Line,
						Comment: buf.String(),
					})
				}

			}
		}
	}
	err := files.Insert()
	if err != nil {
		panic(err)
	}
}

func (gfs GenFiles) AddFileUpdate(fileName string, fu FileUpdate) {
	if gf, exists := gfs[fileName]; exists {
		gf.FileUpdates = append(gf.FileUpdates, fu)
	} else {
		gfs[fileName] = &GenFile{
			Path:        fileName,
			FileUpdates: []FileUpdate{fu},
		}
	}
}

func (gfs GenFiles) Insert() error {
	for _, gf := range gfs {
		file, err := os.OpenFile(gf.Path, os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}
		for i, fu := range gf.FileUpdates {
			lineNum := fu.Line + i
			if lineNum > len(lines)+1 || lineNum < 1 {
				return fmt.Errorf("invalid line number")
			}
			lines = append(lines[:lineNum-1], append([]string{CommentSymbol + " " + fu.Comment}, lines[lineNum-1:]...)...)

			// Rewrite the file with the new content.
			file.Seek(0, 0)  // Reset the file pointer to the beginning.
			file.Truncate(0) // Clear the file content before rewriting.

			writer := bufio.NewWriter(file)
			for _, line := range lines {
				_, err := writer.WriteString(line + "\n")
				if err != nil {
					return fmt.Errorf("failed to write to file: %v", err)
				}
			}
			writer.Flush()
		}
	}
	return nil
}
