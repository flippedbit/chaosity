package instancefilter

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Filter is the default struct to hold all of the ec2.Filters that would be
// constructed based on the flags passed.
type Filter struct {
	f []*ec2.Filter
}

// ByRunnin will add a filter to only gather instances that are running.
func (f *Filter) ByRunning() *Filter {
	f.f = append(f.f, &ec2.Filter{
		Name: aws.String("instance-state-name"),
		Values: []*string{
			aws.String("running"),
		},
	})
	return f
}

// BySubnet will add a filter to the struct for searching within a given
// subnet. Multiple subnets would be given comma separated.
func (f *Filter) BySubnet(s []string) *Filter {
	f.f = append(f.f, &ec2.Filter{
		Name:   aws.String("subnet-id"),
		Values: aws.StringSlice(s),
	})
	return f
}

// ByAvailabilityZone will add a filter to the struct for searching
// within a given availability-zone. Multiple availability-zones would
// be comma separated.
func (f *Filter) ByAvailabilityZone(s []string) *Filter {
	f.f = append(f.f, &ec2.Filter{
		Name:   aws.String("availability-zone"),
		Values: aws.StringSlice(s),
	})
	return f
}

// ByInstance will add a filter to the struct for searching
// by instance-id.Multiple instance-id's would be comma separated.
func (f *Filter) ByInstance(s []string) *Filter {
	f.f = append(f.f, &ec2.Filter{
		Name:   aws.String("instance-id"),
		Values: aws.StringSlice(s),
	})
	return f
}

// ByTag will add a filter to the struct for searching by tag.
// Multiple tags are passed comma separated, each
// tag is represented key=value.
// example: team=devops
func (f *Filter) ByTag(s string) *Filter {
	ss := strings.Split(s, ",")
	for _, t := range ss {
		tag := strings.Split(t, "=")
		tKey := fmt.Sprintf("tag:%s", tag[0])
		f.f = append(f.f, &ec2.Filter{
			Name: aws.String(tKey),
			Values: []*string{
				aws.String(tag[1]),
			},
		})
	}
	return f
}

// Build will construct the AWS ec2 Filter for gathering instances
// based on given flag.
// Returns a list of ec2 filters
func (f *Filter) Build() []*ec2.Filter {
	return f.f
}
