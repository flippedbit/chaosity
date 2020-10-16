package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/flippedbit/chaosity/pkg/aws/instancefilter"
	"github.com/flippedbit/chaosity/pkg/aws/options"
)

func GetInstances(svc *ec2.EC2, o options.AwsOptions) ([]*ec2.Instance, error) {
	var instancesList []*ec2.Instance
	f := &instancefilter.Filter{}
	f = f.ByRunning()

	if s := o.GetSubnets(); len(s) > 0 {
		f = f.BySubnet(s)
	}
	if a := o.GetAvailabilityZones(); len(a) > 0 {
		f = f.ByAvailabilityZone(a)
	}
	if i := o.GetInstances(); len(i) > 0 {
		f = f.ByInstance(i)
	}
	if t := o.GetTags(); len(t) > 0 {
		f = f.ByTag(t)
	}

	input := &ec2.DescribeInstancesInput{
		Filters: f.Build(),
	}
	// AWS describe on all instances with AWS error handling
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
	// cycle through instances and gather instances in order to return a list of ec2.Instance
	for _, reservation := range instancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			instancesList = append(instancesList, instance)
		}
	}

	return instancesList, nil
}

// ApplyChaosSecurityGroupToInstances applies given SecurityGroup ID (sg) to all instances provided in the ec2.Instance list (instances).
// Returns an error
func ApplyChaosSecurityGroupToInstances(svc *ec2.EC2, instances []*ec2.Instance, sg string) error {
	// cycle through given ec2.Instance list
	for _, i := range instances {
		log.Println("Applying SecurityGroups ", sg, " to instance ", *i.InstanceId)
		// modify every instance with the new segurity group ID for the chaos deny sg that was created and passed
		input := &ec2.ModifyInstanceAttributeInput{
			InstanceId: aws.String(*i.InstanceId),
			Groups: []*string{
				aws.String(sg),
			},
		}
		// apply instance changes with AWS error handling
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

// RevertChaosSecurityGroupOnInstances applies the original SecurityGroups to each ec2.Instance within the list (instances).
// Returns error
func RevertChaosSecurityGroupOnInstances(svc *ec2.EC2, instances []*ec2.Instance) error {
	// cycle through given ec2.Instance list
	for _, i := range instances {
		var originalSecurityGroups []*string
		// cycle through grabbing all original security groups for each instance
		for _, s := range i.SecurityGroups {
			log.Println("Applying SecurityGroups ", *s.GroupId, " to instance ", *i.InstanceId)
			originalSecurityGroups = append(originalSecurityGroups, s.GroupId)
		}
		// create the new instance attribute object to be applied combining the instanceID and securitygroupID's
		input := &ec2.ModifyInstanceAttributeInput{
			InstanceId: aws.String(*i.InstanceId),
			Groups:     originalSecurityGroups,
		}
		// apply the original securitygroups to each instance with AWS error handling
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

// RebootInstances reboots all instances given in an ec2.Instance list (instances).
// Returns an error
func RebootInstances(svc *ec2.EC2, instances []*ec2.Instance) error {
	// make sure we get instances otherwise RebootInstances() complains
	if len(instances) == 0 {
		log.Println("No instances passed to reboot.")
		return nil
	}
	var iList []*string
	// cycle through given ec2.Instance list gathering their InstanceID into a list of pointers []*string
	for _, i := range instances {
		log.Println("Rebooting instance ", *i.InstanceId)
		iList = append(iList, i.InstanceId)
	}

	// create the AWS structure of the instances to be rebooted
	input := &ec2.RebootInstancesInput{
		InstanceIds: iList,
	}
	// reboot all the instanceID's and check for AWS errors
	_, err := svc.RebootInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr)
				return aerr
			}
		} else {
			log.Println(err)
			return err
		}
	}
	return nil
}

// ForceShutdownInstances will issue a StopInstances with the force flag call to the AWS API.
// Accepts a pointer to an EC2 service and a list of EC2 instances.
// Returns an error
func ForceShutdownInstances(svc *ec2.EC2, instances []*ec2.Instance) error {
	// make sure we get instances otherwise StopInstances() complains
	if len(instances) == 0 {
		log.Println("No instances passed to shutdown.")
		return nil
	}
	var iList []*string
	t := true
	// cycle through given ec2.Instance list gathering their InstanceID into a list of pointers []*string
	for _, i := range instances {
		log.Println("Force shutting down instance ", *i.InstanceId)
		iList = append(iList, i.InstanceId)
	}

	input := &ec2.StopInstancesInput{
		Force:       &t,
		InstanceIds: iList,
	}

	_, err := svc.StopInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
				return aerr
			}
		} else {
			log.Println(err.Error())
			return err
		}
	}

	return nil
}

// StartInstances will attempt to start the AWS instances passed in the variable list "instances".
func StartInstances(svc *ec2.EC2, instances []*ec2.Instance) error {
	// make sure we get instances otherwise StartInstances() complains
	if len(instances) == 0 {
		log.Println("No instances passed to start.")
		return nil
	}
	var iList []*string
	// cycle through given ec2.Instance list gathering their InstanceID into a list of pointers []*string
	for _, i := range instances {
		log.Println("Starting instance ", *i.InstanceId)
		iList = append(iList, i.InstanceId)
	}

	input := &ec2.StartInstancesInput{
		InstanceIds: iList,
	}
	_, err := svc.StartInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
				return aerr
			}
		} else {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}
