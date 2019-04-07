package main

import "net"

var (
	ElinkTypeKeyNgReq    = "keyngreq"
	ElinkTypeKeyNgAck    = "keyngack"
	ElinkTypeDH          = "dh"
	ElinkTypeDevReg      = "dev_reg"
	ElinkTypeKeepAlive   = "keepalive"
	ElinkTypeAck         = "ack"
	ElinkTypeCfg         = "cfg"
	ElinkTypeGetStatus   = "get_status"
	ElinkTypeStatus      = "status"
	ElinkTypeRealDevInfo = "real_devinfo"
)

type ElinkSession struct {
	net.Conn    `json:"-"`
	UID         string        `json:"uid"`
	Mac         string        `json:"mac"`
	PerMac      string        `json:"peermac"`
	Seq         *Seq          `json:"seq"`
	DevData     DevRegData    `json:"devdata"`
	Status      StatusData    `json:"status"`
	RealDevInfo []RealDevData `json:"realdevinfo"`
	key         []byte        `json:"-"`
}

type Seq struct {
	RecvSeq int `json:"recvseq"`
	SendSeq int `json:"seq"`
}

type Base struct {
	Type string `json:"type"`
	Seq  int    `json:"sequence"`
	Mac  string `json:"mac"`
}

type KeyNgReq struct {
	Base
	Version     string `json:"version"`
	KeyModeList []KeyMode
}

type KeyNgAck struct {
	Base
	Mode string `json:"keymode"`
}

type KeyMode struct {
	KeyMode string `json:"keymode"`
}

type DH struct {
	Base
	Data DHData `json:"data"`
}

type DHData struct {
	DHKey string `json:"dh_key"`
	DHP   string `json:"dh_p"`
	DHG   string `json:"dh_g"`
}

type DevReg struct {
	Base
	Data DevRegData `json:"data"`
}

type DevRegData struct {
	Vendor    string `json:"vendor"`
	Model     string `json:"model"`
	SwVersion string `json:"swversion"`
	HdVersion string `json:"hdversion"`
	DespURL   string `json:"url"`
	Wireless  string `json:"wireless"`
}

/*
{
    "type": "get_status",
    "sequence": 123,
    "mac": "mac",
    "get": [
      {
        "name": "WiFi"
      },
      {
        "name": " WiFiswitch"
      },
      {
        "name": "ledswitch"
      },
      {
        "name": " WiFitimer"
      }
    ]
  }
*/
var (
	StatusNameWiFi       = "wifi"
	StatusNameWiFiSwitch = "wifiswitch"
	StatusNameLedSwitch  = "ledswitch"
	StatusNameWiFiTimer  = "wifitimer"

	APDataNameCPURate       = "cpurate"
	APDataNameMem           = "memoryuserate"
	APDataNameUploadSpeed   = "uploadspeed"
	APDataNameDownloadSpeed = "downloadspeed"
	APDataNameOnlineTime    = "onlineTime"
	APDataNameNum           = "terminalNum"
	APDataNameChannel       = "channel"
	APDataNameLoad          = "load"
)

type GetStatus struct {
	Base
	Get []Get `json:"get"`
}
type Get struct {
	Name string `json:"name"`
}

/*
{
    "type": "status",
    "sequence": 123,
    "mac": "mac",
    "status": {
      "WiFi": [
        {
          "radio": {
            "mode": "2.4G",
            "channel": 123
          },
          "ap": [
            {
              "apidx": 123,
              "enable": "yes",
              "ssid": "ssid",
              "key": " WiFi key",
              "auth": "auth mode",
              "encrypt": "encrypt mode"
            }
          ]
        }
      ],
      "WiFiswitch": {
        "status": "status"
      },
      "ledswitch": {
        "status": "status"
      },
      "WiFitimer": [
        {
          "weekday": "day",
          "time": "time",
          "enable": "enable"
        }
      ]
    }
  }
*/

type Status struct {
	Base
	Status StatusData `json:"status"`
}

