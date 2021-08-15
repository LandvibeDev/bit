package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/chriswalz/complete/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// gitAccount represents the gitAccount command
var gitUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage git account",
	Long:  "Manage git account",
	Run: func(cmd *cobra.Command, args []string) {
		userFunction := ""
		userFunctionSuggestions := UserSuggestions()
		suggestionTree := &complete.CompTree{
			Sub: map[string]*complete.CompTree{
				"user": {
					Dynamic: toAutoCLI(userFunctionSuggestions),
				},
			},
		}
		userFunction = SuggestionPrompt("> bit user ", specificCommandCompleter("user", suggestionTree))
		println("input", userFunction)

		if userFunction == "addUser" {

			username := ""
			survey.AskOne(&survey.Input{
				Message: "input username",
			}, &username)

			email := ""
			survey.AskOne(&survey.Input{
				Message: "input email",
			}, &email)

			token := ""
			survey.AskOne(&survey.Input{
				Message: "input token",
			}, &token)

			addUser(username, email, token)
		} else if userFunction == "deleteUser" {
			// TODO
		} else if userFunction == "resetUser" {
			// TODO
		} else if userFunction == "listUser" {
			listUser()
		}
	},
}

func init() {
	initLocalStorage()
	BitCmd.AddCommand(gitUserCmd)
}

func UserSuggestions() []complete.Suggestion {
	var suggestions []complete.Suggestion
	for _, userFunction := range userFunctions {
		suggestions = append(suggestions, complete.Suggestion{
			Name: userFunction,
		})
	}
	return suggestions
}

func initLocalStorage() {
	// TODO create directory, file if not exist

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/bit")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
func addUser(userName string, email string, token string) {
	// TODO valid username check

	// TODO if user name exist question for overwrite
	viper.Set("users."+userName+".name", userName)
	viper.Set("users."+userName+".email", email)
	viper.Set("users."+userName+".token", token)
	viper.WriteConfig()
}

func listUser() {
	var users = viper.Get("users")

	fmt.Println("users: ", users)
}

var userFunctions = []string{
	"addUser",
	"deleteUser",
	"resetUser",
	"listUser",
}
