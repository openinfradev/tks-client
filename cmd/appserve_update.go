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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
)

var appserveUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an app deployed by AppServing service",
	Long: `Update an app deployed by AppServing service.
  
Example:
tks appserve update <APP_ID> --appserve-config <CONFIGFILE>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("APP_ID is mandatory! Run 'tks appserve update --help'.")
		}

		// Get API Url
		appserveApiUrl := viper.GetString("tksAppServeLcmUrl")
		if appserveApiUrl == "" {
			return errors.New("tks_appserve_api_url is mandatory.")
		}

		var c conf
		if appserveCfgFile == "" {
			return errors.New("--appservce-config is mandatory param.")
		}

		/* Parse command line params */
		appId := args[0]

		// Get Appserving request params from config file
		yamlData, err := os.ReadFile(appserveCfgFile)
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}

		fmt.Printf("*******\nConfig:\n%+s\n*******\n", yamlData)

		// Get application config from file
		if appCfgFile != "" {
			appCfgBytes, err := os.ReadFile(appCfgFile)
			if err != nil {
				return fmt.Errorf("error: %s", err)
			}
			// Add appCfg to existing struct
			c.App_config = string(appCfgBytes)
		}

		// Get application secret from file
		if appSecretFile != "" {
			appSecretBytes, err := os.ReadFile(appSecretFile)
			if err != nil {
				return fmt.Errorf("error: %s", err)
			}
			c.App_secret = string(appSecretBytes)
		}

		// Unmarshal yaml content into struct
		err = yaml.Unmarshal(yamlData, &c)
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}

		// Exclude immutable params
		if c.Name != "" || c.Type != "" || c.App_type != "" {
			fmt.Println(`Following params are immutable in update action, so they're ignored.
- name
- type
- app_type
- target_cluster`)
			c.Name = ""
			c.Type = ""
			c.App_type = ""
			c.Target_cluster_id = ""
		}

		// Convert map to Json
		cBytes, err := json.Marshal(&c)
		if err != nil {
			return fmt.Errorf("Unable to marshal config to JSON: %s", err)
		}

		fmt.Println("========== ")
		fmt.Println("Json data: ")
		fmt.Println(string(cBytes))
		fmt.Println("========== ")

		buff := bytes.NewBuffer(cBytes)

		// Initialize http client
		client := &http.Client{}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, appserveApiUrl+"/apps/"+appId, buff)
		if err != nil {
			return fmt.Errorf("Error while constructing req: %s", err)
		}

		// set the request header Content-Type for json
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("Error from update req: %s", err)
		}

		defer resp.Body.Close()

		// Check response
		respBody, err := io.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)
			fmt.Println("Response:\n", str)
		}

		return nil
	},
}

func init() {
	appserveCmd.AddCommand(appserveUpdateCmd)

	appserveUpdateCmd.Flags().StringVar(&appserveCfgFile, "appserve-config", "", "config file for AppServing service")
	appserveUpdateCmd.Flags().StringVar(&appCfgFile, "app-config", "", "custom config file for user application")
	appserveUpdateCmd.Flags().StringVar(&appSecretFile, "app-secret", "", "custom secret file for user application")
}
