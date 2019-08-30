package lvslb

import (
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Config provider
type Config struct {
	https        bool
	insecure     bool
	vaultEnable  bool
	firewallPort int
	firewallIP   string
	logname      string
	login        string
	password     string
	vaultPath    string
	vaultKey     string
}

// Client configures with Config
func (c *Config) Client() (*Client, error) {
	var client *Client
	if !c.vaultEnable {
		client = NewClient(c.firewallIP, c.firewallPort, c.https, c.insecure, c.logname, c.login, c.password)
	} else {
		login, password := getLoginVault(c.vaultPath, c.firewallIP, c.vaultKey)
		client = NewClient(c.firewallIP, c.firewallPort, c.https, c.insecure, c.logname, login, password)
	}

	return client, nil
}

func getLoginVault(path string, firewallIP string, key string) (string, string) {
	login := ""
	password := ""
	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return "", ""
	}

	c := client.Logical()
	if key != "" {
		secret, err := c.Read(strings.Join([]string{"/secret/", path, "/", key}, ""))
		if err != nil {
			return "", ""
		}
		if secret != nil {
			for key, value := range secret.Data {
				if key == "login" {
					login = value.(string)
				}
				if key == "password" {
					password = value.(string)
				}
			}
		}
	} else {
		secret, err := c.Read(strings.Join([]string{"/secret/", path, "/", firewallIP}, ""))
		if err != nil {
			return "", ""
		}
		if secret != nil {
			for key, value := range secret.Data {
				if key == "login" {
					login = value.(string)
				}
				if key == "password" {
					password = value.(string)
				}
			}
		}
	}
	return login, password
}
