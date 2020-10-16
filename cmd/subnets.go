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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/flippedbit/chaosity/internal/aws/networkacl"
	"github.com/flippedbit/chaosity/internal/aws/subnets"
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
		sess := session.Must(
			session.NewSession(&aws.Config{
				Region:      aws.String(o.Region),
				Credentials: credentials.NewSharedCredentials("", o.Profile),
			}),
		)
		svc := ec2.New(sess)

		var n string
		var err error
		var nAssoc []ec2.NetworkAclAssociation
		var doSomething bool = false

		s, err := subnets.GetSubnets(svc, o)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, sub := range s {
			fmt.Println("Found subnet: ", *sub.SubnetId)
			assoc, err := networkacl.GetNetworkAclAssociation(svc, *sub.SubnetId)
			if err != nil {
				fmt.Println(err)
				return
			}
			nAssoc = append(nAssoc, assoc[0])
		}

		if denyFlag {
			doSomething = true
			n, err = networkacl.CreateDenyNacl(svc, o.VpcID)
			if err != nil {
				fmt.Println(err)
				return
			}
			for i, a := range nAssoc {
				if new, err := networkacl.ReplaceAssociation(svc, *a.NetworkAclAssociationId, n); err != nil {
					fmt.Println(err)
				} else {
					nAssoc[i].NetworkAclAssociationId = &new
				}
			}
		}

		if doSomething {
			fmt.Println("Chaos! Waiting for ", o.Duration, " seconds...")
			time.Sleep(time.Duration(o.Duration) * time.Second)
		} else {
			fmt.Println("Chaos the chaos! Nothing to do so not going to wait...")
		}

		if denyFlag {
			for _, a := range nAssoc {
				if _, err := networkacl.ReplaceAssociation(svc, *a.NetworkAclAssociationId, *a.NetworkAclId); err != nil {
					fmt.Println(err)
				}
			}
			if err := networkacl.DeleteDenyNacl(svc, n); err != nil {
				fmt.Println(err)
				return
			}
		}
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
	subnetsCmd.Flags().BoolVarP(&denyFlag, "deny", "d", false, "Apply empty NetworkACL to deny all traffic to subnets.")
}
