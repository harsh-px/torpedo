package aws

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/portworx/sched-ops/task"
	"github.com/portworx/torpedo/drivers/node"
	"github.com/portworx/torpedo/pkg/awsops"
)

const (
	// DriverName is the name of the aws driver
	DriverName = "aws"
)

type aws struct {
	node.Driver
	region    string
	svc       *ec2.EC2
	svcSsm    *ssm.SSM
	instances []*ec2.Instance
}

func (a *aws) String() string {
	return DriverName
}

func (a *aws) Init() error {
	var err error
	a.instances, err = awsops.Instance().GetAllInstances()
	if err != nil {
		return err
	}

	nodes := node.GetWorkerNodes()
	for _, n := range nodes {
		if err := a.TestConnection(n, node.ConnectionOpts{
			Timeout:         1 * time.Minute,
			TimeBeforeRetry: 10 * time.Second,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (a *aws) TestConnection(n node.Node, options node.ConnectionOpts) error {
	inst, err := awsops.Instance().SearchInstanceByAddresses(n.Addresses)
	if err != nil {
		return &node.ErrFailedToTestConnection{
			Node:  n,
			Cause: fmt.Sprintf("failed to get instance ID to: %v", err),
		}
	}

	t := func() (interface{}, error) {
		_, err = awsops.Instance().RunCommand("uptime", *inst.InstanceId)
		return err
	}

	if _, err := task.DoRetryWithTimeout(t, options.Timeout, options.TimeBeforeRetry); err != nil {
		return &node.ErrFailedToTestConnection{
			Node:  n,
			Cause: fmt.Sprintf("failed to run command due to: %v", err),
		}
	}

	return err
}

func (a *aws) RebootNode(n node.Node, options node.RebootNodeOpts) error {
	inst, err := awsops.Instance().SearchInstanceByAddresses(n.Addresses)
	if err != nil {
		return &node.ErrFailedToRebootNode{
			Node:  n,
			Cause: err.Error(),
		}
	}

	err = awsops.Instance().RebootNode(*inst.InstanceId)
	if err != nil {
		return &node.ErrFailedToRebootNode{
			Node:  n,
			Cause: err.Error(),
		}
	}

	return err
}

func (a *aws) ShutdownNode(n node.Node, options node.ShutdownNodeOpts) error {
	return nil
}

func (a *aws) FindFiles(path string, n node.Node, options node.FindOpts) (string, error) {
	inst, err := awsops.Instance().SearchInstanceByAddresses(n.Addresses)
	if err != nil {
		return "", &node.ErrFailedToFindFileOnNode{
			Node:  n,
			Cause: err.Error(),
		}
	}

	if inst == nil {
		return "", &node.ErrFailedToFindFileOnNode{
			Node:  n,
			Cause: "failed to find AWS instances for node",
		}
	}

	findCmd := "sudo find " + path
	if options.Name != "" {
		findCmd += " -name " + options.Name
	}
	if options.MinDepth > 0 {
		findCmd += " -mindepth " + strconv.Itoa(options.MinDepth)
	}
	if options.MaxDepth > 0 {
		findCmd += " -maxdepth " + strconv.Itoa(options.MaxDepth)
	}

	t := func() (interface{}, error) {
		return awsops.Instance().RunCommand(findCmd, *inst.InstanceId)
	}

	if output, err := task.DoRetryWithTimeout(t, options.Timeout, options.TimeBeforeRetry); err != nil {
		return "", &node.ErrFailedToFindFileOnNode{
			Node:  n,
			Cause: fmt.Sprintf("failed to run command due to: %v", err),
		}
	}

	return output.(string), nil
}

func init() {
	a := &aws{
		Driver: node.NotSupportedDriver,
	}
	node.Register(DriverName, a)
}
