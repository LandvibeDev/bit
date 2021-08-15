package cmd

import (
	"github.com/chriswalz/complete/v3"
	"github.com/spf13/cobra"
)

// gitAccount represents the gitAccount command
var gitAccountCmd = &cobra.Command{
	Use:   "gitAccount",
	Short: "(Pre-alpha) Commit using gitAccount",
	Long:  `bit save gitAccount"`,
	Run: func(cmd *cobra.Command, args []string) {
		gitmojiSuggestions := GitmojiSuggestions()
		suggestionTree := &complete.CompTree{
			Sub: map[string]*complete.CompTree{
				"gitAccount": {
					Dynamic: toAutoCLI(gitmojiSuggestions),
				},
			},
		}
		SuggestionPrompt("> bit gitAccount ", specificCommandCompleter("gitAccount", suggestionTree))
	},
}

func init() {
	BitCmd.AddCommand(gitAccountCmd)
}
