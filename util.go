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

func getMyIp() (ipList []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return
	}
	var data []string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				data = append(data, ipnet.IP.String())
			}
		}
	}
	return data
}

func SaveFile(name string, data []byte) error {
	var file *os.File
	var err error
	// Create file
	if file, err = os.Create(name); err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		log.Println("WriteToStreamFile Error:", err)
		return err
	}

	// Sync to disk
	if err = file.Sync(); err != nil {
		file.Close()
		return err
	}

	file.Close()
	return nil
}

func WriteToFile(filename string, data string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		return err
	}
	return nil
}

func LoadFile(fileName string) (interface{}, error) {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Printf("Load file error: %v\n", e)
		os.Exit(1)
	}

	var data interface{}
	err := json.Unmarshal(file, &data)
	if err != nil {
		log.Printf("load file error:%v \n", err)
		return data, err
	}
	return data, nil
}

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
