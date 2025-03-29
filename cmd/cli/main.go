package main

import (
	"github.com/Pineapple217/TFAnnotate/pkg/comment"
	"github.com/Pineapple217/TFAnnotate/pkg/parser"
	"github.com/Pineapple217/TFAnnotate/pkg/state"
	"github.com/spf13/cobra"
)

func main() {
	var cmdParse = &cobra.Command{
		Use:  "parse [path to parse]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			b := parser.Parse(args[0])
			s := state.Pull(args[0])
			c := comment.GetConfig(args[0])
			comment.Gen(s, b, c)
		},
	}

	var cmdRemove = &cobra.Command{
		Use:   "remove [path to parse]",
		Short: "Clears all comments",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			parser.ClearAll(args[0])
		},
	}

	var rootCmd = &cobra.Command{Use: "tfa"}
	rootCmd.AddCommand(cmdParse, cmdRemove)
	rootCmd.Execute()
}
