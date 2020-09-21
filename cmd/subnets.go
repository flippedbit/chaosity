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

// subnetsCmd represents the subnets command
var subnetsCmd = &cobra.Command{
	Use:   "subnets",
	Short: "A brief description of your command",
	Long: `Generates an empty NetworkACL within AWS which by default will block all traffic.
	NACLs are able to block all traffic that comes in to or out of the subnet but will not
	block any traffic that passes within the subnet itself. That empty NACL is then applied
	to all subnets that are passed using the --subnets parameter comma delimited.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subnets called")
	},
}

func init() {
	awsCmd.AddCommand(subnetsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subnetsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subnetsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
