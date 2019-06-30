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
	"github.com/glassechidna/ec2connect/pkg/sshconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
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
			port, _ := cmd.PersistentFlags().GetInt("port")
			err := connect(instanceId, region, user, port)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.PersistentFlags().String("instance-id", "", "")
	cmd.PersistentFlags().String("region", "", "")
	cmd.PersistentFlags().String("user", "ec2-user", "")
	cmd.PersistentFlags().Int("port", 22, "")

	RootCmd.AddCommand(cmd)
}

func connect(instanceId, region, user string, port int) error {
	conf := sshconfig.DefaultSsh.Get(instanceId, user, port)
	pubKeyBytes := conf.EffectivePublicKey()

	info, err := authorize(instanceId, region, user, string(pubKeyBytes))
	if err != nil {
		if awsErr, ok := errors.Cause(err).(awserr.Error); ok {
			if awsErr.Code() == credentials.ErrNoValidProvidersFoundInChain.Code() {
				_, _ = fmt.Fprintln(os.Stderr, `
No AWS credentials found.

* You can specify one of the profiles from ~/.aws/config by setting the 
  AWS_PROFILE environment variable.

* You can set AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and optionally
  AWS_SESSION_TOKEN.`)
			} else if strings.HasPrefix(awsErr.Code(), "InvalidInstanceID.") {
				_, _ = fmt.Fprintf(os.Stderr, `
No instance found with ID %s. Try specifying an explicit region using the 
AWS_REGION environment variable.

`, instanceId)
			} else {
				return err
			}
			os.Exit(1)
		} else {
			return err
		}
	}

	return tunnel(fmt.Sprintf("%s:%d", info.Address, port))
}

func authorize(instanceId, region, user, sshKey string) (*ec2connect.ConnectionInfo, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		Config:                  *aws.NewConfig().WithRegion(region), //.WithLogLevel(aws.LogDebugWithHTTPBody),
	})
	if err != nil {
		return nil, err
	}

	auth := &ec2connect.Authorizer{Ec2Api: ec2.New(sess), ConnectApi: ec2instanceconnect.New(sess)}
	return auth.Authorize(context.Background(), instanceId, user, sshKey)
}

func tunnel(addr string) error {
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
