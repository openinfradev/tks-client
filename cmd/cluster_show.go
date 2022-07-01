/*
Copyright © 2022 SK Telecom <https://github.com/openinfradev>

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
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jedib0t/go-pretty/table"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// clusterShowCmd represents the list command
var clusterShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cluster details.",
	Long: `Show cluster details. 

Example:
tks cluster show <CLUSTER_ID>`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("You must specify cluster id.")
			return errors.New("Usage: tks cluster show <CLUSTER_ID>")
		}
		var conn *grpc.ClientConn
		tksInfoUrl = viper.GetString("tksInfoUrl")
		if tksInfoUrl == "" {
			return errors.New("You must specify tksInfoUrl at config file")
		}
		conn, err := grpc.Dial(tksInfoUrl, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewClusterInfoServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		data := pb.GetClusterRequest{}
		data.ClusterId = args[0]

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
		r, err := client.GetCluster(ctx, &data)
		if err != nil {
			return fmt.Errorf("Error: %s", err)
		} else {
			printCluster(r)
		}
		return nil
	},
}

func init() {
	clusterCmd.AddCommand(clusterShowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterShowCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

func printCluster(r *pb.GetClusterResponse) {
	c := r.Cluster
	t := table.NewWriter()

	t.AppendHeader(table.Row{"Filed", "Value"})
	tCreatedAt := parseTime(c.CreatedAt)
	tUpdatedAt := parseTime(c.UpdatedAt)
	t.AppendRow(table.Row{"Name", c.Name})
	t.AppendRow(table.Row{"ID", c.Id})
	t.AppendRow(table.Row{"Created At", tCreatedAt})
	t.AppendRow(table.Row{"Updated At", tUpdatedAt})
	t.AppendRow(table.Row{"Status", c.Status})
	t.AppendRow(table.Row{"Contract ID", c.ContractId})
	t.AppendRow(table.Row{"Csp ID", c.CspId})
	t.AppendRow(table.Row{"Conf", c.Conf})
	t.AppendRow(table.Row{"AppGroups", c.AppGroups})
	t.AppendRow(table.Row{"Kubeconfig", c.Kubeconfig})

	fmt.Println(t.Render())
}
