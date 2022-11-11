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
	"io/ioutil"
	"net/http"
)

var appserveCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an app by AppServing service",
	Long: `Create an app by AppServing service.
  
Example:
tks appserve create <APPNAME> [--config CONFIGFILE]`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		//if len(args) == 0 {
		//	return errors.New("Usage: tks appserve create <APPNAME>")
		//}

		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			viper.SetConfigName("appserve-config") // name of config file (without extension)
			viper.SetConfigType("yaml")            // REQUIRED if the config file does not have the extension in the name
			viper.AddConfigPath(".")               // optionally look for config in the working directory
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				return errors.New("Config file not found. Aborting..")
			} else {
				// Config file was found but another error was produced
				return fmt.Errorf("Error while reading config: %s", err)
			}
		} else {
			fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
		}

		appserveApiUrl := viper.GetString("tks_appserve_api_url")
		if appserveApiUrl == "" {
			return errors.New("tks_appserve_api_url is mandatory in config file")
		}

		c := viper.AllSettings()
		//fmt.Printf("viper map: %v\n\n", c)
		delete(c, "tks_appserve_api_url")
		cBytes, err := json.Marshal(c) // 맵을 JSON 문서로 변환
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
		respBody, err := ioutil.ReadAll(resp.Body)
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

	appserveCreateCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is ./appserve-config.yaml)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appserveCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appserveCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
