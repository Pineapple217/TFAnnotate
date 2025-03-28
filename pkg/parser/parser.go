package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Pineapple217/TFAnnotate/pkg/comment"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func Parse(path string) []*hclsyntax.Block {
	rootDir := path

	files, err := findHCLFiles(rootDir)
	if err != nil {
		fmt.Println("Error scanning directory:", err)
		return nil
	}
	blocks := []*hclsyntax.Block{}

	parser := hclparse.NewParser()
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		hclFile, diags := parser.ParseHCL(content, file)
		if diags.HasErrors() {
			continue
		}
		blocks = append(blocks, hclFile.Body.(*hclsyntax.Body).Blocks...)
	}
	return blocks
}

func findHCLFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tf" {
			ClearComments(path)
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ClearComments(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
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
		fmt.Println("Error reading file:", err)
		return
	}

	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}
