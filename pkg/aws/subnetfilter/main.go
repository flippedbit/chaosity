package subnetfilter

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SubnetFilter struct {
	f []*ec2.Filter
}

func (sf *SubnetFilter) BySubnet(s []string) *SubnetFilter {
	sf.f = append(sf.f, &ec2.Filter{
		Name:   aws.String("subnet-id"),
		Values: aws.StringSlice(s),
	})
	return sf
}

func (sf *SubnetFilter) ByAvailabilityZone(s []string) *SubnetFilter {
	sf.f = append(sf.f, &ec2.Filter{
		Name:   aws.String("availability-zone"),
		Values: aws.StringSlice(s),
	})
	return sf
}

func (sf *SubnetFilter) ByTag(s []string) *SubnetFilter {
	for _, t := range s {
		tag := strings.Split(t, "=")
		tKey := fmt.Sprintf("tag:%s", tag[0])
		sf.f = append(sf.f, &ec2.Filter{
			Name: aws.String(tKey),
			Values: []*string{
				aws.String(tag[1]),
			},
		})
	}
	return sf
}

func (sf *SubnetFilter) Build() []*ec2.Filter {
	return sf.f
}
