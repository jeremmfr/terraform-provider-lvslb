# terraform-provider-lvslb
[![GoDoc](https://godoc.org/github.com/jeremmfr/terraform-provider-lvslb?status.svg)](https://godoc.org/github.com/jeremmfr/terraform-provider-lvslb) [![Go Report Card](https://goreportcard.com/badge/github.com/jeremmfr/terraform-provider-lvslb)](https://goreportcard.com/report/github.com/jeremmfr/terraform-provider-lvslb)

Terraform's provider for generate keepalived virtual_server with lvslb-api (https://github.com/jeremmfr/lvslb-api)

Compile:
========

export GO111MODULE=on
go build -o terraform-provider-lvslb && mv terraform-provider-lvslb /usr/bin/

Config:
=======

Setup information for contact server :
```
provider "lvsnetwork" {
    firewall_ip = "192.168.0.1"
    port		= 9443
	https		= true
	insecure	= true
	vault_enable = true
}
```
* **firewall_ip** : (Required) IP for firewall API (lvsnetwork-api)
* **port** : (Optional) [Def: 8080] Port for firewall API (lvsnetwork-api)
* **https** : (Optional) [Def: false] Use HTTPS for firewall API
* **insecure** : (Optional) [Def: false] Don't check certificate for HTTPS
* **login** : (Optional) [Def: ""] User for http basic authentication
* **password** : (Optional) [Def: ""] Password for http basic authentication
* **vault_enable** : (Optional) [Def: false] Read login/password in secret/$vault_path/$firewall_ip or secret/$vault_path/$vault_key (For server and token, read environnement variables "VAULT_ADDR", "VAULT_TOKEN") Priority on login/password
* **vault_path** : (Optional) [Def: "lvs"] Path where the key are
* **vault_key** : (Optional) [Def: ""] Name of key in vault path

Resource:
=========

** ipvs **
---------------

Create ipvs lb  configuration on server with backend

```
resource lvslb_ipvs "test" {
        ip = "203.0.113.1"
        port = 80
        backends {
                ip = [ "10.0.0.129", "10.0.0.130" ]
                weight = 2
        }
        backends {
                ip = [ "10.0.0.132" ]
        }
}
```

* **ip** : (Required) IPv4 for load balancer
* **port** : (Required) Port for load balancer
* **protocol** : (Optional) [Def: "TCP"] Protocol for load balancer (TCP|UDP|SCTP)
* **type** : (Optional) [Def: "NAT"] Type for load balancer (NAT|DR|TUN)
* **algo** : (Optional) [Def: "wlc"] Algorithm for load balancer (wlc|lc|rr|wrr|lblc|sh|dh)
* **persistence_timeout** : (Optional) [Def: 0 ] Persistence for choice backend compared client IP
* **timer_check** : (Optional) [Def: 5 ] number of secondes between healthcheck
* **sorry_server_ip** : (Optional) IP of sorry server if all backend is out of pool
* **sorry_server_port** : (Optional) Port of sorry server if all backend is out of pool
* **virtualhost** : (Optional) Vhost for healthchecker if HTTP_GET or SSL_GET
* **monitoring_period**: (Optional) Period options for add/change monitoring
* **backends** (Required) block supports :
	* **ip** : (Required) list of IP for backends
	* **port** : (Optional) [ Default: port of load balancer ] port of backends
	* **weight** : (Optional) [ Default: 1 ] weight for backends
	* **check_type** : (Optional) [ Default: "TCP_CHECK" ] Type of check for healthchecker (TCP_CHECK|HTTP_GET|SSL_GET|MISC_CHECK|NONE)
	* **check_port** : (Optional) [ Default: port of backends ] port for healthchecker if different of port backends
	* **check_timeout** : (Optional) [ Default: 3 ] timeout of secondes for healthchecker
	* **nb_get_retry** : (Optional) [ Default: 3 ] number of retry after healthcheck failed
	* **delay_before_retry** : (Optional) [ Default: 3 ] number of secondes before new healthcheck after healthcheck failed
	* **check_url** : (Optional) Url for healthchecker when type is HTTP_GET or SSL_GET
	* **check_digest** : (Optional) md5sum of response when type is HTTP_GET or SSL_GET
	* **check_status_code** : (Optional) HTTP Code of response when type is HTTP_GET or SSL_GET
	* **misc_path** : (Optional) Path for script when type is MISC_CHECK
