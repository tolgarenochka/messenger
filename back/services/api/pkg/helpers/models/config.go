package models

type Configuration struct {
	PathToServerKey string `json:"path_to_server_key"`
	PathToServerCrt string `json:"path_to_server_crt"`
	HTTPPort        string `json:"http_port"`
	HTTPSPort       string `json:"https_port"`
	WSPort          string `json:"ws_port"`
	WSSPort         string `json:"wss_port"`
	IP              string `json:"ip"`
	Salt            string `json:"salt"`
}
