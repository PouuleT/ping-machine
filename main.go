package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/tatsushid/go-fastping"
	"github.com/urfave/negroni"
)

// PingResult represents a response to a ping
type PingResult struct {
	Response string        `json:"response"`
	IP       string        `json:"IP"`
	Time     time.Duration `json:"time"`
}

// PingError represents a problem to a ping
type PingError struct {
	Response string `json:"response"`
	Why      string `json:"Why"`
}

func main() {
	log.Print("Serving server on /ping")

	n := negroni.Classic()
	// Serve static files
	n.Use(negroni.NewStatic(http.Dir("public")))
	// Compress responses
	router := httprouter.New()
	router.GET("/ping/:times/times/:addr", HandlePing)

	n.UseHandler(router)
	n.Run(":8080")
}

func ping(addr string, nbToPing int) ([]*PingResult, error) {
	responses := []*PingResult{}

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", addr)
	if err != nil {
		log.Print("got an error whil resolving ", err)
		return responses, err
	}
	var resp *PingResult

	p.AddIPAddr(ra)
	p.MaxRTT = time.Millisecond * 100
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		log.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		resp = &PingResult{"OK", addr.String(), rtt}
		responses = append(responses, resp)
	}

	p.OnIdle = func() {
		if resp == nil {
			resp = &PingResult{"KO", addr, time.Millisecond * 0}
			responses = append(responses, resp)
		}
	}

	for i := 0; i < nbToPing; i++ {
		err = p.Run()
		if err != nil {
			log.Println("error:", err)
			return responses, err
		}
		resp = nil
	}

	return responses, err
}

// HandlePing is a Route function to ping an ip address
func HandlePing(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	addrToPing := ps.ByName("addr")
	argNbToPing := ps.ByName("times")
	nbToPing, err := strconv.Atoi(argNbToPing)
	if err != nil {
		json.NewEncoder(w).Encode(PingError{"KO", "You didn't give a number of pings to make"})
		return
	}
	if nbToPing > 100 {
		json.NewEncoder(w).Encode(PingError{"KO", "Unable to make more than 100 pings"})
		return
	}

	res, err := ping(addrToPing, int(nbToPing))
	if err != nil {
		log.Print("Got an error", err)
		json.NewEncoder(w).Encode(err)
		return
	}

	// Write the result in the json response
	json.NewEncoder(w).Encode(res)
}
