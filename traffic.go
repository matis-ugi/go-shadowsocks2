package main

import (
	"fmt"
	"sync"

	"github.com/robfig/cron"
)

type TrafficManager struct {
	Mutex *sync.RWMutex
	List  map[string]*Traffic
}

type Traffic struct {
	Host     string `json:"host"`
	Inbound  int64  `json:"inbound"`
	Outbound int64  `json:"outbound"`
	Count    int64  `json:"count"`
}

func (t *Traffic) Reset() {
	t.Count = 0
	t.Inbound = 0
	t.Outbound = 0
}

func NewTrafficManager() *TrafficManager {
	var tm TrafficManager
	tm.Mutex = &sync.RWMutex{}
	tm.List = make(map[string]*Traffic)
	go tm.CronJob()
	return &tm
}
func (tm *TrafficManager) CronJob() {
	c := cron.New()
	c.AddFunc(CONFIGS.RecordTime, tm.SaveTraffic)
	c.Start()
}

func (tm *TrafficManager) CalcGlobalTraffic() Traffic {
	var t Traffic
	for _, v := range tm.List {
		t.Count += v.Count
		t.Inbound += v.Inbound
		t.Outbound += v.Outbound
	}
	return t
}
func (tm *TrafficManager) SaveTraffic() {
	if len(tm.List) > 0 {
		tm.Mutex.Lock()
		global := tm.CalcGlobalTraffic()
		mongo.AddGlobalTraffic(&global)
		mongo.AddTrafficList()
		for _, v := range tm.List {
			mongo.AddTraffic(v)
		}
		/*data, err := json.Marshal(tm.List)
		if err != nil {
			log.Println("SaveTraffic Error:", err)
		}*/
		for _, v := range tm.List {
			v.Reset()
		}
		tm.Mutex.Unlock()
		//SaveFile(fmt.Sprintf("record/%s.json", t.Format("2006-01-02-15:04")), data)
	}
}

func (tm *TrafficManager) Add(host string, inbound int64, outbound int64, count int) {
	//log.Println(host, inbound, outbound, count)
	tm.Mutex.Lock()
	if v, ok := tm.List[host]; !ok {
		var traffic Traffic
		traffic.Host = host
		if inbound > 0 {
			traffic.Inbound += (inbound)
		}
		if outbound > 0 {
			traffic.Outbound += (outbound)
		}
		if count > 0 {
			traffic.Count += int64(count)
		}
		tm.List[host] = &traffic
	} else {
		if inbound > 0 {
			v.Inbound += int64(inbound)
		}
		if outbound > 0 {
			v.Outbound += int64(outbound)
		}
		if count > 0 {
			v.Count += int64(count)
		}
		fmt.Println(host, v)
	}

	tm.Mutex.Unlock()
}

func (tm *TrafficManager) Get(host string) *Traffic {
	tm.Mutex.RLock()
	var traffic *Traffic
	if v, ok := tm.List[host]; ok {
		traffic = v
	}
	tm.Mutex.RUnlock()
	return traffic
}
