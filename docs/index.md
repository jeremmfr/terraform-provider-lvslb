# lvslb Provider

Terraform's provider for generate keepalived virtual_server with [lvslb-api](https://github.com/jeremmfr/lvslb-api)

## Example Usage

```hcl
provider "lvslb" {
  firewall_ip  = "192.168.0.1"
  port         = 9443
  https        = true
  insecure     = true
  vault_enable = true
}
```

## Argument Reference

* **firewall_ip** : (Required) IP for firewall API (lvslb-api)
* **port** : (Optional) [Def: 8080] Port for firewall API (lvslb-api)
* **https** : (Optional) [Def: false] Use HTTPS for firewall API
* **insecure** : (Optional) [Def: false] Don't check certificate for HTTPS
* **login** : (Optional) [Def: ""] User for http basic authentication
* **password** : (Optional) [Def: ""] Password for http basic authentication
* **vault_enable** : (Optional) [Def: false] Read login/password in secret/$vault_path/$firewall_ip or secret/$vault_path/$vault_key  
(For server and token, read environnement variables "VAULT_ADDR", "VAULT_TOKEN") Conflict With login/password
* **vault_path** : (Optional) [Def: "lvs"] Path where the key are
* **vault_key** : (Optional) [Def: ""] Name of key in vault path
