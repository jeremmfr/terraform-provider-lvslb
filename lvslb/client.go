package lvslb

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Client = provider configuration
type Client struct {
	FirewallIP string
	Port       int
	HTTPS      bool
	Insecure   bool
	Logname    string
	Login      string
	Password   string
}

type ipvs struct {
	IP                 string       `json:"IP"`
	Port               string       `json:"Port"`
	Protocol           string       `json:"Protocol"`
	DelayLoop          string       `json:"Delay_loop"`
	LbAlgo             string       `json:"Lb_algo"`
	LbKind             string       `json:"Lb_kind"`
	PersistenceTimeout string       `json:"Persistence_timeout"`
	SorryIP            string       `json:"Sorry_IP"`
	SorryPort          string       `json:"Sorry_port"`
	Backends           ipvsBackends `json:"Backends"`
	Virtualhost        string       `json:"Virtualhost"`
	MonPeriod          string       `json:"Mon_period"`
}

type ipvsBackend struct {
	IP               string `json:"IP"`
	Port             string `json:"Port"`
	Weight           string `json:"Weight"`
	CheckType        string `json:"Check_type"`
	CheckPort        string `json:"Check_port"`
	CheckTimeout     string `json:"Check_timeout"`
	NbGetRetry       string `json:"Nb_get_retry"`
	DelayBeforeRetry string `json:"Delay_before_retry"`
	URLPath          string `json:"Url_path"`
	URLDigest        string `json:"Url_digest"`
	URLStatusCode    string `json:"Url_status_code"`
	MiscPath         string `json:"Misc_path"`
}

type ipvsBackends []ipvsBackend

// NewClient configure
func NewClient(firewallIP string, firewallPort int, https bool, insecure bool, logname string, login string, password string) *Client {
	client := &Client{
		FirewallIP: firewallIP,
		Port:       firewallPort,
		HTTPS:      https,
		Insecure:   insecure,
		Logname:    logname,
		Login:      login,
		Password:   password,
	}
	return client
}

func (client *Client) newRequest(uri string, ipvs *ipvs) (int, string, error) {
	urlString := "http://" + client.FirewallIP + ":" + strconv.Itoa(client.Port) + uri + "?&logname=" + client.Logname
	if client.HTTPS {
		urlString = strings.Replace(urlString, "http://", "https://", -1)
	}
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(ipvs)
	if err != nil {
		return 500, "", err
	}
	req, err := http.NewRequest("POST", urlString, body)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	if client.Login != "" && client.Password != "" {
		req.SetBasicAuth(client.Login, client.Password)
	}
	if err != nil {
		return 500, "", err
	}
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	if client.Insecure {
		tr = &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}
	}
	httpClient := &http.Client{Transport: tr}
	log.Printf("[DEBUG] Request API (%v) %v", urlString, body)
	resp, err := httpClient.Do(req)
	if err != nil {
		return 500, "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 500, "", err
	}
	log.Printf("[DEBUG] Response API (%v) %v => %v", urlString, resp.StatusCode, string(respBody))
	return resp.StatusCode, string(respBody), nil
}

func (client *Client) requestAPI(action string, ipvsSend *ipvs) (ipvs, error) {
	var ipvsReturn ipvs
	switch action {
	case "ADD":
		uriString := "/add_ipvs/" + ipvsSend.Protocol + "/" + ipvsSend.IP + "/" + ipvsSend.Port + "/"
		statuscode, body, err := client.newRequest(uriString, ipvsSend)
		if err != nil {
			return ipvsReturn, err
		}
		if statuscode == 401 {
			return ipvsReturn, fmt.Errorf("you are Unauthorized")
		}
		if statuscode != 200 {
			return ipvsReturn, fmt.Errorf(body)
		}
		return ipvsReturn, nil
	case "REMOVE":
		uriString := "/remove_ipvs/" + ipvsSend.Protocol + "/" + ipvsSend.IP + "/" + ipvsSend.Port + "/"
		statuscode, body, err := client.newRequest(uriString, ipvsSend)
		if err != nil {
			return ipvsReturn, err
		}
		if statuscode == 401 {
			return ipvsReturn, fmt.Errorf("you are Unauthorized")
		}
		if statuscode != 200 {
			return ipvsReturn, fmt.Errorf(body)
		}
		return ipvsReturn, nil
	case "CHECK":
		uriString := "/check_ipvs/" + ipvsSend.Protocol + "/" + ipvsSend.IP + "/" + ipvsSend.Port + "/"
		statuscode, body, err := client.newRequest(uriString, ipvsSend)
		if err != nil {
			return ipvsReturn, err
		}
		if statuscode == 401 {
			return ipvsReturn, fmt.Errorf("you are Unauthorized")
		}
		if statuscode == 404 {
			ipvsReturn.IP = nullStr
			ipvsReturn.Protocol = nullStr
			ipvsReturn.Port = nullStr
			return ipvsReturn, nil
		}

		errDecode := json.Unmarshal([]byte(body), &ipvsReturn)
		if errDecode != nil {
			return ipvsReturn, fmt.Errorf("[ERROR] decode json API response (%v) %v", errDecode, body)
		}
		return ipvsReturn, nil
	case "CHANGE":
		uriString := "/change_ipvs/" + ipvsSend.Protocol + "/" + ipvsSend.IP + "/" + ipvsSend.Port + "/"
		statuscode, body, err := client.newRequest(uriString, ipvsSend)
		if err != nil {
			return ipvsReturn, err
		}
		if statuscode == 401 {
			return ipvsReturn, fmt.Errorf("you are Unauthorized")
		}
		if statuscode != 200 {
			return ipvsReturn, fmt.Errorf(body)
		}
		return ipvsReturn, nil
	}
	return ipvsReturn, fmt.Errorf("internal error => unknown action for requestAPI")
}
