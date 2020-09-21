# Chaosity

A general chaos engineering platform to perform different chaos testing
scenarios on your environment. Currently the testing only consists of basic
scenarios for AWS but will support Kubernetes testing as well as others
in the future.

### Usage
```
./chaosity aws instances --profile <profile> --regions <region> --vpc-id <vpc> --subnets subnet-1,subnet-2
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
