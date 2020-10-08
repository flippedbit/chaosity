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
	internalAWS "github.com/flippedbit/chaosity/internal/aws"
	"github.com/spf13/cobra"
)

var rebootFlag bool
var denyFlag bool
var shutdownFlag bool

// instancesCmd represents the instances command
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Apply an empty security group to all instances within an availability zone or subnets",
	Long: `Generates an emtpy security group in order to deny all traffic.
All instances are then gathered within the provided subnets or availability zone given.
If the subnets parameter is used, individual subnet IDs should be provided comma delimited
and they will be separated out. After the duration period, provided by the --duration flag,
all previous security groups are re-applied to each instance. Finally, the empty secuirty
group is deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		doSomething := false
		sess := session.Must(
			session.NewSession(&aws.Config{
				Region:      aws.String(options.region),
				Credentials: credentials.NewSharedCredentials("", options.profile),
			}),
		)
		svc := ec2.New(sess)

		var instances []*ec2.Instance
		var denySG string

		if options.subnets != "" {
			instances, _ = internalAWS.GetInstancesBySubnet(svc, &options.subnets)
		}
		if denyFlag {
			doSomething = true
			denySG, err := internalAWS.GenerateDenySecurityGroup(svc, &options.vpcID)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Created SecurityGroup ", denySG)
			if err := internalAWS.ApplyChaosSecurityGroupToInstances(svc, instances, denySG); err != nil {
				fmt.Println(err)
				return
			}
		}
		if rebootFlag {
			//doSomething = true
			internalAWS.RebootInstances(svc, instances)
		} else if shutdownFlag {
			//doSomething = true
			internalAWS.ForceShutdownInstances(svc, instances)
		}
		if doSomething {
			fmt.Println("Chaos! Waiting for ", options.duration, " seconds...")
			time.Sleep(time.Duration(options.duration) * time.Second)
		} else {
			fmt.Println("Chaos the chaos! Nothing to do so not going to wait...")
		}
		if denyFlag {
			if err := internalAWS.RevertChaosSecurityGroupOnInstances(svc, instances); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Deleting SecurityGroup ", denySG)
			internalAWS.DeleteDenySecurityGroup(svc, denySG)
		}
	},
}

func init() {
	awsCmd.AddCommand(instancesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// instancesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// instancesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	instancesCmd.Flags().BoolVarP(&rebootFlag, "reboot", "r", false, "Reboot selected instances from subnets or availability-zone.")
	instancesCmd.Flags().BoolVarP(&denyFlag, "deny", "d", false, "Apply deny security group to instances.")
	instancesCmd.Flags().BoolVarP(&shutdownFlag, "shutdown", "s", false, "Force stop selected instances from subnets or availability-zone.")
	instancesCmd.MarkFlagRequired("profile")
	instancesCmd.MarkFlagRequired("vpc-id")
	instancesCmd.MarkFlagRequired("region")
}
