package networkacl

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateDenyNacl creates a blank AWS NetworkACL within the provided VPC
// in order to deny all inbound/outbound traffic to the subnet it is later
// applied to.
func CreateDenyNacl(svc *ec2.EC2, vpc string) (string, error) {
	input := &ec2.CreateNetworkAclInput{
		VpcId: aws.String(vpc),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("network-acl"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("app"),
						Value: aws.String("chaosity"),
					},
				},
			},
		},
	}
	result, err := svc.CreateNetworkAcl(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return "", aerr
			}
		} else {
			return "", err
		}
	}
	n := *result.NetworkAcl.NetworkAclId
	fmt.Println("Created NACL: ", n)
	return n, nil
}

// DeleteDenyNacl removes the AWS NetworkACL that is passed in.
func DeleteDenyNacl(svc *ec2.EC2, n string) error {
	input := &ec2.DeleteNetworkAclInput{
		NetworkAclId: aws.String(n),
	}
	_, err := svc.DeleteNetworkAcl(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return aerr
			}
		} else {
			return err
		}
	}
	fmt.Println("Removing NACL: ", n)
	return nil
}
