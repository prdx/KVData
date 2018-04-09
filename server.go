package main

import (
	status "./status"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func request_handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	contents, _ := ioutil.ReadAll(r.Body)
	switch r.URL.Path {
	case "/set":
		handle_set(w, r, contents)
	case "/get":
		handle_get(w, r, contents)
	}
}

type Value struct {
	Encoding string `json:"encoding"`
	Data     string `json:"data"`
}

type KVData struct {
	Key   string `json:"key"`
	Value `json:"value"`
}

type Queries struct {
	Keys []string `json:"keys"`
}

var dataStore = map[string]Value{}

func handle_set(w http.ResponseWriter, r *http.Request, contents []uint8) {
	d := json_to_object_post(contents)
	save(d)
	fmt.Println(dataStore)
}

func handle_get(w http.ResponseWriter, r *http.Request, contents []uint8) {
	switch r.Method {
	case "GET":
		js, err := json.Marshal(dataStore)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "POST":
		ks := Queries{}
		err := json.Unmarshal(contents, &ks)
		fmt.Println(ks)
		search(ks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func search(ks Queries) (int, []KVData) {
	res := []KVData{}
	code := status.SUCCESS

	for _, k := range ks.Keys {
		if val, ok := dataStore[k]; ok {
			temp := KVData{
				Key:   k,
				Value: val,
			}
			res = append(res, temp)
		} else {
			code = status.PARTIAL_SUCCESS
		}
	}
	return code, res
}

func json_to_object_post(contents []uint8) []KVData {
	var d []KVData
	err := json.Unmarshal(contents, &d)
	if err != nil {
		fmt.Println("Error when extracting json")
		fmt.Println(err)
		os.Exit(1)
	}
	return d
}

func save(d []KVData) {
	for _, el := range d {
		dataStore[el.Key] = el.Value
	}
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "IP address")
	port := flag.String("port", "9191", "Port")
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		request_handler(w, r)
	})

	fmt.Println("Starting server at " + *ip + ":" + *port)
	log.Fatal(http.ListenAndServe(*ip+":"+*port, nil))
}
