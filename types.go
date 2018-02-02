package main

import "time"

type Configs struct {
	Debug      bool   `json:"Debug"`
	HTTP       string `json:"HTTP"`
	RecordTime string `json:"Record"`
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
	UserList   []User `json:"UserList"`
	DB         struct {
		Addr   string `json:"Address"`
		DbName string `json:"DBName"`
		User   string `json:"User"`
		Pass   string `json:"Password"`
	} `json:"DB"`
}

type User struct {
	Account  string   `json:"Account"`
	Password string   `json:"Password"`
	Role     []string `json:"Role"`
}

type RequestUser struct {
	Account  string `json:"account"`
	Password string `json:"pwd"`
	Salt     string `json:"salt"`
}

type UserInfo struct {
	Account string `json:"account"`
	Salt    string `json:"salt"`
	Token   string `json:"token"`
}

type ResponseObject struct {
	State string      `json:"State"`
	Token string      `json:"Token"`
	Data  interface{} `json:"Data"`
	Error string      `json:"error"`
}

type DBServer struct {
	Server string    `json:"server" bson:"server"`
	Time   time.Time `json:"time" bson:"time"`
}

type DBTraffic struct {
	Time    time.Time `json:"time" bson:"time"`
	Traffic *Traffic  `json:"traffic" bson:"traffic"`
}

type WriteTrafficList struct {
	Time        time.Time           `json:"time" bson:"time"`
	TrafficList map[string]*Traffic `json:"traffic" bson:"traffic"`
}

type DBTrafficList struct {
	Time        time.Time          `json:"time" bson:"time"`
	TrafficList map[string]Traffic `json:"traffic" bson:"traffic"`
}
