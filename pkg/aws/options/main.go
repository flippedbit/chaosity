package options

import "strings"

type AwsOptions struct {
	Region   string
	VpcID    string
	Profile  string
	Subnets  string
	Az       string
	Duration int
}

func (a *AwsOptions) GetSubnets() []string {
	return strings.Split(a.Subnets, ",")
}

func (a *AwsOptions) GetAvailabilityZones() []string {
	return strings.Split(a.Az, ",")
}
