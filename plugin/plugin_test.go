// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"
	"github.com/google/go-cmp/cmp"

	"gopkg.in/h2non/gock.v1"
)

var noContext = context.Background()

func TestPlugin(t *testing.T) {
	defer gock.Off()

	gock.New("https://ec2.us-east-1.amazonaws.com").
		Post("/").
		MatchHeader("Content-Type", "application/x-amz-json-1.1").
		MatchHeader("X-Amz-Target", "AmazonSSM.GetParameter").
		Reply(200).
		File("testdata/parameter.json")

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	config := defaults.Config()
	config.HTTPClient = client
	config.Region = "us-east-1"
	config.EndpointResolver = aws.ResolveWithEndpoint(aws.Endpoint{
		SigningRegion: config.Region,
		URL:           "https://ec2.us-east-1.amazonaws.com",
	})

	const paramName = "/staging/webapp/DATABASE_URL"
	req := &secret.Request{
		Name: paramName,
	}

	ssmClient := ssm.New(config)
	plugin := New(ssmClient)

	got, err := plugin.Find(noContext, req)

	if err != nil {
		t.Error(err)
	}

	want := &drone.Secret{
		Data: "mysql://fakedburl",
		Name: paramName,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
		return
	}
}

func TestPlugin_InvalidParamName(t *testing.T) {
	req := &secret.Request{
		Name: "",
	}
	_, err := New(nil).Find(noContext, req)

	if err == nil {
		t.Error("Invalid parameter name error is expected")
	}

	if got, want := err.Error(), "invalid or missing secret name"; got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestPlugin_ParameterNotFound(t *testing.T) {
	defer gock.Off()

	gock.New("https://ec2.us-east-1.amazonaws.com").
		Post("/").
		MatchHeader("Content-Type", "application/x-amz-json-1.1").
		MatchHeader("X-Amz-Target", "AmazonSSM.GetParameter").
		Reply(400).
		AddHeader("x-amzn-requestid", "fakerequestid").
		BodyString(`{"__type":"ParameterNotFound"}`)

	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	config := defaults.Config()
	config.HTTPClient = client
	config.Region = "us-east-1"
	config.EndpointResolver = aws.ResolveWithEndpoint(aws.Endpoint{
		SigningRegion: config.Region,
		URL:           "https://ec2.us-east-1.amazonaws.com",
	})

	const paramName = "test"
	req := &secret.Request{
		Name: paramName,
	}

	ssmClient := ssm.New(config)
	plugin := New(ssmClient)

	_, err := plugin.Find(noContext, req)

	if err == nil {
		t.Error("Invalid parameter name error is expected")
	}

	const notFoundMsg = `couldn't retrieve parameter from SSM: ParameterNotFound: 
	status code: 400, request id: fakerequestid`

	if got, want := err.Error(), notFoundMsg; got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}
