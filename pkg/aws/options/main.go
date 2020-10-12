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
	if a.Subnets != "" {
		return strings.Split(a.Subnets, ",")
	}
	return []string{}
}

func (a *AwsOptions) GetAvailabilityZones() []string {
	if a.Az != "" {
		return strings.Split(a.Az, ",")
	}
	return []string{}
}

func (a *AwsOptions) GetInstances() []string {
	if a.Instances != "" {
		return strings.Split(a.Instances, ",")
	}
	return []string{}
}
