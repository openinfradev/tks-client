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
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/openinfradev/tks-proto/tks_pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clusterImportCmd represents the import command
var clusterImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a TKS Cluster.",
	Long: `Import a TKS Cluster.
  
Example:
tks cluster import <CLUSTERNAME> [--contract-id CONTRACTID --kubeconfig KUBECONFIG]`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("You must specify cluster name.")
			return errors.New("Usage: tks cluster import <CLUSTERNAME> --contract-id <CONTRACTID>")
		}
		var conn *grpc.ClientConn
		tksClusterLcmUrl = viper.GetString("tksClusterLcmUrl")
		if tksClusterLcmUrl == "" {
			return errors.New("You must specify tksClusterLcmUrl at config file")
		}
		conn, err := grpc.Dial(tksClusterLcmUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Could not connect to LCM server: %s", err)
		}
		defer conn.Close()

		client := pb.NewClusterLcmServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		/* Parse command line arguments */
		ClusterName := args[0]
		ContractId, _ := cmd.Flags().GetString("contract-id")
		creator, _ := cmd.Flags().GetString("creator")
		description, _ := cmd.Flags().GetString("description")
		templateName, _ := cmd.Flags().GetString("template")
		kubeconfigPath, _ := cmd.Flags().GetString("kubeconfig-path")
		kubeconfig, err := os.ReadFile(kubeconfigPath)
		if err != nil {
			log.Fatalf("Failed to read kubeconfig from [%s] path", err)
			log.Fatalf("Failed to read kubeconfig from [%s] path", kubeconfigPath)
		}

		/* Construct request map */
		data := pb.ImportClusterRequest{
			Name:         ClusterName,
			ContractId:   ContractId,
			Kubeconfig:   kubeconfig,
			TemplateName: templateName,
			Creator:      creator,
			Description:  description,
		}

		/* Comment for security
		m := protojson.MarshalOptions{
			Indent:        "  ",
			UseProtoNames: true,
		}
		jsonBytes, _ := m.Marshal(&data)
		fmt.Println("Proto Json data: ")
		fmt.Println(string(jsonBytes))
		*/

		r, err := client.ImportCluster(ctx, &data)
		fmt.Println("Response:\n", r)
		if err != nil {
			return fmt.Errorf("Error: %s", err)
		} else {
			fmt.Println("Success: The request to import cluster ", args[0], " was accepted.")
		}

		return nil
	},
}

func init() {
	clusterCmd.AddCommand(clusterImportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	clusterImportCmd.Flags().String("contract-id", "", "Contract ID")
	clusterImportCmd.Flags().String("kubeconfig-path", "", "Path of Kubeconfig for importing cluster")
	_ = clusterImportCmd.MarkFlagRequired("kubeconfig-path")
	clusterImportCmd.Flags().String("template", "aws-reference", "Template name for the cluster")

	clusterImportCmd.Flags().String("creator", "", "Uuid of creator")
	clusterImportCmd.Flags().String("description", "", "Description of cluster")
}
