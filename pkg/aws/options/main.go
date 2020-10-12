package options

import "strings"

type AwsOptions struct {
	Region    string
	VpcID     string
	Profile   string
	Subnets   string
	Az        string
	Duration  int
	Instances string
}

func (a *AwsOptions) GetSubnets() []string {
	return strings.Split(a.Subnets, ",")
}

func (a *AwsOptions) GetAvailabilityZones() []string {
	return strings.Split(a.Az, ",")
}

func (a *AwsOptions) GetInstances() []string {
	return strings.Split(a.Instances, ",")
}
