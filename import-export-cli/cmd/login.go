package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var loginUsername string
var loginPassword string
var loginPasswordStdin bool

const loginCmdLiteral = "login [environment] [flags]"
const loginCmdShortDesc = "Login to an API Manager"
const loginCmdLongDesc = `Login to an API Manager using credentials`
const loginCmdExamples = utils.ProjectName + " login dev -u admin -p admin\n" +
	utils.ProjectName + " login dev -u admin\n" +
	"cat ~/.mypassword | " + utils.ProjectName + " login dev -u admin"

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     loginCmdLiteral,
	Short:   loginCmdShortDesc,
	Long:    loginCmdLongDesc,
	Example: loginCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environment := args[0]

		if loginPassword != "" {
			fmt.Println("Warning: Using --password in CLI is not secure. Use --password-stdin")
			if loginPasswordStdin {
				fmt.Println("--password and --password-stdin are mutual exclusive")
				os.Exit(1)
			}
		}

		if loginPasswordStdin {
			if loginUsername == "" {
				fmt.Println("An username is required to use password-stdin")
				os.Exit(1)
			}

			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			loginPassword = strings.TrimRight(strings.TrimSuffix(string(data), "\n"), "\r")
		}

		err := runLogin(environment, loginUsername, loginPassword)
		if err != nil {
			fmt.Println("Error occurred while login : ", err)
			os.Exit(1)
		}
	},
}

func runLogin(environment, username, password string) error {
	if !utils.EnvExistsInMainConfigFile(environment, utils.MainConfigFilePath) {
		fmt.Println(environment, "does not exists. Add it using add-env")
		os.Exit(1)
	}

	if username == "" {
		fmt.Print("Username:")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			username = scanner.Text()
		}
	}

	if password == "" {
		fmt.Print("Password:")
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return err
		}
		password = string(pass)
		fmt.Println()
	}

	store, err := credentials.GetDefaultCredentialStore()
	if err != nil {
		return err
	}

	registrationEndpoint := utils.GetRegistrationEndpointOfEnv(environment, utils.MainConfigFilePath)
	clientId, clientSecret, err := utils.GetClientIDSecret(username, password, registrationEndpoint)
	if err != nil {
		return err
	}

	fmt.Println("Logged into", environment, "environment")
	err = store.Set(environment, username, password, clientId, clientSecret)
	if err != nil {
		return err
	}

	return nil
}

// init using Cobra
func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Username for login")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password for login")
	loginCmd.Flags().BoolVarP(&loginPasswordStdin, "password-stdin", "", false, "Get password from stdin")
}
