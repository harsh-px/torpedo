package awsops

import (
	"fmt"
	"os"
	"sync"
	"time"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/portworx/sched-ops/task"
	"github.com/sirupsen/logrus"
)

// Ops are the various AWS operations supported by the package
type Ops interface {
	EC2Ops
}

// EC2Ops are AWS EC2 operations
type EC2Ops interface {
	// RebootNode reboots given ec2 instance
	RebootNode(instanceID string) error
	// ShutdownNode shuts down given ec2 instance
	ShutdownNode(instanceID string) error
	// GetAllInstances fetches all ec2 instances
	GetAllInstances() ([]*ec2.Instance, error)
	// SearchInstanceByAddresses searches ec2 instances matching given IP addresses
	SearchInstanceByAddresses(addresses []string) (*ec2.Instance, error)
	// RunCommand runs given command on the given instance and returns the output
	RunCommand(command string, instanceID string) (string, error)
}

var (
	instance Ops
	once     sync.Once
)

// Instance gives access to the aws ops singleton
func Instance() Ops {
	once.Do(func() {
		instance = &awsOps{}
	})

	return instance
}

type awsOps struct {
	svc    *ec2.EC2
	svcSsm *ssm.SSM
	region string
}

func (a *awsOps) initService() error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	creds := credentials.NewEnvCredentials()
	a.region = os.Getenv("AWS_REGION")
	if a.region == "" {
		return fmt.Errorf("env AWS_REGION not found")
	}
	config := &aws.Config{Region: aws.String(a.region)}
	config.WithCredentials(creds)
	a.svc = ec2.New(sess, config)
	a.svcSsm = ssm.New(sess, aws.NewConfig().WithRegion(a.region))

	return nil
}

func (a *awsOps) RebootNode(instanceID string) error {
	if err := a.initService(); err != nil {
		return err
	}

	rebootInput := &ec2.RebootInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	if _, err := a.svc.RebootInstances(rebootInput); err != nil {
		return err
	}

	return nil
}

func (a *awsOps) ShutdownNode(instanceID string) error {
	return nil
}

func (a *awsOps) GetAllInstances() ([]*ec2.Instance, error) {
	if err := a.initService(); err != nil {
		return nil, err
	}

	params := &ec2.DescribeInstancesInput{}
	resp, err := a.svc.DescribeInstances(params)
	if err != nil {
		return nil, fmt.Errorf("there was an error listing instances in %s. Error: %v", a.region, err)
	}

	instances := []*ec2.Instance{}
	for _, resv := range resp.Reservations {
		for _, ins := range resv.Instances {
			instances = append(instances, ins)
		}
	}
	return instances, nil
}

func (a *awsOps) SearchInstanceByAddresses(addresses []string) (*ec2.Instance, error) {
	instances, err := a.GetAllInstances()
	if err != nil {
		return nil, err
	}

	for _, i := range instances {
		for _, addr := range addresses {
			if aws.StringValue(i.PrivateIpAddress) == addr {
				return i, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to get for given addresses: %v", addresses)
}

func (a *awsOps) RunCommand(command string, instanceID string) (string, error) {
	if err := a.initService(); err != nil {
		return "", err
	}

	param := make(map[string][]*string)
	param["commands"] = []*string{
		aws.String(command),
	}
	s3BucketName := "awsops"
	s3KeyPrefix := fmt.Sprintf("%s", instanceID)

	sendCommandInput := &ssm.SendCommandInput{
		Comment:            aws.String(command),
		DocumentName:       aws.String("AWS-RunShellScript"),
		Parameters:         param,
		OutputS3BucketName: &s3BucketName,
		OutputS3KeyPrefix:  &s3KeyPrefix,
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	sendCommandOutput, err := a.svcSsm.SendCommand(sendCommandInput)
	if err != nil {
		return "", fmt.Errorf("failed to send command to instance %s: %v", instanceID, err)
	}

	if sendCommandOutput.Command == nil || sendCommandOutput.Command.CommandId == nil {
		return "", fmt.Errorf("no command returned after sending command to %s", instanceID)
	}

	t := func() (interface{}, error) {
		var status string

		listCmdsInput := &ssm.ListCommandInvocationsInput{
			CommandId: sendCommandOutput.Command.CommandId,
		}
		listCmdInvsOutput, _ := a.svcSsm.ListCommandInvocations(listCmdsInput)
		for _, cmd := range listCmdInvsOutput.CommandInvocations {
			status = strings.TrimSpace(*cmd.StatusDetails)
			if status == "Success" {
				logrus.Infof("[debug] cmd: %v", cmd)

				// TODO : Read command output from S3
				return "", nil
			}
		}

		return "", fmt.Errorf("Failed to connect. Command status is %s", status)
	}

	output, err := task.DoRetryWithTimeout(t, 1*time.Minute, 10*time.Second)
	if err != nil {
		return "", err
	}

	return output.(string), nil
}
