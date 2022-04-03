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

	"github.com/jedib0t/go-pretty/table"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
)

// serviceListCmd represents the create command
var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show list of service.",
	Long: `Show list of service.

Example:
tks service list <CLUSTER ID> (--long)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify cluster ID.")
			fmt.Println("Usage: tks service list <CLUSTER ID>")
			os.Exit(1)
		}
		var conn *grpc.ClientConn
		tksInfoUrl = viper.GetString("tksInfoUrl")
		if tksInfoUrl == "" {
			fmt.Println("You must specify tksInfoUrl at config file")
			os.Exit(1)
		}
		conn, err := grpc.Dial(tksInfoUrl, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewAppInfoServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		ClusterId := args[0]

		data := pb.IDRequest{}
		data.Id = ClusterId

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
		r, err := client.GetAppGroupsByClusterID(ctx, &data)
		if err != nil {
			fmt.Println(err)
		} else {
			long, _ := cmd.Flags().GetBool("long")
			printAppGroups(r, long)
		}
	},
}

func init() {
	serviceCmd.AddCommand(serviceListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serviceListCmd.Flags().BoolP("long", "l", false, "Print detail information")
}

func printAppGroups(r *pb.GetAppGroupsResponse, long bool) {
	t := table.NewWriter()
	tTemp := table.Table{}
	tTemp.Render()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	if long {
		t.AppendHeader(table.Row{"TYPE", "SERVICE_ID", "STATUS", "CREATED_AT", "UPDATED_AT", "UPDATE_DESC"})
		for _, s := range r.AppGroups {
			tCreatedAt := parseTime(s.CreatedAt)
			tUpdatedAt := parseTime(s.UpdatedAt)

			t.AppendRow(table.Row{s.Type, s.AppGroupId, s.Status, tCreatedAt, tUpdatedAt, s.StatusDesc})
		}
	} else {
		t.AppendHeader(table.Row{"TYPE", "SERVICE_ID", "STATUS", "CREATED_AT", "UPDATED_AT"})
		for _, s := range r.AppGroups {
			tCreatedAt := parseTime(s.CreatedAt)
			tUpdatedAt := parseTime(s.UpdatedAt)

			t.AppendRow(table.Row{s.Type, s.AppGroupId, s.Status, tCreatedAt, tUpdatedAt})
		}
	}

	if len(r.AppGroups) > 0 {
		fmt.Println(t.Render())
	} else {
		fmt.Println("No services found.")
	}
}
