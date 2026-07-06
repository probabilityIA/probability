package awsec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/secamc93/probability/back/central/services/modules/woostore/internal/domain"
)

type Client struct {
	ec2        *ec2.Client
	instanceID string
	storeURL   string
}

func New(region, key, secret, instanceID, storeURL string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(key, secret, "")),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		ec2:        ec2.NewFromConfig(cfg),
		instanceID: instanceID,
		storeURL:   storeURL,
	}, nil
}

func (c *Client) Start(ctx context.Context) (*domain.PowerState, error) {
	if _, err := c.ec2.StartInstances(ctx, &ec2.StartInstancesInput{InstanceIds: []string{c.instanceID}}); err != nil {
		return nil, err
	}
	return c.Status(ctx)
}

func (c *Client) Stop(ctx context.Context) (*domain.PowerState, error) {
	if _, err := c.ec2.StopInstances(ctx, &ec2.StopInstancesInput{InstanceIds: []string{c.instanceID}}); err != nil {
		return nil, err
	}
	return c.Status(ctx)
}

func (c *Client) Status(ctx context.Context) (*domain.PowerState, error) {
	out, err := c.ec2.DescribeInstances(ctx, &ec2.DescribeInstancesInput{InstanceIds: []string{c.instanceID}})
	if err != nil {
		return nil, err
	}
	st := &domain.PowerState{InstanceID: c.instanceID, StoreURL: c.storeURL}
	if len(out.Reservations) > 0 && len(out.Reservations[0].Instances) > 0 {
		inst := out.Reservations[0].Instances[0]
		if inst.State != nil {
			st.State = string(inst.State.Name)
		}
		if inst.PublicIpAddress != nil {
			st.PublicIP = *inst.PublicIpAddress
		}
	}
	return st, nil
}
