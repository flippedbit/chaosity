# Chaosity

A general chaos engineering platform to perform different chaos testing
scenarios on your environment. Currently the testing only consists of basic
scenarios for AWS but will support Kubernetes testing as well as others
in the future.

### Usage
```
Usage:
  chaosity update [flags]

Flags:
      --bucket string    Bucket to download update from(required)
  -h, --help             help for update
      --profile string   Profile used to download from S3 - Any standard profile will work (required)
      --region string    Region specified defaults to us-east-1 (default "us-east-1")  chaosity update [flags]
```
```
Usage:
  chaosity aws instances [flags]

Flags:
  -d, --deny               Apply deny security group to instances.
  -h, --help               help for instances
      --instances string   Individual AWS Instance IDs to perform chaos on, comma separated.
  -r, --reboot             Reboot selected instances from subnets or availability-zone.
  -s, --shutdown           Force stop selected instances from subnets or availability-zone.
  -n, --stopnetwork        Stops network interface on target machines leaving applications in tact

Global Flags:
      --availability-zone string   AWS Availibility-Zone to perform chaos on.
      --duration int               How long to perform chaos testing for in seconds (default 300)
      --profile string             AWS credentials profile to use in order to connect (required)
      --region string              AWS region to perform chaos in (required)
      --subnets string             AWS Subnet IDs to perform chaos on (comma separated)
      --tags string                Tags applied to resource to match on. Format: key1=value1,key2=value2
      --vpc-id string              AWS VPC to perform chaos in (required)
```
Example:
```
chaosity aws instances --profile default --region us-east-1 --subnets subnet-75502f38,subnet-0dc0f76bf7ca00009 --duration 10 --deny --reboot
Created SecurityGroup  sg-02035af5bc9fc6afb
Applying SecurityGroups  sg-02035af5bc9fc6afb  to instance  i-0cc126878ba6b1610
Applying SecurityGroups  sg-02035af5bc9fc6afb  to instance  i-0f8dd4db1319d540e
Rebooting instance  i-0cc126878ba6b1610
Rebooting instance  i-0f8dd4db1319d540e
Chaos! Waiting for  10  seconds...
Applying SecurityGroups  sg-0470f76d5d6316aec  to instance  i-0cc126878ba6b1610
Applying SecurityGroups  sg-0dfc38f6b07adbace  to instance  i-0f8dd4db1319d540e
Deleting SecurityGroup  
```