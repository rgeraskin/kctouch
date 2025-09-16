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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ansxuman/go-touchid"
	"github.com/keybase/go-keychain"
	"github.com/spf13/cobra"
)

const (
	authMethod  = touchid.DeviceTypeAny
	accessGroup = ""
	cacheKey    = "/kctouch/cache"
)

type cacheEntry struct {
	ForAttempts int       `json:"forAttempts"`
	ForTime     time.Time `json:"forTime"`
}

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
	// use strings for auth cache because want to handle empty values
	rootCmd.PersistentFlags().
		String("cache-for", "", "Cache auth for time (e.g. 1h, 10m, 10s), use 0 to invalidate")
	rootCmd.PersistentFlags().
		String("cache-n", "", "Cache auth for N auth requests, use 0 to invalidate")

	// Add verbose flag
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")

	// Set up logging based on verbose flag
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		if !verbose {
			log.SetOutput(io.Discard)
		}
	}
}

func updateCache(
	cache *cacheEntry,
	account string,
	cacheForStr string,
	cacheAttemptsStr string,
) error {
	var err error
	log.Printf("updating cache for account '%s': %+v", account, cache)

	if cacheAttemptsStr != "" {
		cache.ForAttempts, err = strconv.Atoi(cacheAttemptsStr)
		if err != nil {
			return fmt.Errorf("failed to convert --cache-n to int: %w", err)
		}
	}

	if cacheForStr != "" {
		cacheFor, err := time.ParseDuration(cacheForStr)
		if err != nil {
			return fmt.Errorf("failed to convert --cache-for to duration: %w", err)
		}
		cache.ForTime = time.Now().Add(cacheFor)
	}

	cacheBytes, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal kctouch json: %w", err)
	}

	item := createKeychainItem(cacheKey, account, cacheKey, string(cacheBytes))

	err = keychain.AddItem(item)
	if err == keychain.ErrorDuplicateItem {
		err = keychain.UpdateItem(item, item)
	}

	if err != nil {
		return fmt.Errorf("failed to add/update keychain item: %w", err)
	}

	return nil
}

func getCache(account string) (*cacheEntry, error) {
	var cache cacheEntry
	cacheBytes, err := keychain.GetGenericPassword(cacheKey, account, cacheKey, accessGroup)

	// something went wrong
	if err != nil {
		return nil, fmt.Errorf("failed to get cached auth: %w", err)
	}

	// if cache bytes are nil, it's ok
	if cacheBytes == nil {
		log.Printf("cached auth is empty for account '%s'", account)
		return &cache, nil
	}

	// decode cache json
	err = json.Unmarshal(cacheBytes, &cache)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal kctouch json: %w", err)
	}

	log.Printf("found cached auth for account '%s': %+v", account, cache)
	return &cache, nil
}

func auth(cmd *cobra.Command, reason string) error {
	cacheForStr := cmd.Flag("cache-for").Value.String()
	cacheAttemptsStr := cmd.Flag("cache-n").Value.String()
	account := cmd.Flag("account").Value.String()

	// try to retrieve cached auth
	cache, err := getCache(account)
	if err != nil {
		return fmt.Errorf("failed to get cached auth: %w", err)
	}

	switch {
	// auth time is not expired
	case cache.ForTime.After(time.Now()):
		log.Printf("cached auth time is not expired for account '%s'", account)

	// auth attempts are not expired
	case cache.ForAttempts > 0:
		log.Printf("cached auth attempts are not expired for account '%s'", account)
		// decrement auth attempts if cacheN is not set
		if cacheAttemptsStr == "" {
			cacheAttemptsStr = fmt.Sprintf("%d", cache.ForAttempts-1)
		}
	// have to auth
	default:
		log.Printf("have to auth for account '%s'", account)
		success, err := touchid.Auth(authMethod, reason)
		if err != nil {
			return fmt.Errorf("authentication error: %v", err)
		}

		if !success {
			return fmt.Errorf("authentication failed")
		}
	}

	// update cache if required
	if cacheForStr != "" || cacheAttemptsStr != "" {
		err = updateCache(cache, account, cacheForStr, cacheAttemptsStr)
		if err != nil {
			return fmt.Errorf("failed to update cache: %w", err)
		}
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
