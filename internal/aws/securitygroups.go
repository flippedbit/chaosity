package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GenerateDenySecurityGroup used to create an emtpy SecurityGroup within the AWS VPC given.
// Returns SecurityGroup-ID as a string and an error
func GenerateDenySecurityGroup(svc *ec2.EC2, vpc *string) (string, error) {
	// generate default template for a security group, by default will deny all ingress and permit all egress
	groupInput := &ec2.CreateSecurityGroupInput{
		Description: aws.String("Deny security group used for chaos testing"),
		GroupName:   aws.String("chaos-deny-security-group"),
		VpcId:       vpc,
	}
	// create the security group and check for AWS errors
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

	// describe the newly created security group so we can make sure to remove the exact egress IpPermissionsEgress late on
	/**
	 * @todo
	 * @body create a generic SecurityGroup.IpPermissionsEgress structure so we do not need to make an additoinal api call
	 */
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

	// build struct to remove default egress rule from generated security group
	revokeSecurityGroupInput := &ec2.RevokeSecurityGroupEgressInput{
		GroupId:       aws.String(sgID),
		IpPermissions: sgIPPermissionsEgress,
	}
	// remove egress rule from security group and check for AWS errors
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

// DeleteDenySecurityGroup removes the given SecurityGroup ID from AWS.
// Returns an error
func DeleteDenySecurityGroup(svc *ec2.EC2, sg string) error {
	// create struct to for deleting the originally created security group
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sg),
	}
	// delete the originally created security group and check for AWS errors
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
