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
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// showBYOHAgentGuideCmd represents command for showing BYOH agent installation guide.
var showByohAgentGuideCmd = &cobra.Command{
	Use:   "show-byoh-node-agent-guide",
	Short: "Show BYOH node agent installation guide.",
	Long: `Show BYOH node agent installation guide.

Example:
$ tks cluster show-byoh-node-agent-guide --type=$NODE_TYPE

NODE_TYPE should be one of these: ['controlplane', 'tks', 'worker']
Standard reference architecture is as follows.
 - controlplane (for k8s controlplane) x 3
 - tks (for tks built-in module such as LMA) x 3
 - worker (for user applications) x 3

Among these types, the 'worker' nodes might needs to be scaled out based on your application size.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		nodeType, _ := cmd.Flags().GetString("type")
		if nodeType == "" {
			return fmt.Errorf("Usage: tks cluster show-byoh-node-agent-guide --type=$NODE_TYPE\n")
		}

		if nodeType != "controlplane" && nodeType != "tks" && nodeType != "worker" {
			return fmt.Errorf("Wrong node type '%s': please refer to help message.\n", nodeType)
		}

		fmt.Printf("*************************************************\n")
		fmt.Printf("******** BYOH Agent Installation Process ********\n")
		fmt.Printf("*************************************************\n\n")
		fmt.Printf("Follow these steps to install BYOH agent on your machine.\n\n")

		// Print mgmt cluster's kubeconfig
		cmdStr := "cat ~/.kube/config | base64"
		out, err := exec.Command("bash", "-c", cmdStr).Output()
		if err != nil {
			return fmt.Errorf("Error: %s", err)
		}

		fmt.Printf("Encoded kubeconfig for MGMT cluster:\n%s\n", string(out))

		// Print steps for installing the agent
		guide_str := `# Decode the above string using base64 command as follows.
$ echo $ENCODED_CONFIG | base64 -d > mgmt-cluster.conf

# Install essential packages
$ sudo apt-get install socat ebtables ethtool conntrack

# Make sure hostname is listed in hostsfile
$ echo "127.0.0.1 $(hostname)" >> /etc/hosts

# Download agent
$ curl -LO https://github.com/openinfradev/tks-file-repo/releases/download/stable/byoh-hostagent-linux-amd64
$ chmod a+x ./byoh-hostagent-linux-amd64

# Set proper label depending on the node role
LABEL_OPT="--label role=%s"

# Run agent (NOTE: THE 'namespace' param MUST be substituted to your cluster uuid!!)
sudo killall byoh-hostagent-linux-amd64 || true
sudo ./byoh-hostagent-linux-amd64 --kubeconfig mgmt-cluster.conf --namespace $YOUR_CLUSTER_UUID $LABEL_OPT --v 20 2>&1 | tee byoh-agent.log

# Install binary bundle
$ curl -L -o k8s-v1.22.3-bundle.tar.xz https://github.com/openinfradev/tks-file-repo/releases/download/stable/k8s-v1.22.3-bundle.tar.xz
$ sudo mkdir -p /var/lib/byoh/bundles/projects.registry.vmware.com.cluster_api_provider_bringyourownhost
$ sudo tar xvfJ k8s-v1.22.3-bundle.tar.xz -C /var/lib/byoh/bundles/projects.registry.vmware.com.cluster_api_provider_bringyourownhost

That's it! Enjoy BYOH provider!
*******************************
`
		fmt.Printf(guide_str, nodeType)
		return nil
	},
}

func init() {
	clusterCmd.AddCommand(showByohAgentGuideCmd)
	showByohAgentGuideCmd.Flags().String("type", "", "[mandatory] node type in <controlplane|tks|worker>")
}
