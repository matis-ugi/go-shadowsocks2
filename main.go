package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
)

var (
	CONFIGS    Configs
	modTime    time.Time
	configPath string = "configs.json"
	MyIP       string = "127.0.0.1"
	mongo      *MongoDB
	TM         *TrafficManager
)

func logf(f string, v ...interface{}) {
	if CONFIGS.Debug {
		log.Printf(f, v...)
	}
}

func main() {
	flag.StringVar(&configPath, "c", "configs.json", " confis json file path")
	flag.Parse()
	var err error
	CONFIGS, err = LoadConfigFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	//get local ip
	tmpIP := getMyIp()
	fmt.Println("Internal IP:", tmpIP)
	if len(tmpIP) > 0 {
		MyIP = tmpIP[0]
	}

	TM = NewTrafficManager()

	if CONFIGS.Keygen > 0 {
		key := make([]byte, CONFIGS.Keygen)
		io.ReadFull(rand.Reader, key)
		fmt.Println(base64.URLEncoding.EncodeToString(key))
		return
	}

	if CONFIGS.Client == "" && CONFIGS.Server == "" {
		flag.Usage()
		return
	}

	var key []byte
	if CONFIGS.Key != "" {
		k, err := base64.URLEncoding.DecodeString(CONFIGS.Key)
		if err != nil {
			log.Fatal(err)
		}
		key = k
	}

	if CONFIGS.Client != "" { // client mode
		addr := CONFIGS.Client
		cipher := CONFIGS.Cipher
		password := CONFIGS.Password
		var err error

		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				log.Fatal(err)
			}
		}

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			log.Fatal(err)
		}

		if CONFIGS.UDPTun != "" {
			for _, tun := range strings.Split(CONFIGS.UDPTun, ",") {
				p := strings.Split(tun, "=")
				go udpLocal(p[0], addr, p[1], ciph.PacketConn)
			}
		}

		if CONFIGS.TCPTun != "" {
			for _, tun := range strings.Split(CONFIGS.TCPTun, ",") {
				p := strings.Split(tun, "=")
				go tcpTun(p[0], addr, p[1], ciph.StreamConn)
			}
		}

		if CONFIGS.Socks != "" {
			socks.UDPEnabled = CONFIGS.UDPSocks
			go socksLocal(CONFIGS.Socks, addr, ciph.StreamConn)
			if CONFIGS.UDPSocks {
				go udpSocksLocal(CONFIGS.Socks, addr, ciph.PacketConn)
			}
		}

		if CONFIGS.RedirTCP != "" {
			go redirLocal(CONFIGS.RedirTCP, addr, ciph.StreamConn)
		}

		if CONFIGS.RedirTCP6 != "" {
			go redir6Local(CONFIGS.RedirTCP6, addr, ciph.StreamConn)
		}
	}

	if CONFIGS.Server != "" { // server mode
		//Start MongoDB
		mongo = NewMongoDB(CONFIGS.DB.Addr, CONFIGS.DB.DbName, CONFIGS.DB.User, CONFIGS.DB.Pass)
		mongo.AddServer()
		addr := CONFIGS.Server
		cipher := CONFIGS.Cipher
		password := CONFIGS.Password
		var err error

		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				log.Fatal(err)
			}
		}

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			log.Fatal(err)
		}

		go udpRemote(addr, ciph.PacketConn)
		go tcpRemote(addr, ciph.StreamConn)
		go WebServer()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func parseURL(s string) (addr, cipher, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return
}
