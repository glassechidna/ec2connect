package ec2connect

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAuthorizer_Authorize(t *testing.T) {
	t.Run("err describing", func(t *testing.T) {
		connApi := &mockConnect{}
		ec2Api := &mockEc2{}
		ec2Api.
			On("DescribeInstancesWithContext", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.AnythingOfType("[]request.Option")).
			Return(nil, errors.New("err"))

		auth := &Authorizer{Ec2Api: ec2Api, ConnectApi: connApi}
		_, err := auth.Authorize(context.Background(), "i-012abc", "", "")
		assert.Error(t, err)
	})

	t.Run("no instances", func(t *testing.T) {
		connApi := &mockConnect{}
		ec2Api := &mockEc2{}
		ec2Api.
			On("DescribeInstancesWithContext", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.AnythingOfType("[]request.Option")).
			Return(&ec2.DescribeInstancesOutput{}, nil)

		auth := &Authorizer{Ec2Api: ec2Api, ConnectApi: connApi}
		_, err := auth.Authorize(context.Background(), "i-012abc", "", "")
		assert.Error(t, err)
	})

	t.Run("err sending ssh key", func(t *testing.T) {
		ec2Api := &mockEc2{}
		ec2Api.
			On("DescribeInstancesWithContext", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.AnythingOfType("[]request.Option")).
			Return(&ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								Placement: &ec2.Placement{
									AvailabilityZone: aws.String("ap-southeast-2b"),
								},
							},
						},
					},
				},
			}, nil)

		connApi := &mockConnect{}
		connApi.
			On("SendSSHPublicKeyWithContext", mock.Anything, mock.AnythingOfType("*ec2instanceconnect.SendSSHPublicKeyInput"), mock.AnythingOfType("[]request.Option")).
			Return(nil, errors.New("err"))

		auth := &Authorizer{Ec2Api: ec2Api, ConnectApi: connApi}
		_, err := auth.Authorize(context.Background(), "i-012abc", "", "")
		assert.Error(t, err)
	})

	t.Run("unsuccessful send ssh key", func(t *testing.T) {
		ec2Api := &mockEc2{}
		ec2Api.
			On("DescribeInstancesWithContext", mock.Anything, mock.AnythingOfType("*ec2.DescribeInstancesInput"), mock.AnythingOfType("[]request.Option")).
			Return(&ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								Placement: &ec2.Placement{
									AvailabilityZone: aws.String("ap-southeast-2b"),
								},
							},
						},
					},
				},
			}, nil)

		connApi := &mockConnect{}
		connApi.
			On("SendSSHPublicKeyWithContext", mock.Anything, mock.AnythingOfType("*ec2instanceconnect.SendSSHPublicKeyInput"), mock.AnythingOfType("[]request.Option")).
			Return(&ec2instanceconnect.SendSSHPublicKeyOutput{
				Success: aws.Bool(false),
			}, nil)

		auth := &Authorizer{Ec2Api: ec2Api, ConnectApi: connApi}
		_, err := auth.Authorize(context.Background(), "i-012abc", "", "")
		assert.Error(t, err)
	})
}

type mockEc2 struct {
	mock.Mock
	ec2iface.EC2API
}

func (m *mockEc2) DescribeInstancesWithContext(ctx aws.Context, input *ec2.DescribeInstancesInput, opts ...request.Option) (*ec2.DescribeInstancesOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*ec2.DescribeInstancesOutput)
	return output, f.Error(1)
}

type mockConnect struct {
	mock.Mock
	ec2instanceconnectiface.EC2InstanceConnectAPI
}

func (m *mockConnect) SendSSHPublicKeyWithContext(ctx aws.Context, input *ec2instanceconnect.SendSSHPublicKeyInput, opts ...request.Option) (*ec2instanceconnect.SendSSHPublicKeyOutput, error) {
	f := m.Called(ctx, input, opts)
	output, _ := f.Get(0).(*ec2instanceconnect.SendSSHPublicKeyOutput)
	return output, f.Error(1)
}
