package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/chriswalz/complete/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"regexp"
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

			userName := ""
			survey.AskOne(&survey.Input{
				Message: "input username",
			}, &userName)

			var isUserNameString = regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString
			if !isUserNameString(userName) || len(userName) == 0 {
				fmt.Println("UserName is invalid: ", userName)
				return
			}

			email := ""
			survey.AskOne(&survey.Input{
				Message: "input email",
			}, &email)

			var isEmailString = regexp.MustCompile(`^[_a-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$`).MatchString
			if isEmailString(email) {
				fmt.Println("Email is invalid: ", email)
				return
			}

			token := ""
			survey.AskOne(&survey.Input{
				Message: "input token",
			}, &token)

			addUser(userName, email, token)
		} else if userFunction == "deleteUser" {
			// TODO
		} else if userFunction == "resetUser" {
			resetUsers()
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

func resetUsers() {
	viper.Set("users", "")
	viper.WriteConfig()
	fmt.Println("All users are deleted from bit")
}

var userFunctions = []string{
	"addUser",
	"deleteUser",
	"resetUser",
	"listUser",
}
