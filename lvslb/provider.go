package lvslb

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const nullStr string = "null"
const (
	defaultFirewallPort = 8080
)

// Provider lvslb for terraform.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"firewall_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  defaultFirewallPort,
			},
			"https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"login": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_enable": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"login", "password"},
			},
			"vault_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "lvs",
			},
			"vault_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"lvslb_ipvs": resourceIpvs(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		firewallIP:   d.Get("firewall_ip").(string),
		firewallPort: d.Get("port").(int),
		https:        d.Get("https").(bool),
		insecure:     d.Get("insecure").(bool),
		logname:      os.Getenv("USER"),
		login:        d.Get("login").(string),
		password:     d.Get("password").(string),
		vaultEnable:  d.Get("vault_enable").(bool),
		vaultPath:    d.Get("vault_path").(string),
		vaultKey:     d.Get("vault_key").(string),
	}

	return config.Client()
}
