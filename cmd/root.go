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
	"os"

	"github.com/ansxuman/go-touchid"
	"github.com/spf13/cobra"
)

const (
	authMethod  = touchid.DeviceTypeAny
	accessGroup = ""
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kctouch",
	Short: "A tool for managing keychain items with TouchID",
	Long:  "A tool for managing generic-password macOS Keychain items with TouchID authentication.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("account", "a", "", "Account name")
	rootCmd.PersistentFlags().StringP("service", "s", "", "Service name (required for add, rm)")
	rootCmd.PersistentFlags().
		StringP("label", "l", "", "Label (if omitted, service name will be used)")
}

func auth(reason string) error {
	// Authenticate using any available method (Touch ID or passcode)
	success, err := touchid.Auth(authMethod, reason)
	if err != nil {
		return fmt.Errorf("authentication error: %v", err)
	}

	if !success {
		return fmt.Errorf("authentication failed")
	}

	return nil
}

func composeForMsg(service, account, label string) string {
	forMsg := fmt.Sprintf("service='%s'", service)
	if account != "" {
		forMsg += fmt.Sprintf(" account='%s'", account)
	}
	if label != "" {
		forMsg += fmt.Sprintf(" label='%s'", label)
	}
	return forMsg
}
