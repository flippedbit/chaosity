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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

type awsOptions struct {
	region   string
	vpcID    string
	profile  string
	subnets  string
	az       string
	duration int
}

var options awsOptions

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Test programmatic connectivity to AWS",
	Long: `This module will only perform a simple AWS API connection to confirm that you are able to connect with the given credentials.
	You are required to provide --profile --vpc-id and --region parameters.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws called")
		fmt.Println(&options)
		sess := session.Must(
			session.NewSession(&aws.Config{
				Region:      aws.String(options.region),
				Credentials: credentials.NewSharedCredentials("", options.profile),
			}),
		)
		svc := ec2.New(sess)
		if *svc.Config.Region == options.region {
			fmt.Println("Able to connect to AWS EC2 - profile: ", options.profile, "\tregion: ", options.region)
		} else {
			fmt.Println("Did not get proper response back from AWS")
		}
	},
}

func init() {
	rootCmd.AddCommand(awsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// awsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	awsCmd.PersistentFlags().StringVar(&options.region, "region", "", "AWS region to perform chaos in (required)")
	awsCmd.PersistentFlags().StringVar(&options.vpcID, "vpc-id", "", "AWS VPC to perform chaos in (required)")
	awsCmd.PersistentFlags().StringVar(&options.profile, "profile", "", "AWS credentials profile to use in order to connect (required)")
	awsCmd.PersistentFlags().StringVar(&options.subnets, "subnets", "", "AWS Subnet IDs to perform chaos on (comma separated)")
	awsCmd.PersistentFlags().StringVar(&options.az, "availability-zone", "", "AWS Availibility-Zone to perform chaos on.")
	awsCmd.PersistentFlags().IntVar(&options.duration, "duration", 300, "How long to perform chaos testing for in seconds")

	awsCmd.MarkFlagRequired("profile")
	awsCmd.MarkFlagRequired("vpc-id")
	awsCmd.MarkFlagRequired("region")
}
