package main

type Configs struct {
	Debug      bool   `json:"Debug"`
	HTTP       string `json:"HTTP"`
	Client     string `json:"Client"`
	Server     string `json:"Server"`
	Cipher     string `json:"Cipher"`
	Key        string `json:"Key"`
	Password   string `json:"Password"`
	Keygen     int    `json:"Keygen"`
	Socks      string `json:"Socks"`
	RedirTCP   string `json:"RedirTCP"`
	RedirTCP6  string `json:"RedirTCP6"`
	TCPTun     string `json:"TCPTun"`
	UDPTun     string `json:"UDPTun"`
	UDPTimeout int    `json:"UDPTimeout"`
	UDPSocks   bool   `json:"UDPSocks"`
}
