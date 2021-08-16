package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/chriswalz/complete/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// gitAccount represents the gitAccount command
var gitUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage git account",
	Long:  "Manage git account",
	Run: func(cmd *cobra.Command, args []string) {
		selectedSuggestion := ""
		userFunctionSuggestions := UserSuggestions()

		suggestionTree := &complete.CompTree{
			Sub: map[string]*complete.CompTree{
				"user": {
					Dynamic: toAutoCLI(userFunctionSuggestions),
				},
			},
		}
		selectedSuggestion = SuggestionPrompt("> bit user ", specificCommandCompleter("user", suggestionTree))

		if selectedSuggestion == "addUser" {
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
		} else if selectedSuggestion == "deleteUser" {
			deleteUser()
		} else if selectedSuggestion == "resetUser" {
			// TODO
		} else if selectedSuggestion == "listUser" {
			listUser()
		} else { // selected user
			users := readUsers()
			flag := selectFlag(users[selectedSuggestion])
			if user, ok := users[selectedSuggestion]; ok {
				applyGitAccount(user, flag)
			}
		}
	},
}

func init() {
	initLocalStorage()
	BitCmd.AddCommand(gitUserCmd)
}

func UserSuggestions() []complete.Suggestion {
	users := readUsers()
	userSuggestion := parseToSuggestion(users)

	var suggestions []complete.Suggestion
	for _, userFunction := range userFunctions {
		suggestions = append(suggestions, complete.Suggestion{
			Name: userFunction,
		})
	}

	return append(userSuggestion, suggestions...)
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
	writeUsers()
}

func listUser() {
	var users = viper.Get("users")

	fmt.Println("users: ", users)
}

func deleteUser() {
	users := readUsers()
	userSuggestion := parseToSuggestion(users)
	if len(userSuggestion) == 0 {
		fmt.Println("Has no account")
		return
	}

	deletingUsername := selectDeletingUser(userSuggestion)
	delete(viper.Get("users").(map[string]interface{}), deletingUsername)
	writeUsers()
}

func selectDeletingUser(userSuggestion []complete.Suggestion) string {
	suggestionTree := &complete.CompTree{
		Sub: map[string]*complete.CompTree{
			"deleteUser": {
				Dynamic: toAutoCLI(userSuggestion),
			},
		},
	}

	return SuggestionPrompt("> bit user deleteUser ", specificCommandCompleter("deleteUser", suggestionTree))
}

func selectFlag(user User) string {
	gitLocalConfigFile, err := os.Stat("./.git/config")
	if os.IsNotExist(err) || gitLocalConfigFile.IsDir() {
		return "--global"
	}

	suggestions := []complete.Suggestion{
		{Name: "--global", Desc: "set up global "},
		{Name: "--local", Desc: "set up local "},
	}

	suggestionTree := &complete.CompTree{
		Sub: map[string]*complete.CompTree{
			user.Name: {
				Dynamic: toAutoCLI(suggestions),
			},
		},
	}

	return SuggestionPrompt("> bit user "+user.Name+" ", specificCommandCompleter(user.Name, suggestionTree))
}

func applyGitAccount(user User, flag string) {
	if _, err := execCommand("git", "config", flag, "user.name", user.Name).CombinedOutput(); err != nil {
		log.Debug().Err(err).Send()
	} else {
		println("apply user name")
	}

	if _, err := execCommand("git", "config", flag, "user.email", user.Email).CombinedOutput(); err != nil {
		log.Debug().Err(err).Send()
	} else {
		println("apply user email")
	}

	if _, err := execCommand("git", "config", flag, "user.signingkey", user.Token).CombinedOutput(); err != nil {
		log.Debug().Err(err).Send()
	} else {
		println("apply user token")
	}
}

func readUsers() map[string]User {
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Read Error")
		return map[string]User{}
	}

	var users map[string]User
	err = viper.UnmarshalKey("users", &users)
	if err != nil {
		fmt.Println("Get Error")
		return map[string]User{}

	}

	return users
}

func writeUsers() {
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println("Write Error")
	}
}

func parseToSuggestion(userMap map[string]User) []complete.Suggestion {
	var userList []complete.Suggestion
	for _, user := range userMap {
		userList = append(userList, complete.Suggestion{Name: "" + user.Name, Desc: user.Email})
	}

	return userList
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

var userFunctions = []string{
	"addUser",
	"deleteUser",
	"resetUser",
	"listUser",
}
