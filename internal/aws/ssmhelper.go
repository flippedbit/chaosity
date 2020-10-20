package aws

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type commandStruct struct {
	command   string
	profile   string
	ssmClient *ssm.SSM
	commandID string
	os        string
	encode    bool
}

var c commandStruct

//Checks instance to ensure it is online and rectrieve its OS type
func SendCommandToSSM(ssmClient *ssm.SSM, instances []*ec2.Instance, mode string) {
	c.ssmClient = ssmClient
	for _, instance := range instances {
		//error handle this
		log.Printf("Getting SSM Information about instanceID: %v", *instance.InstanceId)
		c.os = c.checkInstance(*instance.InstanceId)
		if c.os == "" {
			log.Fatalln("Unable to determine information from SSM about the selected instance")
		}
		if mode == "stop" && c.os != "Windows" {
			c.command = "sudo dhclient -r"
		} else if mode == "stop" && c.os == "Windows" {
			c.command = "ipconfig /release"
		}
		log.Printf("Got OS: %v for instance: %v", c.os, *instance.InstanceId)
		c.runCommand(*instance.InstanceId)
		if mode != "stop" {
			log.Printf("Polling command instanceId: %v", *instance.InstanceId)
			c.pollCommand(*instance.InstanceId)
		} else {
			log.Printf("Due to IP Relase mode not polling ssm command: %v", *instance.InstanceId)
		}
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

func (c *commandStruct) runCommand(instanceId string) {
	document := ""
	commandToRun := c.command
	if c.os == "Linux" {
		document = "AWS-RunShellScript"
		if c.encode {
			commandToRun = "cd /tmp ; base64 -d <<< " + c.command + " | sh"
		}
	} else {
		document = "AWS-RunPowerShellScript"
		if c.encode {
			commandToRun = "cd C:/ ; powershell.exe -EncodedCommand " + c.command
		}

	}

	output, err := c.ssmClient.SendCommand(&ssm.SendCommandInput{
		DocumentName: aws.String(document),
		InstanceIds:  []*string{aws.String(instanceId)},
		Parameters: map[string][]*string{
			"commands": {aws.String(commandToRun)},
		},
		TimeoutSeconds: aws.Int64(600),
	})
	if err != nil {
		log.Println(err)
	}

	c.commandID = *output.Command.CommandId
	log.Printf("Got CommandID: %v", c.commandID)
}

func (c *commandStruct) pollCommand(instanceId string) {
	wait := true
	for wait {
		result, err := c.ssmClient.GetCommandInvocation(&ssm.GetCommandInvocationInput{
			CommandId:  aws.String(c.commandID),
			InstanceId: aws.String(instanceId),
		})
		if err != nil {
			continue
		} else {
			time.Sleep(time.Second * 2)
			if *result.Status == "Success" {
				log.Printf("Command Status: %v, Output: %v", *result.Status, *result.StandardOutputContent)
				wait = false
			} else if *result.Status == "Failed" {
				log.Printf("Command Status: %v, Output: %v ", *result.Status, *result.StandardErrorContent)
				wait = false

			}
		}

	}
}
