package subnets

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/flippedbit/chaosity/pkg/aws/options"
	"github.com/flippedbit/chaosity/pkg/aws/subnetfilter"
)

func GetSubnets(svc *ec2.EC2, o options.AwsOptions) ([]*ec2.Subnet, error) {
	sf := &subnetfilter.SubnetFilter{}

	if s := o.GetSubnets(); len(s) > 0 {
		sf = sf.BySubnet(s)
	}
	if a := o.GetAvailabilityZones(); len(a) > 0 {
		sf = sf.ByAvailabilityZone(a)
	}
	if t := o.GetTags(); len(t) > 0 {
		sf = sf.ByTag(t)
	}

	input := &ec2.DescribeSubnetsInput{
		Filters: sf.Build(),
	}
	result, err := svc.DescribeSubnets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return []*ec2.Subnet{}, aerr
			}
		} else {
			return []*ec2.Subnet{}, err
		}
	}
	return result.Subnets, nil
}
