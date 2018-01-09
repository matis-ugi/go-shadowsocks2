package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"syscall"
	"time"
)

func LoadConfigFile(fileName string) (Configs, error) {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Printf("Load config file error: %v\n", e)
		os.Exit(1)
	}

	var config Configs
	err := json.Unmarshal(file, &config)
	if err != nil {
		log.Printf("Config load error:%v \n", err)
		return config, err
	}
	return config, nil
}

func ConfigWatcher() {
	file, err := os.Open(configPath) // For read access.
	if err != nil {
		log.Println("configWatcher error:", err)
	}
	info, err := file.Stat()
	if err != nil {
		log.Println("configWatcher error:", err)
	}
	if modTime.Unix() == -62135596800 {
		modTime = info.ModTime()
	}

	if info.ModTime() != modTime {
		log.Printf("Config file changed. Reolad config file.\n")
		modTime = info.ModTime()
		CONFIGS, err = LoadConfigFile(configPath)
		if err != nil {
			log.Printf("configWatcher error:%v\n", err)
		}
	}
	defer file.Close()
}

func IPv4Verify(ip string) (bool, error) {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false, fmt.Errorf("%v is not an IPv4 address\n", trial)
	}
	return true, nil
}

func SetUlimit(number uint64) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Getting Rlimit ", err)
	}
	rLimit.Max = number
	rLimit.Cur = number
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Setting Rlimit ", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Getting Rlimit ", err)
	}
	log.Println("set file limit done:", rLimit)
}

func decodeString(data string) string {
	keybyte, err := hex.DecodeString(string(data[:2]))
	if err != nil {
		return ""
	}
	key := uint(keybyte[0])
	decodeStr := ""
	if len(data)%2 != 0 {
		return ""
	}
	for i := 2; i < len(data)-1; i += 2 {
		value := data[i : i+2]
		byte, err := hex.DecodeString(value)
		if err != nil {
			continue
		}
		decodeStr += string(uint(byte[0]) ^ key)
	}
	return decodeStr

}

func encodeString(data string, key int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	if key == 0 || key > 99 {
		key = r.Intn(99)
	}
	encodeStr := fmt.Sprintf("%02s", fmt.Sprintf("%x", key))

	for _, v := range data {
		encodeStr += fmt.Sprintf("%02s", fmt.Sprintf("%x", int(v)^key))
	}
	return encodeStr
}
