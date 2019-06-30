package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/glassechidna/ec2connect/pkg/ec2connect"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func init() {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "SSH ProxyCommand implementation",
		Run: func(cmd *cobra.Command, args []string) {
			instanceId, _ := cmd.PersistentFlags().GetString("instance-id")
			region, _ := cmd.PersistentFlags().GetString("region")
			user, _ := cmd.PersistentFlags().GetString("user")
			sshKeyPath, _ := cmd.PersistentFlags().GetString("ssh-key")

			info, err := authorize(instanceId, region, user, sshKeyPath)
			if err != nil {
				if awsErr, ok := errors.Cause(err).(awserr.Error); ok {
					if awsErr.Code() == credentials.ErrNoValidProvidersFoundInChain.Code() {
						fmt.Fprintln(os.Stderr, `
No AWS credentials found.

* You can specify one of the profiles from ~/.aws/config by setting the 
  AWS_PROFILE environment variable.

* You can set AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and optionally
  AWS_SESSION_TOKEN.`)
					} else if strings.HasPrefix(awsErr.Code(), "InvalidInstanceID.") {
						fmt.Fprintf(os.Stderr, `
No instance found with ID %s. Try specifying an explicit region using the 
AWS_REGION environment variable.

`, instanceId)
					}
					return
				} else {
					panic(err)
				}
			}

			err = connect(info.Address + ":22")
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.PersistentFlags().String("instance-id", "", "")
	cmd.PersistentFlags().String("region", "", "")
	cmd.PersistentFlags().String("user", "ec2-user", "")
	cmd.PersistentFlags().String("ssh-key", "~/.ssh/ec2connect/id_rsa.pub", "")

	RootCmd.AddCommand(cmd)
}

func authorize(instanceId, region, user, sshKeyPath string) (*ec2connect.ConnectionInfo, error) {
	path, err := homedir.Expand(sshKeyPath)
	if err != nil {
		return nil, err
	}

	sshKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		Config:                  *aws.NewConfig().WithRegion(region), //.WithLogLevel(aws.LogDebugWithHTTPBody),
	})
	if err != nil {
		return nil, err
	}

	key, err := ec2connect.NormalizeKey(string(sshKey))
	if err != nil {
		return nil, err
	}

	auth := &ec2connect.Authorizer{Ec2Api: ec2.New(sess), ConnectApi: ec2instanceconnect.New(sess)}
	return auth.Authorize(context.Background(), instanceId, user, key)
}

func connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "establishing connection to %s", addr)
	}

	go func() {
		io.Copy(os.Stdout, conn)
		conn.Close()
	}()

	io.Copy(conn, os.Stdin)
	conn.Close()

	return nil
}
