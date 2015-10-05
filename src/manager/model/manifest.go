package model

type Manifest struct {
	Autonomous string `json:"autonomous"`
	Ports string `json:"ports"`
	Controls string `json:"controls"`
	Settings string `json:"settings"`
	Logs string `json:"logs"`
	Binary string `json:"binary"`
}