type Radio struct {
	Mode    string `json:"mode"`
	Channel int    `json:"channel"`
}

type Ap struct {
	Apidx   int    `json:"apidx"`
	Enable  string `json:"enable"`
	Ssid    string `json:"ssid"`
	Key     string `json:"key"`
	Auth    string `json:"auth"`
	Encrypt string `json:"encrypt"`
}

type WiFi struct {
	Radio Radio `json:"radio"`
	Ap    []Ap  `json:"ap"`
}

type WiFiswitch struct {
	Status string `json:"status"`
}

type Ledswitch struct {
	Status string `json:"status"`
}

type Wpsswitch struct {
	Status string `json:"status"`
}

type WiFitimer struct {
	Weekday string `json:"weekday"`
	Time    string `json:"time"`
	Enable  string `json:"enable"`
}

type StatusData struct {
	WiFi       []WiFi      `json:"wifi"`
	WiFiswitch WiFiswitch  `json:"wifiswitch"`
	Ledswitch  Ledswitch   `json:"ledswitch"`
	Wpsswitch  Wpsswitch   `json:"wpsswitch"`
	WiFitimer  []WiFitimer `json:"wifitimer"`
}

/*
{
    "type": "cfg",
    "sequence": 123,
    "mac": "mac",
    "set": {
      "WiFiswitch": {
        "status": "status"
      },
      "ledswitch": {
        "status": "status"
      },
      "WiFitimer": [
        {
          "weekday": "day",
          "time": "time",
          "enable": "enable"
        }
      ]
    }
  }
*/

type SwitchConfig struct {
	Base
	Set SwitchSet `json:"set"`
}

type SwitchSet struct {
	WiFiswitch WiFiswitch  `json:"wifiswitch"`
	Ledswitch  Ledswitch   `json:"ledswitch"`
	Wpsswitch  Wpsswitch   `json:"wpsswitch"`
	WiFitimer  []WiFitimer `json:"wifitimer"`
}

type UpgradeConfig struct {
	Base
	Set UpgradeSet `json:"set"`
}

type UpgradeSet struct {
	Upgrade Upgrade `json:"upgrade"`
}

type Upgrade struct {
	Downloadurl string `json:"downurl"`
	IsReboot    string `json:"isreboot"`
}

/*
{
    "type": "cfg",
    "sequence": 123,
    "mac": "mac",
    "status": {
      "WiFi": [
        {
          "radio": {
            "mode": "2.4G",
            "channel": 123
          }
        }
      ]
    },
    "set": {
      "WiFi": [
        {
          "radio": {
            "mode": "2.4G",
            "channel": 123
          },
          "ap": [
            {
              "apidx": 123,
              "enable": "yes",
              "ssid": "ssid",
              "key": "WiFi key",
              "auth": "auth mode",
              "encrypt": "encrypt mode"
            }
          ]
        }
      ]
    }
  }
*/

type WifiConfig struct {
	Base
	Status StatusData `json:"status"`
	Set    WiFiSet    `json:"set"`
}

type WiFiSet struct {
	WiFi []WiFi `json:"wifi"`
}

type Reboot struct {
	Base
	Set RebootSet `json:"set"`
}
type RebootSet struct {
	Ctrlcommand string `json:"ctrlcommand"`
}

type RealDevInfo struct {
	Base
	RealDev []RealDevData `json:"real_devinfo"`
}

type RealDevData struct {
	Mac           string `json:"mac"`
	Rssi          string `json:"rssi"`
	OnlineTime    string `json:"onlineTime"`
	UploadSpeed   string `json:"uploadSpeed"`
	DownloadSpeed string `json:"downloadSpeed"`
	Band          string `json:"band"`
}

type KeepAlive struct {
	Type string `json:"type"`
	Seq  int    `json:"sequence"`
	Mac  string `json:"mac"`
}

type Ack struct {
	Type string `json:"type"`
	Seq  int    `json:"sequence"`
	Mac  string `json:"mac"`
}
