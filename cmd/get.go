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
	"fmt"

	"github.com/keybase/go-keychain"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g", "find"},
	Short:   "Get a generic password item",
	Long:    "Get a keychain entry with the specified service, account, and label.",
	Example: `  kctouch get -s "MyService" -a "johndoe"
  kctouch g -s "GitHub" -l "GitHub Token" -a "myusername"
  kctouch find -s "this/is/a/service"`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		// Authenticate before allowing access
		return auth("get keychain entry")
	},
	RunE: get,
}

func get(cmd *cobra.Command, args []string) error {
	service := cmd.Flag("service").Value.String()
	account := cmd.Flag("account").Value.String()
	label := cmd.Flag("label").Value.String()

	passwordBytes, err := keychain.GetGenericPassword(service, account, label, accessGroup)
	if err != nil {
		return fmt.Errorf("failed to get keychain item: %w", err)
	}
	password := string(passwordBytes)

	if password == "" {
		return fmt.Errorf("no password found or password is empty")
	}

	fmt.Println(password)
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)
}
