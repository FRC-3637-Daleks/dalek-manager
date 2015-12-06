package model

type Manifest struct {
	Templates struct {
				  Configs struct {
							  Ports     string `json:"ports"`
							  Controls  struct {
											Available    string `json:"available"`
											Requirements string `json:"requirements"`
										}
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
				  WebRoot string `json:"web-root"`
				  Port    int     `json:"port"`
			  } `json:"server"`
}