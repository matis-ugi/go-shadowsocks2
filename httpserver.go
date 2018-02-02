package main

import (
	_ "crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var WebRouter *mux.Router

func WebServer() {
	WebRouter = mux.NewRouter()
	WebRouter.HandleFunc("/monitor", monitorHandler)
	WebRouter.HandleFunc("/state/{host}", stateHandler)
	WebRouter.HandleFunc("/api/{cmd}", apiHandler)
	WebRouter.HandleFunc("/api/{cmd}/{host}/{start}/{end}", apiHandler)
	WebRouter.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("html/"))))

	log.Println("API Service starting.", CONFIGS.HTTP)
	//WebSocketInit()
	if (CONFIGS.HTTP) != "" {
		err := http.ListenAndServe(CONFIGS.HTTP, WebRouter)
		if err != nil {
			log.Fatal("Web Service start failed.", err)
		}
	}

}

func monitorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//vars := mux.Vars(r)
	jsonParser(TM.List, w)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data ResponseObject
	data.State = "fail"
	switch vars["cmd"] {
	case "login":
		var rUser RequestUser
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&rUser)
		if err != nil {
			data.Error = err.Error()
		}
		defer r.Body.Close()
		if rUser.Account != "" && rUser.Password != "" && rUser.Salt != "" {
			for _, v := range CONFIGS.UserList {
				if v.Account == rUser.Account && v.Password == rUser.Password {
					data.State = "success"
					data.Token = GenToken(rUser.Account, rUser.Salt)
					break
				}
			}
		} else {
			data.State = "fail"
			data.Error = "Some fields incorrect."
		}
		jsonParser(data, w)
	case "global_chart":
		if !VerifyUser(r) {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001]:User verify failed.")
			return
		}
		if !verifyFields([]string{"host", "start", "end"}, vars) {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001]:Fields verify failed.")
			return
		}
		start, err := TimeParser(vars["start"])
		if err != nil {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001]:Start Time field verify failed."+err.Error())
			return
		}
		end, err := TimeParser(vars["end"])
		if err != nil {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001]:End Time field verify failed."+err.Error())
			return
		}
		list, err := mongo.GetAllGlobalTrafficList(start, end)
		if err != nil {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001][GetAllGlobalTrafficList]:"+err.Error())
			return
		}
		jsonParser(list, w)
	case "chart":
		if !VerifyUser(r) {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001]:User verify failed.")
			return
		}
		if !verifyFields([]string{"host", "start", "end"}, vars) {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001]:Fields verify failed.")
			return
		}
		start, err := TimeParser(vars["start"])
		if err != nil {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001]:Start Time field verify failed."+err.Error())
			return
		}
		end, err := TimeParser(vars["end"])
		if err != nil {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001]:End Time field verify failed."+err.Error())
			return
		}
		list, err := mongo.GetAllTrafficList(start, end)
		if err != nil {
			w.WriteHeader(501)
			fmt.Fprint(w, "[1001][GetAllGlobalTrafficList]:"+err.Error())
			return
		}
		jsonParser(list, w)

	}

}

func VerifyUser(r *http.Request) bool {
	var ui UserInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ui)
	if err != nil {
		return false
	}
	defer r.Body.Close()
	if ui.Account != "" && ui.Salt != "" && ui.Token != "" {
		return VerifyToken(ui.Account, ui.Salt, ui.Token)
	}
	return false
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//vars := mux.Vars(r)
}

func jsonParser(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if data != nil {
		json, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(500)
			log.Println("Error generating json", err)
			fmt.Fprintln(w, "Could not generate JSON")
			return
		}
		fmt.Fprint(w, string(json))
	} else {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 no data can be find.")
	}
}

func verifyFields(fields []string, vars map[string]string) bool {
	verify := true
	for _, field := range fields {
		if v, ok := vars[field]; ok {
			if len(v) == 0 {
				verify = false
			}
		} else {
			verify = false
		}
	}
	return verify
}

func TimeParser(timeStr string) (time.Time, error) {
	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	newTime := time.Unix(i, 0)
	return newTime, nil
}
