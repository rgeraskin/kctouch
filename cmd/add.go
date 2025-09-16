/*
Copyright © 2025 Roman Geraskin

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/keybase/go-keychain"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "put", "set"},
	Short:   "Add a generic password item",
	Long:    "Add a new keychain entry with the specified service, account, label, and password.",
	Example: `  kctouch add -s "MyService" -a "johndoe"
  kctouch add -s "GitHub" -l "GitHub Token" -a "myusername" -p PLAIN_SECRET_PASS
  echo "mypassword" | kctouch add -s "this/is/a/service" -p -`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		// Authenticate before allowing access
		return auth(cmd, "add new keychain entry")
	},
	RunE: add,
}

func add(cmd *cobra.Command, args []string) error {
	service := cmd.Flag("service").Value.String()
	account := cmd.Flag("account").Value.String()
	label := cmd.Flag("label").Value.String()
	passwordFlag := cmd.Flag("password").Value.String()
	updateFlag, _ := cmd.Flags().GetBool("update")

	if service == "" {
		return fmt.Errorf("required flag(s) \"service\" not set")
	}

	password, err := getPassword(passwordFlag)
	if err != nil {
		return err
	}

	item := createKeychainItem(service, account, label, password)
	return addKeychainItem(item, service, label, account, updateFlag)
}

func getPassword(passwordFlag string) (string, error) {
	var password string

	switch passwordFlag {
	case "":
		// Prompt for password
		fmt.Print("Enter password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println() // Add newline after password input
		password = string(passwordBytes)
	case "-":
		// Read password from stdin
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			password = strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read password from stdin: %w", err)
		}
	default:
		password = passwordFlag
	}

	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	return password, nil
}

func createKeychainItem(service, account, label, password string) keychain.Item {
	// Create generic password item with service, account, label, password, access group
	item := keychain.NewGenericPassword(
		service,
		account,
		label,
		[]byte(password),
		accessGroup,
	)

	// don't know what this is for, it's from the example in the repo :)
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	return item
}

func addKeychainItem(item keychain.Item, service, label, account string, update bool) error {
	forMsg := composeForMsg(service, account, label)
	action := "added"

	err := keychain.AddItem(item)

	if err == keychain.ErrorDuplicateItem {
		if update {
			err = keychain.UpdateItem(item, item)
			action = "updated"
		} else {
			return fmt.Errorf(
				"keychain item already exists for %s, use --update to update it",
				forMsg,
			)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to add/update keychain item: %w", err)
	}

	fmt.Printf(
		"Successfully %s keychain entry for: %s\n",
		action,
		forMsg,
	)
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().
		StringP("password", "p", "", "Password (if omitted: will be prompted (safest way), use '-' for stdin)")
	addCmd.Flags().BoolP("update", "u", false, "Update existing keychain item if it exists already")
}
