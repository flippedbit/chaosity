package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetInstancesBySubnet gathers all Instance objects from a list of subnets provided
// it will separate out multiple subnets with a comma (,) delimiter
// Returns a list of ec2.Instances and an error
func GetInstancesBySubnet(svc *ec2.EC2, s *string) ([]*ec2.Instance, error) {
	var instancesList []*ec2.Instance
	subnets := strings.Split(*s, ",")
	for _, subnet := range subnets {
		input := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("subnet-id"),
					Values: []*string{
						aws.String(subnet),
					},
				},
			},
		}
		instancesOutput, err := svc.DescribeInstances(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					return instancesList, fmt.Errorf(aerr.Error())
				}
			} else {
				return instancesList, fmt.Errorf(err.Error())
			}
		}
		for _, reservation := range instancesOutput.Reservations {
			for _, instance := range reservation.Instances {
				instancesList = append(instancesList, instance)
			}
		}
	}
	return instancesList, nil
}

// ApplyChaosSecurityGroupToInstances applies given SecurityGroup ID (sg) to all instances provided in the ec2.Instance list (instances)
// Returns an error
func ApplyChaosSecurityGroupToInstances(svc *ec2.EC2, instances []*ec2.Instance, sg string) error {
	for _, i := range instances {
		fmt.Println("Applying SecurityGroups ", sg, " to instance ", *i.InstanceId)
		input := &ec2.ModifyInstanceAttributeInput{
			InstanceId: aws.String(*i.InstanceId),
			Groups: []*string{
				aws.String(sg),
			},
		}
		if _, err := svc.ModifyInstanceAttribute(input); err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					return aerr
				}
			} else {
				return err
			}
		}
	}
	return nil
}

// RevertChaosSecurityGroupOnInstances applies the original SecurityGroups to each ec2.Instance within the list (instances)
// Returns error
func RevertChaosSecurityGroupOnInstances(svc *ec2.EC2, instances []*ec2.Instance) error {
	for _, i := range instances {
		var originalSecurityGroups []*string
		for _, s := range i.SecurityGroups {
			fmt.Println("Applying SecurityGroups ", *s.GroupId, " to instance ", *i.InstanceId)
			originalSecurityGroups = append(originalSecurityGroups, s.GroupId)
		}
		input := &ec2.ModifyInstanceAttributeInput{
			InstanceId: aws.String(*i.InstanceId),
			Groups:     originalSecurityGroups,
		}
		if _, err := svc.ModifyInstanceAttribute(input); err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					return aerr
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func RebootInstances(svc *ec2.EC2, instances []*ec2.Instance) error {
	var iList []*string
	for _, i := range instances {
		fmt.Println("Rebooting instance ", *i.InstanceId)
		iList = append(iList, i.InstanceId)
	}

	input := &ec2.RebootInstancesInput{
		InstanceIds: iList,
	}
	_, err := svc.RebootInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr)
				return aerr
			}
		} else {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
