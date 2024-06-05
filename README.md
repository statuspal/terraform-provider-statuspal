# Terraform Provider StatusPal (Terraform Plugin Framework)

This a Terraform provider for interacting with [StatusPal API](https://www.statuspal.io/api-docs).

**Visit [statuspal.io](https://www.statuspal.io/).**

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install .` from the root directory. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to run the full suite of Acceptance tests, run `TF_ACC=1 go test -v -cover ./internal/provider/` from the root directory.

To test manually the resource or data source, run from the root directory:
- apply the terraform plan: `terraform -chdir=./examples/<resource_or_data_source_name>/ apply --auto-approve`
- destroy the resource: `terraform -chdir=./examples/<resource_or_data_source_name>/ destroy --auto-approve`

To generate or update documentation, run `go generate ./...` from the root directory.

> [!NOTE]
> **For more information visit [Implement a provider with the Terraform Plugin Framework](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider)**:
>
> - Skip the Docker part in [**Set up your development environment**](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#set-up-your-development-environment) section and you have to run [statuspal/statushq](https://github.com/statuspal/statushq) locally.
> - In the [**Prepare Terraform for local provider install**](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install) section, the `~/.terraformrc` file should look like this:
>   ```terraform
>   provider_installation {
>
>     dev_overrides {
>         "registry.terraform.io/statuspal/statuspal" = "<PATH>"
>     }
>
>     # For all other providers, install them directly from their origin provider
>     # registries as normal. If you omit this, Terraform will _only_ use
>     # the dev_overrides block, and so no other providers will be available.
>     direct {}
>   }
>   ```

> [!IMPORTANT]
> Run `golangci-lint run` before you push a commit and fix all the showed errors.

## Create a provider release

- Add changes into [CHANGELOG.md](https://github.com/statuspal/terraform-provider-statuspal/blob/main/CHANGELOG.md) file.
- Follow the instructions in the [Create a provider release](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-release-publish#create-a-provider-release) section.
