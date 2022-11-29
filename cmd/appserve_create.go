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
	"os"
	"net/http"
)

var appserveCfgFile string

// Member variables are named as snake_case on purpose
// to be marshalled into json object later.
type conf struct {
	Contract_id  string `yaml:"contract_id"`
	Name         string `yaml:"name"`
	Version      string `yaml:"version"`
	Task_type    string `yaml:"task_type"`
	Artifact_url string `yaml:"artifact_url"`
	Port         string `yaml:"port"`
	Profile      string `yaml:"profile"`

	Resource_spec     string `yaml:"resource_spec"`
	Target_cluster_id string `yaml:"target_cluster_id"`

	Pv_enabled       bool   `yaml:"pv_enabled"`
	Pv_storage_class string `yaml:"pv_storage_class"`
	Pv_access_mode   string `yaml:"pv_access_mode"`
	Pv_size          string `yaml:"pv_size"`
	Pv_mount_path    string `yaml:"pv_mount_path"`
}

var appserveCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an app by AppServing service",
	Long: `Create an app by AppServing service.
  
Example:
tks appserve create --appserve-config CONFIGFILE`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get API Url
		appserveApiUrl := viper.GetString("tksAppServeLcmUrl")
		if appserveApiUrl == "" {
			return errors.New("tks_appserve_api_url is mandatory.")
		}

		var c conf
		if appserveCfgFile == "" {
			return errors.New("--appservce-config is mandatory param.")
		}

		// Get Appserving request params from config file
		yamlData, err := os.ReadFile(appserveCfgFile)
		if err != nil {
			return fmt.Errorf("error: %s", err)
		}

		fmt.Printf("*******\nConfig:\n%+s\n*******\n", yamlData)

		err = yaml.Unmarshal(yamlData, &c)
		if err != nil {
			return fmt.Errorf("error: %s", err)
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
		resp, err := http.Post(appserveApiUrl+"/apps", "application/json", buff)
		if err != nil {
			return fmt.Errorf("Error from POST req: %s", err)
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
	appserveCmd.AddCommand(appserveCreateCmd)

	appserveCreateCmd.Flags().StringVar(&appserveCfgFile, "appserve-config", "", "config file for AppServing service")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appserveCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appserveCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
