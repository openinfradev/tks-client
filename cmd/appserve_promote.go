/*
Copyright © 2021 SK Telecom <https://github.com/openinfradev>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var appservePromoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "Promote an app paused in blue-green state",
	Long: `Promote an app paused in blue-green state.

Example:
tks appserve promote <APP_ID>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("APP_ID is mandatory! Run 'tks appserve promote --help'.")
		}

		// Get API Url
		appserveApiUrl := viper.GetString("tksAppServeLcmUrl")
		if appserveApiUrl == "" {
			return errors.New("tks_appserve_api_url is mandatory.")
		}

		/* Parse command line params */
		appId := args[0]

		// Prepare request body
		data := url.Values{}
		data.Set("promote", "true")

		// Initialize http client
		client := &http.Client{}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, appserveApiUrl+"/apps/"+appId, strings.NewReader(data.Encode()))
		if err != nil {
			return fmt.Errorf("Error while constructing req: %s", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("Error from update req: %s", err)
		}

		defer resp.Body.Close()

		// Check response
		respBody, err := io.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)
			// TODO: after test run, fix this output msg.
			fmt.Println("Response:\n", str)
		}

		return nil
	},
}

func init() {
	appserveCmd.AddCommand(appservePromoteCmd)
}
