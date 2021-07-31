# terraform-provider-lvslb

![GitHub release (latest by date)](https://img.shields.io/github/v/release/jeremmfr/terraform-provider-lvslb)
[![Registry](https://img.shields.io/badge/registry-doc%40latest-lightgrey?logo=terraform)](https://registry.terraform.io/providers/jeremmfr/lvslb/latest/docs)
[![Go Status](https://github.com/jeremmfr/terraform-provider-lvslb/workflows/Go%20Tests/badge.svg)](https://github.com/jeremmfr/terraform-provider-lvslb/actions)
[![Lint Status](https://github.com/jeremmfr/terraform-provider-lvslb/workflows/GolangCI-Lint/badge.svg)](https://github.com/jeremmfr/terraform-provider-lvslb/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeremmfr/terraform-provider-lvslb)](https://goreportcard.com/report/github.com/jeremmfr/terraform-provider-lvslb)

Terraform's provider to generate keepalived virtual_server with [lvslb-api](https://github.com/jeremmfr/lvslb-api)

## Automatic install (Terraform 0.13 and later)

Add source information inside the Terraform configuration block for automatic provider installation:

```hcl
terraform {
  required_providers {
    lvslb = {
      source = "jeremmfr/lvslb"
    }
  }
}
```

## Documentation

[registry.terraform.io](https://registry.terraform.io/providers/jeremmfr/lvslb/latest/docs)

or in docs :

[terraform-provider-lvslb](docs/index.md)  

Resources:

* [lvslb_ipvs](docs/resources/ipvs.md)

## Compile

```shell
go build
```
