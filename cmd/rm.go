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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"d", "del", "delete", "remove"},
	Short:   "Delete a generic password item",
	Long:    "Remove a keychain entry with the specified service and account.",
	Example: `  kctouch del -s "MyService"
  kctouch rm -s "GitHub" -a "myusername"`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		// Authenticate before allowing access
		return auth("remove keychain entry")
	},
	RunE: rm,
}

func rm(cmd *cobra.Command, args []string) error {
	service := cmd.Flag("service").Value.String()
	account := cmd.Flag("account").Value.String()

	if service == "" {
		return fmt.Errorf("required flag(s) \"service\" not set")
	}

	err := keychain.DeleteGenericPasswordItem(service, account)
	if err == keychain.ErrorItemNotFound {
		return fmt.Errorf("keychain item not found")
	}
	if err != nil {
		return fmt.Errorf("failed to delete keychain item: %w", err)
	}

	// label is ignored for deletion, so we pass an empty string
	forMsg := composeForMsg(service, account, "")

	fmt.Printf("Successfully deleted keychain entry for: %s\n", forMsg)
	return nil
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
