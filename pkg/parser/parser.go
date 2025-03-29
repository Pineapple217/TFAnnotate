package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pineapple217/TFAnnotate/pkg/comment"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func Parse(path string) ([]*hclsyntax.Block, error) {
	rootDir := path

	files, err := findHCLFiles(rootDir)
	if err != nil {
		return nil, err
	}
	blocks := []*hclsyntax.Block{}

	parser := hclparse.NewParser()
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		hclFile, diags := parser.ParseHCL(content, file)
		if diags.HasErrors() {
			return nil, err
		}
		blocks = append(blocks, hclFile.Body.(*hclsyntax.Body).Blocks...)
	}
	return blocks, nil
}

func ClearAll(path string) {
	findHCLFiles(path)
}

func findHCLFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tf" {
			err = ClearComments(path)
			if err != nil {
				return err
			}
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ClearComments(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, comment.CommentSymbol) {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	return err
}
