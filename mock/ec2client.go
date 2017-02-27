package mock

import "github.com/aws/aws-sdk-go/service/ec2"

// The EC2Client struct holds the mock implementation of the EC2Client, to
// facilitate testing
type EC2Client struct {
	DescribeInstancesFn        func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	DescribeInstancesFnInvoked bool
}

// DescribeInstances is a mock implementation of ec2.DescribeInstances
func (m *EC2Client) DescribeInstances(params *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	m.DescribeInstancesFnInvoked = true
	if m.DescribeInstancesFn == nil {
		return m.defaultDescribeInstancesFn(params)
	}
	return m.DescribeInstancesFn(params)
}

func (m *EC2Client) defaultDescribeInstancesFn(params *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	instances := []*ec2.Instance{
		{},
		{},
	}
	reservation := &ec2.Reservation{
		Instances: instances,
	}
	reservations := []*ec2.Reservation{reservation}
	return &ec2.DescribeInstancesOutput{
		Reservations: reservations,
	}, nil
}
