package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GenerateDenySecurityGroup used to create an emtpy SecurityGroup within the AWS VPC given
// Returns SecurityGroup-ID as a string and an error
func GenerateDenySecurityGroup(svc *ec2.EC2, vpc *string) (string, error) {
	groupInput := &ec2.CreateSecurityGroupInput{
		Description: aws.String("Deny security group used for chaos testing"),
		GroupName:   aws.String("chaos-deny-security-group"),
		VpcId:       vpc,
	}
	groupResult, err := svc.CreateSecurityGroup(groupInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return "", aerr
			}
		} else {
			return "", err
		}
		return "", err
	}
	sgID := *groupResult.GroupId

	describeGroupInput := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(sgID),
		},
	}
	describeGroupResult, err := svc.DescribeSecurityGroups(describeGroupInput)
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
	sgIPPermissionsEgress := describeGroupResult.SecurityGroups[0].IpPermissionsEgress

	revokeSecurityGroupInput := &ec2.RevokeSecurityGroupEgressInput{
		GroupId:       aws.String(sgID),
		IpPermissions: sgIPPermissionsEgress,
	}
	if _, err = svc.RevokeSecurityGroupEgress(revokeSecurityGroupInput); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return "", aerr
			}
		} else {
			return "", err
		}
	}

	return sgID, nil
}

// DeleteDenySecurityGroup removes the given SecurityGroup ID from AWS
// Returns an error
func DeleteDenySecurityGroup(svc *ec2.EC2, sg string) error {
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sg),
	}
	_, err := svc.DeleteSecurityGroup(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return fmt.Errorf(aerr.Error())
			}
		} else {
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}
