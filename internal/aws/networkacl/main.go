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

func GetNetworkAclAssociation(svc *ec2.EC2, s string) ([]ec2.NetworkAclAssociation, error) {
	var nacls []ec2.NetworkAclAssociation
	input := &ec2.DescribeNetworkAclsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("association.subnet-id"),
				Values: []*string{
					aws.String(s),
				},
			},
		},
	}
	result, err := svc.DescribeNetworkAcls(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return nacls, aerr
			}
		} else {
			return nacls, err
		}
	}
	for _, n := range result.NetworkAcls {
		for _, a := range n.Associations {
			if *a.SubnetId == s {
				nacls = append(nacls, *a)
			}
		}
	}
	//fmt.Println(nacls)
	return nacls, nil
}

func ReplaceAssociation(svc *ec2.EC2, a string, n string) (string, error) {
	input := &ec2.ReplaceNetworkAclAssociationInput{
		AssociationId: aws.String(a),
		NetworkAclId:  aws.String(n),
	}
	result, err := svc.ReplaceNetworkAclAssociation(input)
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
	new := *result.NewAssociationId
	fmt.Println("Replaced NetworkACL on association ", a, " with ", n, ". New associationID ", new)
	return new, nil
}
