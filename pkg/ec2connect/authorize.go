package ec2connect

import (
	"context"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type Authorizer struct {
	Ec2Api     ec2iface.EC2API
	ConnectApi ec2instanceconnectiface.EC2InstanceConnectAPI
}

type ConnectionInfo struct {
	Address   string
	RequestId string
}

func (c *Authorizer) Authorize(ctx context.Context, instanceId, user, sshKey string) (*ConnectionInfo, error) {
	r, err := c.Ec2Api.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{InstanceIds: []*string{&instanceId}})
	if err != nil {
		return nil, errors.Wrap(err, "describing instance")
	}

	if len(r.Reservations) == 0 || len(r.Reservations[0].Instances) == 0 {
		return nil, errors.Errorf("no instance with id %s", instanceId)
	}

	instance := r.Reservations[0].Instances[0]
	az := instance.Placement.AvailabilityZone

	r2, err := c.ConnectApi.SendSSHPublicKeyWithContext(ctx, &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: az,
		InstanceId:       &instanceId,
		InstanceOSUser:   &user,
		SSHPublicKey:     &sshKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sending ssh key to instance")
	}

	if !*r2.Success {
		return nil, errors.Errorf("sending ssh key to instance")
	}

	ip := *instance.PrivateIpAddress
	if instance.PublicIpAddress != nil { // TODO: they might *want* the private ip
		ip = *instance.PublicIpAddress
	}

	return &ConnectionInfo{
		Address: ip,
		RequestId: *r2.RequestId,
	}, nil
}

func NormalizeKey(input string) (string, error) {
	s, err := ssh.ParsePrivateKey([]byte(input))
	if err != nil {
		return "", errors.Wrap(err, "parsing private key")
	}

	return string(ssh.MarshalAuthorizedKey(s.PublicKey())), nil
}