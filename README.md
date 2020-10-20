# Chaosity

A general chaos engineering platform to perform different chaos testing
scenarios on your environment. Currently the testing only consists of basic
scenarios for AWS but will support Kubernetes testing as well as others
in the future.

### Usage
#### Update
* Only works currently for AWS S3 Buckets, could also be reworked for github releases
```
Usage:
  chaosity update [flags]

Flags:
      --bucket string    Bucket to download update from(required)
  -h, --help             help for update
      --profile string   Profile used to download from S3 - Any standard profile will work (required)
      --region string    Region specified defaults to us-east-1 (default "us-east-1")  chaosity update [flags]
```
#### General
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
#### Reboot
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
#### Network Stop

```
chaosity aws instances --instances  i-077be821f09369884 --vpc-id vpc-05ef195e1cd6dd497  --region us-east-1 --profile default --duration 30 -n
2020/10/20 01:32:03 Getting SSM Information about instanceID: i-077be821f09369884
2020/10/20 01:32:03 Got OS: Linux for instance: i-077be821f09369884
2020/10/20 01:32:03 Got CommandID: bc42212d-d239-4d6e-923e-f7527721eea8
2020/10/20 01:32:03 Due to IP Relase mode not polling ssm command: i-077be821f09369884
2020/10/20 01:32:03 Chaos! Waiting for  30  seconds...
2020/10/20 01:32:33 Force shutting down instance  i-077be821f09369884
2020/10/20 01:32:33 Waiting 120 seconds for instances to stop.
2020/10/20 01:34:33 Starting Instances back up.
2020/10/20 01:34:33 Starting instance  i-077be821f09369884
