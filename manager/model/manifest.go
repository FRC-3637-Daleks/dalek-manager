package model

type Manifest struct {
	Templates struct {
				  Configs struct {
							  Ports     string `json:"ports"`
							  Controls  string `json:"controls"`
							  Settings  string `json:"settings"`
							  LogConfig string `json:"log-config"`
						  } `json:"configs"`
			  } `json:"templates"`
	Runtime   struct {
				  Configs    struct {
								 Ports    string `json:"ports"`
								 Controls string `json:"controls"`
								 Settings string `json:"settings"`
								 Logs     string `json:"logs"`
							 } `json:"configs"`
				  Autonomous string `json:"autonomous"`
				  Binary     string `json:"binary"`
			  } `json:"runtime"`
	Server    struct {
				  WebRoot  string `json:"web-root"`
				  Port     int     `json:"port"`
				  MqttHost string `json:"mqttHost"`
				  MqttPort int `json:"mqttPort"`
			  } `json:"server"`
}