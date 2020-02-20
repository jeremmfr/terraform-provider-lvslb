package lvslb

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	one                     = 1
	maxInternetPort         = int(1<<16 - one)
	defaultTimerCheck       = 5
	maxTimerCheck           = 60
	defaultCheckTimeout     = 3
	maxCheckTimeout         = 60
	defaultNbGetRetry       = 3
	maxNbGetRetry           = 10
	defaultDelayBeforeRetry = 3
	maxDelayBeforeRetry     = 60
	maxPersistenceTimeout   = 86400
	maxBackendWeight        = 1000
	minStatusCode           = 100
	maxStatusCode           = 600
)

func resourceIpvs() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpvsCreate,
		Read:   resourceIpvsRead,
		Update: resourceIpvsUpdate,
		Delete: resourceIpvsDelete,

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					testInput := net.ParseIP(value)
					if testInput.To16() == nil {
						errors = append(errors, fmt.Errorf("[ERROR] %q %v isn't an IPv4 or IPv6", k, value))
					}
					return
				},
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, maxInternetPort),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TCP",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := strings.ToUpper(v.(string))
					if value != "TCP" && value != "UDP" && value != "SCTP" {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be TCP or UDP", k))
					}
					return
				},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NAT",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := strings.ToUpper(v.(string))
					if value != "NAT" && value != "DR" && value != "TUN" {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be NAT, DR or TUN", k))
					}
					return
				},
			},
			"algo": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "wlc",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := strings.ToLower(v.(string))
					if value != "wlc" && value != "lc" && value != "rr" && value != "wrr" &&
						value != "lblc" && value != "sh" && value != "dh" {
						errors = append(errors, fmt.Errorf("[ERROR] %q must be wlc, lc, rr, wrr, lblc, sh or dh", k))
					}
					return
				},
			},
			"persistence_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, maxPersistenceTimeout),
			},
			"timer_check": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultTimerCheck,
				ValidateFunc: validateIntegerInRange(one, maxTimerCheck),
			},
			"sorry_server_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					testInput := net.ParseIP(value)
					if testInput.To16() == nil {
						errors = append(errors, fmt.Errorf("[ERROR] %q %v isn't an IPv4 or IPv6", k, value))
					}
					return
				},
			},
			"sorry_server_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, maxInternetPort),
			},
			"virtualhost": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"monitoring_period": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"backends": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(0, maxInternetPort),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      one,
							ValidateFunc: validateIntegerInRange(one, maxBackendWeight),
						},
						"check_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "TCP_CHECK",
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := strings.ToUpper(v.(string))
								if value != "TCP_CHECK" && value != "HTTP_GET" && value != "SSL_GET" &&
									value != "MISC_CHECK" && value != "NONE" {
									errors = append(errors, fmt.Errorf("[ERROR] %q must be TCP_CHECK, HTTP_GET, SSL GET, MISC_CHECK or NONE", k))
								}
								return
							},
						},
						"check_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(one, maxInternetPort),
						},
						"check_timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      defaultCheckTimeout,
							ValidateFunc: validateIntegerInRange(one, maxCheckTimeout),
						},
						"nb_get_retry": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      defaultNbGetRetry,
							ValidateFunc: validateIntegerInRange(one, maxNbGetRetry),
						},
						"delay_before_retry": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      defaultDelayBeforeRetry,
							ValidateFunc: validateIntegerInRange(one, maxDelayBeforeRetry),
						},
						"check_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"check_digest": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"check_status_code": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(minStatusCode, maxStatusCode),
						},
						"misc_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceIpvsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := validateIPBackend(d)
	if err != nil {
		return err
	}
	Ipvs := createStrucIpvs(d)
	_, err = client.requestAPI("ADD", &Ipvs)
	if err != nil {
		return err
	}
	d.SetId(d.Get("ip").(string) + "_" + strings.ToUpper(d.Get("protocol").(string)) +
		"_" + strconv.Itoa(d.Get("port").(int)))
	return nil
}

func resourceIpvsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	Ipvs := createStrucIpvs(d)
	IpvsRead, err := client.requestAPI("CHECK", &Ipvs)
	if err != nil {
		return err
	}
	if IpvsRead.IP == nullStr {
		d.SetId("")
	}
	if len(IpvsRead.Backends) == 0 {
		emptyBackend := map[string]interface{}{
			"ip":                 make([]string, 0),
			"port":               0,
			"weight":             0,
			"check_type":         "",
			"check_port":         0,
			"check_timeout":      0,
			"nb_get_retry":       0,
			"delay_before_retry": 0,
			"check_url":          "",
			"check_digest":       "",
			"check_status_code":  0,
			"misc_path":          "",
		}
		tfErr := d.Set("backends", []map[string]interface{}{emptyBackend})
		if tfErr != nil {
			panic(tfErr)
		}
	}
	return nil
}

func resourceIpvsUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)
	client := m.(*Client)
	err := validateIPBackend(d)
	if err != nil {
		return err
	}
	switch {
	case d.HasChange("ip") || d.HasChange("port"):
		IpvsOld := createStrucIpvs(d)
		oldIP, _ := d.GetChange("ip")
		IpvsOld.IP = oldIP.(string)
		oldPort, _ := d.GetChange("port")
		IpvsOld.Port = strconv.Itoa(oldPort.(int))
		oldProtocol, _ := d.GetChange("protocol")
		IpvsOld.Protocol = strings.ToUpper(oldProtocol.(string))
		_, err := client.requestAPI("REMOVE", &IpvsOld)
		if err != nil {
			return err
		}
		d.SetId("")
		Ipvs := createStrucIpvs(d)
		_, err = client.requestAPI("ADD", &Ipvs)
		if err != nil {
			return err
		}
		d.SetId(d.Get("ip").(string) + "_" + strings.ToUpper(d.Get("protocol").(string)) +
			"_" + strconv.Itoa(d.Get("port").(int)))
	case d.HasChange("protocol"):
		oldProtocol, _ := d.GetChange("protocol")
		if strings.EqualFold(oldProtocol.(string), d.Get("protocol").(string)) {
			Ipvs := createStrucIpvs(d)
			_, err := client.requestAPI("CHANGE", &Ipvs)
			if err != nil {
				return err
			}
		} else {
			IpvsOld := createStrucIpvs(d)
			IpvsOld.Protocol = strings.ToUpper(oldProtocol.(string))
			_, err := client.requestAPI("REMOVE", &IpvsOld)
			if err != nil {
				return err
			}
			d.SetId("")
			Ipvs := createStrucIpvs(d)
			_, err = client.requestAPI("ADD", &Ipvs)
			if err != nil {
				return err
			}
			d.SetId(d.Get("ip").(string) + "_" + strings.ToUpper(d.Get("protocol").(string)) +
				"_" + strconv.Itoa(d.Get("port").(int)))
		}
	default:
		Ipvs := createStrucIpvs(d)
		_, err := client.requestAPI("CHANGE", &Ipvs)
		if err != nil {
			return err
		}
	}
	d.Partial(false)
	return nil
}

func resourceIpvsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	Ipvs := createStrucIpvs(d)
	_, err := client.requestAPI("REMOVE", &Ipvs)
	if err != nil {
		return err
	}
	return nil
}

func validateIPBackend(d *schema.ResourceData) error {
	if v, ok := d.GetOk("backends"); ok {
		backendSet := v.([]interface{})
		for _, dataBackend := range backendSet {
			backend := dataBackend.(map[string]interface{})
			for _, backendIP := range backend["ip"].([]interface{}) {
				testInputIP := net.ParseIP(d.Get("ip").(string))
				if testInputIP.To4() == nil {
					testInput := net.ParseIP(backendIP.(string))
					if testInput.To16() == nil || !strings.Contains(backendIP.(string), ":") {
						return fmt.Errorf("[ERROR] backend %v isn't an IPv6 for IPv6 virtual server", backendIP)
					}
				} else {
					testInput := net.ParseIP(backendIP.(string))
					if testInput.To4() == nil {
						return fmt.Errorf("[ERROR] backend %v isn't an IPv4 for IPv4 virtual server", backendIP)
					}
				}
			}
		}
	}
	return nil
}
func createStrucIpvs(d *schema.ResourceData) ipvs {
	var backends []ipvsBackend
	if v, ok := d.GetOk("backends"); ok {
		backendSet := v.([]interface{})
		for _, dataBackend := range backendSet {
			backend := dataBackend.(map[string]interface{})
			for _, backendIP := range backend["ip"].([]interface{}) {
				var backendPort string
				if backend["port"].(int) != 0 {
					backendPort = strconv.Itoa(backend["port"].(int))
				} else {
					backendPort = strconv.Itoa(d.Get("port").(int))
				}
				var checkPort string
				switch {
				case backend["check_port"] != 0:
					checkPort = strconv.Itoa(backend["check_port"].(int))
				case backend["port"].(int) != 0:
					checkPort = strconv.Itoa(backend["port"].(int))
				default:
					checkPort = strconv.Itoa(d.Get("port").(int))
				}
				var statusCode string
				if backend["check_status_code"].(int) != 0 {
					statusCode = strconv.Itoa(backend["check_status_code"].(int))
				} else {
					statusCode = ""
				}

				IpvsBackend := ipvsBackend{
					IP:               backendIP.(string),
					Port:             backendPort,
					Weight:           strconv.Itoa(backend["weight"].(int)),
					CheckType:        strings.ToUpper(backend["check_type"].(string)),
					CheckPort:        checkPort,
					CheckTimeout:     strconv.Itoa(backend["check_timeout"].(int)),
					NbGetRetry:       strconv.Itoa(backend["nb_get_retry"].(int)),
					DelayBeforeRetry: strconv.Itoa(backend["delay_before_retry"].(int)),
					URLPath:          backend["check_url"].(string),
					URLDigest:        backend["check_digest"].(string),
					URLStatusCode:    statusCode,
					MiscPath:         backend["misc_path"].(string),
				}
				backends = append(backends, IpvsBackend)
			}
		}
	}
	Ipvs := ipvs{
		IP:                 d.Get("ip").(string),
		Port:               strconv.Itoa(d.Get("port").(int)),
		Protocol:           strings.ToUpper(d.Get("protocol").(string)),
		DelayLoop:          strconv.Itoa(d.Get("timer_check").(int)),
		LbAlgo:             strings.ToLower(d.Get("algo").(string)),
		LbKind:             strings.ToUpper(d.Get("type").(string)),
		PersistenceTimeout: strconv.Itoa(d.Get("persistence_timeout").(int)),
		SorryIP:            d.Get("sorry_server_ip").(string),
		SorryPort:          strconv.Itoa(d.Get("sorry_server_port").(int)),
		Virtualhost:        d.Get("virtualhost").(string),
		MonPeriod:          d.Get("monitoring_period").(string),
		Backends:           backends,
	}
	return Ipvs
}
