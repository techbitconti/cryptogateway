package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"net"
	"network"
)

func ReSponseMessage(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {

	defer wg.Done()

	api := r.URL.Path[1:]

	if _, ok := network.PROCCESSING_MAP[api]; ok {
		params, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		network.PROCCESSING_MAP[api](r.RemoteAddr, w, params)

	} else {
		w.Write([]byte("Error 404"))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	ip, _ := FromRequest(r)

	fmt.Println("ip : ", ip)

	//	if _, ok := config.Ips[ip]; !ok {
	//		return
	//	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method != "POST" {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	payload := Payload{Conn: w, Message: r, Wait: &wg}
	work := Job{Payload: payload}
	fmt.Println("JobQueu <- work")
	JobQueu <- work
	fmt.Println("work done: ")

	wg.Wait()
}

func FromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

func Create(host string) {

	dispatcher := NewDispatcher()
	dispatcher.Run()
	fmt.Println("Worker run ....................")

	http.HandleFunc("/", handler)
	http.ListenAndServe(host, nil)

}
