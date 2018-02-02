package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
	Session *mgo.Session
	DB      *mgo.Database
	DBName  string
	Addr    string
	User    string
	Pass    string
	Cron    *cron.Cron
}

func NewMongoDB(addr string, dbName string, user string, pass string) *MongoDB {
	info := &mgo.DialInfo{
		Addrs:    []string{addr},
		Timeout:  60 * time.Second,
		Database: dbName,
		Username: user,
		Password: pass,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatalf("Connect to db failed. Address:%s DB:%s user:%s info:+%v err:%s\n", addr, dbName, user, info, err)
	}
	/*session, err := mgo.Dial(addr)
	if err != nil {
		log.Fatalf("Connect to db failed. Address:%s DB:%s user:%s err:%s\n", addr, dbName, user, err)
	}*/

	db := session.DB(dbName)
	if user != "" {
		err := db.Login(user, pass)
		if err != nil {
			log.Fatalf("Login to db failed. Address:%s DB:%s user:%s password:%s err:%s\n", addr, dbName, user, pass, err)
		}
	}
	log.Printf("Connect to db succesful. Address:%s DB:%s user:%s\n", addr, dbName, user)
	var mongoDb MongoDB
	mongoDb.Session = session
	mongoDb.DB = db
	mongoDb.DBName = dbName
	mongoDb.Addr = addr
	mongoDb.User = user
	mongoDb.Pass = pass
	return &mongoDb
}

func (db *MongoDB) Close() {
	if db.Session != nil {
		db.Session.Close()
	}
}
func (db *MongoDB) AddServer() (err error) {
	if db.Session == nil {
		return fmt.Errorf("DB Session has nil\n")
	}
	//copy new seesion to concurrent using.
	childSession := db.Session.Copy()
	defer childSession.Close()
	childDB := childSession.DB(db.DBName)
	c := childDB.C("Traffic-Server")
	count, err := c.Find(bson.M{"server": MyIP}).Count()
	if count > 0 {
		err = c.Update(bson.M{"server": MyIP}, bson.M{"$set": bson.M{"time": time.Now()}})
		if err != nil {
			return err
		}
	} else {
		// Insert Datas
		err = c.Insert(bson.M{"server": MyIP, "time": time.Now()})
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *MongoDB) AddTrafficList() (err error) {
	if db.Session == nil {
		return fmt.Errorf("DB Session has nil\n")
	}
	//copy new seesion to concurrent using.
	childSession := db.Session.Copy()
	defer childSession.Close()
	childDB := childSession.DB(db.DBName)
	c := childDB.C(fmt.Sprintf("Traffic-%s", MyIP))
	// Insert Datas
	var trafficList WriteTrafficList
	trafficList.Time = time.Now()
	trafficList.TrafficList = TM.List
	err = c.Insert(&trafficList)

	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) AddGlobalTraffic(traffic *Traffic) (err error) {
	if db.Session == nil {
		return fmt.Errorf("DB Session has nil\n")
	}
	//copy new seesion to concurrent using.
	childSession := db.Session.Copy()
	defer childSession.Close()
	childDB := childSession.DB(db.DBName)
	c := childDB.C(fmt.Sprintf("Traffic-Global-%s", MyIP))
	// Insert Datas
	var dbTraffic DBTraffic
	dbTraffic.Time = time.Now()
	dbTraffic.Traffic = traffic
	err = c.Insert(&dbTraffic)

	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) AddTraffic(traffic *Traffic) (err error) {
	if db.Session == nil {
		return fmt.Errorf("DB Session has nil\n")
	}
	//copy new seesion to concurrent using.
	childSession := db.Session.Copy()
	defer childSession.Close()
	childDB := childSession.DB(db.DBName)
	c := childDB.C(fmt.Sprintf("Traffic-%s-%s", traffic.Host, MyIP))
	// Insert Datas
	err = c.Insert(&DBTraffic{Time: time.Now(), Traffic: traffic})

	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) GetAllServers() ([]string, error) {
	if db.Session == nil {
		return nil, fmt.Errorf("DB Session has nil\n")
	}
	//copy new seesion to concurrent using.
	childSession := db.Session.Copy()
	defer childSession.Close()
	childDB := childSession.DB(db.DBName)
	c := childDB.C("Traffic-Server")
	var servers []string
	err := c.Find(bson.M{}).All(&servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (db *MongoDB) GetTrafficList(table string, start time.Time, end time.Time) ([]DBTraffic, error) {
	var trafficList []DBTraffic
	if db.Session == nil {
		return trafficList, fmt.Errorf("DB Session has nil\n")
	}
	//copy new seesion to concurrent using.
	childSession := db.Session.Copy()
	defer childSession.Close()
	childDB := childSession.DB(db.DBName)
	c := childDB.C(table)

	err := c.Find(bson.M{"time": bson.M{"$gte": start, "$lte": end}}).All(&trafficList)
	if err != nil {
		return trafficList, err
	}
	log.Println(trafficList)
	return trafficList, nil
}

func (db *MongoDB) GetAllTrafficList(start time.Time, end time.Time) (map[string][]DBTraffic, error) {
	if db.Session == nil {
		return nil, fmt.Errorf("DB Session has nil\n")
	}
	//copy new seesion to concurrent using.
	childSession := db.Session.Copy()
	defer childSession.Close()
	childDB := childSession.DB(db.DBName)
	c := childDB.C("Traffic-Server")
	var servers []DBServer
	err := c.Find(bson.M{}).All(&servers)
	if err != nil {
		return nil, err
	}
	trafficList := make(map[string][]DBTraffic)
	if len(servers) > 0 {
		for _, v := range servers {
			list, err := db.GetTrafficList(fmt.Sprintf("Traffic-Server-%s", v.Server), start, end)
			if err != nil {
				log.Printf("[GetAllTrafficList][%s][Error]:%s\n", v.Server, err.Error())
				continue
			}
			trafficList[v.Server] = list
		}
	}
	return trafficList, nil
}

func (db *MongoDB) GetAllGlobalTrafficList(start time.Time, end time.Time) (map[string][]DBTraffic, error) {
	if db.Session == nil {
		return nil, fmt.Errorf("DB Session has nil\n")
	}
	//copy new seesion to concurrent using.
	childSession := db.Session.Copy()
	defer childSession.Close()
	childDB := childSession.DB(db.DBName)
	c := childDB.C("Traffic-Server")
	var servers []DBServer
	err := c.Find(bson.M{}).All(&servers)
	if err != nil {
		return nil, err
	}
	trafficList := make(map[string][]DBTraffic)
	if len(servers) > 0 {
		for _, v := range servers {
			list, err := db.GetTrafficList(fmt.Sprintf("Traffic-Global-%s", v.Server), start, end)
			if err != nil {
				log.Printf("[GetAllTrafficList][%s][Error]:%s\n", v.Server, err.Error())
				continue
			}
			trafficList[v.Server] = list
		}
	}
	return trafficList, nil
}
