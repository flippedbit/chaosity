/*
Copyright © 2020 Michael Straughan <straughan.michael@gmail.com>

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

	"github.com/spf13/cobra"
)

// azCmd represents the az command
var azCmd = &cobra.Command{
	Use:   "az",
	Short: "Block all traffic to an AWS availability zone usign a NetworkACL",
	Long: `Generates an empty NetworkACL within AWS which by default will block all traffic.
NACLs are able to block all traffic that comes in to or out of the subnet but will not
block any traffic that passes within the subnet itself. All subnets are gathered within the
availability zone provided using the --availibility-zone flag. That empty NACL is
then applied to subnets previously gathered. After the duration period, given using the --duration flag,
the original NACLs are then re-applied to the subnets that were previously gathered
and the empty NACL is deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("az called")
	},
}

func init() {
	awsCmd.AddCommand(azCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// azCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// azCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
