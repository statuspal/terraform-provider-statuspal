package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// providerConfig is a shared configuration to combine with the actual test configuration.
func providerConfig(testUrl *string) *string {
	providerConfig := `
		provider "statuspal" {
			test_url = "` + *testUrl + `"
		}
	`
	return &providerConfig
}

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"statuspal": providerserver.NewProtocol6WithError(New("test")()),
	}
)
