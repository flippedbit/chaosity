# Chaosity

A general chaos engineering platform to perform different chaos testing
scenarios on your environment. Currently the testing only consists of basic
scenarios for AWS but will support Kubernetes testing as well as others
in the future.

### Usage
```
Usage:
  chaosity aws instances [flags]

Flags:
  -d, --deny     Apply deny security group to instances
  -h, --help     help for instances
  -r, --reboot   Reboot selected instances from subnets or availability-zone.

Global Flags:
  -a, --author string              author name for copyright attribution (default "Michael Straughan")
      --availability-zone string   AWS Availibility-Zone to perform chaos on.
      --config string              config file (default is $HOME/.chaosity.yaml)
      --duration int               How long to perform chaos testing for in seconds (default 300)
      --profile string             AWS credentials profile to use in order to connect (required)
      --region string              AWS region to perform chaos in (required)
      --subnets string             AWS Subnet IDs to perform chaos on (comma separated)
      --vpc-id string              AWS VPC to perform chaos in (required)
```
Example:
```
chaosity aws instances --profile default --region us-east-1 --subnets subnet-75502f38,subnet-0dc0f76bf7ca00009 --duration 10
Created SecurityGroup  sg-0d1324bc8103607c0
Applying SecurityGroups  sg-0d1324bc8103607c0  to instance  i-0cc126878ba6b1610
Applying SecurityGroups  sg-0d1324bc8103607c0  to instance  i-0f8dd4db1319d540e
Applying SecurityGroups  sg-0470f76d5d6316aec  to instance  i-0cc126878ba6b1610
Applying SecurityGroups  sg-0dfc38f6b07adbace  to instance  i-0f8dd4db1319d540e
Deleting SecurityGroup  sg-0d1324bc8103607c0
```