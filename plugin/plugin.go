// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"
)

// New returns a new secret plugin.
// Takes client param that is a type which implements ssmiface.ClientAPI interface
func New(client ssmiface.ClientAPI) secret.Plugin {
	return &plugin{
		client: client,
	}
}

type plugin struct {
	client ssmiface.ClientAPI
}

func (p *plugin) Find(ctx context.Context, req *secret.Request) (*drone.Secret, error) {

	if req.Name == "" {
		return nil, errors.New("invalid or missing secret name")
	}

	response, err := p.client.GetParameterRequest(
		&ssm.GetParameterInput{
			Name:           aws.String(req.Name),
			WithDecryption: aws.Bool(true),
		},
	).Send(ctx)

	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve parameter from SSM: %s", err)
	}

	return &drone.Secret{
		Name: req.Name,
		Data: *response.Parameter.Value,
	}, nil
}
