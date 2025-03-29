package main

import (
	"fmt"

	"github.com/Pineapple217/TFAnnotate/pkg/comment"
	"github.com/Pineapple217/TFAnnotate/pkg/parser"
	"github.com/Pineapple217/TFAnnotate/pkg/state"
	"github.com/spf13/cobra"
)

func main() {
	var cmdVersion = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TFAnnotate v0.1.0")
		},
	}

	var cmdParse = &cobra.Command{
		Use:  "parse [path to parse]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := parser.Parse(args[0])
			if err != nil {
				return err
			}

			s, err := state.Pull(args[0])
			if err != nil {
				return err
			}

			c, err := comment.GetConfig(args[0])
			if err != nil {
				return err
			}

			comment.Gen(s, b, c)
			return nil
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
	rootCmd.AddCommand(cmdVersion, cmdParse, cmdRemove)
	rootCmd.Execute()
}
