package ssm

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type commandStruct struct {
	command   string
	profile   string
	ssmClient *ssm.SSM
	commandID string
	os        string
}

var c commandStruct

// Checks instance to ensure it is online and rectrieve its OS type
func SendCommandToSSM(ssmClient *ssm.SSM, instances []string, mode string) {
	//Loop Through Instances
	//Check Instance
	//SendCommand
	//Poll Command
	//Return Result
	c.ssmClient = ssmClient
	for _, instance := range instances {
		fmt.Println(instance)
	}
}
func (c *commandStruct) checkInstance(instance string) (os string) {
	instanceData, err := c.ssmClient.DescribeInstanceInformation(&ssm.DescribeInstanceInformationInput{
		Filters: []*ssm.InstanceInformationStringFilter{
			{
				Key:    aws.String("InstanceIds"),
				Values: []*string{aws.String(instance)},
			},
		},
	})
	if err != nil {
		log.Fatalln("Unable to retrieve instance information please ensure your instanceid exists within the given profile.")

	}
	if len(instanceData.InstanceInformationList) > 0 {
		status := instanceData.InstanceInformationList[0].PingStatus
		os := *instanceData.InstanceInformationList[0].PlatformType

		if *status == "Online" {
			return os
		}

	}
	if len(instanceData.InstanceInformationList) == 0 {
		log.Fatalln("Couldn't Find Instance in SSM.. Ensure ssm-agent is installed on the instance provided.")
		return ""
	}

	return os

}
