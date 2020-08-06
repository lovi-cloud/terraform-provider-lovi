package lovi

import (
	"errors"
	"fmt"

	satelitpb "github.com/whywaita/satelit/api/satelit"
	"google.golang.org/grpc"
)

// Config is config of terraform-provider-lovi
type Config struct {
	APIEndpoint string

	SatelitClient satelitpb.SatelitClient
}

// LoadAndValidate performs to connect and init configuration
func (c *Config) LoadAndValidate() error {
	if c.APIEndpoint == "" {
		return errors.New("SATELIT_API_ENDPOINT must be set")
	}

	conn, err := grpc.Dial(c.APIEndpoint, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect satelit server")
	}
	client := satelitpb.NewSatelitClient(conn)

	c.SatelitClient = client

	return nil
}
