package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
	"net/http"
	"os/exec"
)

type Options struct {
	Listenport []int `long:"listenport" env:"LISTENPORT" default:"2121" description:"Listening port"`
}

var (
	opts Options
)

type proxyresponse struct {
	Success  bool
	Errormsg error
	Response mvgresponse
}

type mvgresponse struct {
	Station            string
	Server_time        string
	Transports         []string
	Using_config_file  bool
	Result_sorted      []mvgresult
	Result_display     []mvgresult
	Station_unknown    bool
	Station_alternates []string
}

type mvgresult struct {
	Line        string
	Destination string
	Minutes     int32
}

func main() {

	log.Info("App started")

	flags.Parse(&opts)

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"requestURI":  r.RequestURI,
			"remoteAddr ": r.RemoteAddr,
		}).Error("failed to find page")
		// fmt.Println(r.RequestURI)
		http.ServeFile(w, r, "public/index.html")
	})

	r.HandleFunc("/query/{station}", Querymvg)

	log.WithFields(log.Fields{
		"port": opts.Listenport[0],
		"host": "0.0.0.0",
	}).Info("starting listender")

	erro := http.ListenAndServe(fmt.Sprintf(":%v", opts.Listenport[0]), r)

	log.Error("Failed to open socket:", erro)
}

func Querymvg(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	log.WithFields(log.Fields{
		"station": vars["station"],
	}).Info("Received request")

	// out, err := exec.Command("cat", "marienplazfailed.json").Output()
	out, err := exec.Command("mvdwg_json", vars["station"]).Output()
	if err != nil {
		log.Error(err)
	}

	var queryresponse proxyresponse
	var outdata *mvgresponse

	err = json.Unmarshal(out, &outdata)
	queryresponse.Response = *outdata

	if err != nil {
		queryresponse.Success = false
		log.WithFields(log.Fields{
			"station": vars["station"],
		}).Error("Response unmarshall failed: ", err)
	} else {
		queryresponse.Success = true
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)

	queryresponse.Errormsg = err
	js, _ := json.Marshal(queryresponse)
	resp.Write(js)
}
