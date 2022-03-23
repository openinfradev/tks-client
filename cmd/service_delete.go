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
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/openinfradev/tks-proto/tks_pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

// serviceDeleteCmd represents the create command
var serviceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a TKS Service.",
	Long: `Delete a TKS Service.

Example:
tks service delete <SERVICE ID>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify service ID.")
			fmt.Println("Usage: tks service delete <SERVICE ID>")
			os.Exit(1)
		}
		var conn *grpc.ClientConn
		tksClusterLcmUrl = viper.GetString("tksClusterLcmUrl")
		if tksClusterLcmUrl == "" {
			fmt.Println("You must specify tksClusterLcmUrl at config file")
			os.Exit(1)
		}
		conn, err := grpc.Dial(tksClusterLcmUrl, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewClusterLcmServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		serviceId := args[0]

		data := pb.UninstallAppGroupsRequest{}
		data.AppGroupIds = []string{serviceId}

		m := protojson.MarshalOptions{
			Indent:        "  ",
			UseProtoNames: true,
		}
		jsonBytes, _ := m.Marshal(&data)
		verbose, err := rootCmd.PersistentFlags().GetBool("verbose")
		if verbose {
			fmt.Println("Proto Json data...")
			fmt.Println(string(jsonBytes))
		}
		r, err := client.UninstallAppGroups(ctx, &data)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r)
		}
	},
}

func init() {
	serviceCmd.AddCommand(serviceDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
