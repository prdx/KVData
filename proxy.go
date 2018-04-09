package main

import (
	status "./status"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	ips, ports []string
)

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

type ResponseItem struct {
	Key      string `json:"key"`
	Encoding string `json:"encoding"`
}

type Response struct {
	Status int
	Items  []ResponseItem `json:items`
}

// To maintain where each value is stored in each server
var addressBook = map[string]string{}

func request_handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
	contents, _ := ioutil.ReadAll(r.Body)
	switch r.URL.Path {
	case "/set":
		handle_set(r, contents)
	case "/get":
		handle_get(r, contents)
	}
}

func handle_get(r *http.Request, contents []uint8) {
	switch r.Method {
	case "POST":
		handle_get_post(r, contents)
	case "GET":
		handle_get_get(r, contents)
	}
}

func handle_set(r *http.Request, contents []uint8) {
	switch r.Method {
	case "POST":
		handle_set_post(r, contents)
	case "PUT":
		handle_set_put(r, contents)
	}
}

func handle_set_post(r *http.Request, contents []uint8) {
	d := build_kvdata_array(contents)
	_, destinations := build_destination_list(d, "POST")

	var wg sync.WaitGroup
	wg.Add(len(destinations))
	respChan := make(chan *http.Response)
	resps := make([]*http.Response, 0)

	for address, data := range destinations {
		json_obj, _ := json.Marshal(data)
		url := strings.Join([]string{"http://", address, r.URL.Path}, "")
		response, err := http.NewRequest("POST", url, bytes.NewBuffer(json_obj))
		if err != nil {
			os.Exit(2)
		} else {
			go func(response *http.Request) {
				defer response.Body.Close()
				defer wg.Done()
				response.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp_received, err := client.Do(response)
				if err != nil {
					panic(err)
				} else {
					respChan <- resp_received
				}
				time.Sleep(time.Second * 2)
			}(response)

			for _, d := range data {
				addressBook[d.Key] = address
			}
		}
	}

	go func() {
		for response := range respChan {
			resps = append(resps, response)
		}
	}()
	wg.Wait()
	fmt.Println(resps)
}

func handle_set_put(r *http.Request, contents []uint8) {
	d := build_kvdata_array(contents)
	_, destinations := build_destination_list(d, "PUT")

	var wg sync.WaitGroup
	wg.Add(len(destinations))
	respChan := make(chan *http.Response)
	resps := make([]*http.Response, 0)

	for address, data := range destinations {
		json_obj, _ := json.Marshal(data)
		url := strings.Join([]string{"http://", address, r.URL.Path}, "")
		response, err := http.NewRequest("PUT", url, bytes.NewBuffer(json_obj))
		if err != nil {
			os.Exit(2)
		} else {
			go func(response *http.Request) {
				defer response.Body.Close()
				defer wg.Done()
				response.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp_received, err := client.Do(response)
				if err != nil {
					panic(err)
				} else {
					respChan <- resp_received
				}
				time.Sleep(time.Second * 2)
			}(response)

			for _, d := range data {
				addressBook[d.Key] = address
			}
		}
	}

	go func() {
		for response := range respChan {
			resps = append(resps, response)
		}
	}()
	wg.Wait()
	fmt.Println(resps)
}

func handle_get_get(r *http.Request, contents []uint8) {
	i := 0
	var wg sync.WaitGroup
	resps := make([]*http.Response, 0)
	respChan := make(chan *http.Response)
	wg.Add(len(ips))
	for i < len(ips) {
		url := strings.Join([]string{"http://", string(ips[i]), ":", string(ports[i]), r.URL.Path}, "")
		response, err := http.NewRequest("GET", url, nil)
		if err != nil {
			os.Exit(2)
		} else {
			go func(response *http.Request) {
				defer wg.Done()
				response.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp_received, err := client.Do(response)
				if err != nil {
					panic(err)
				} else {
					respChan <- resp_received
				}
				time.Sleep(time.Second * 2)
			}(response)
		}
		i++
	}
	go func() {
		for response := range respChan {
			resps = append(resps, response)
		}
	}()
	wg.Wait()
	fmt.Println(resps)
}

func handle_get_post(r *http.Request, contents []uint8) {
	ks := build_queries_object(contents)
	json_obj, _ := json.Marshal(ks)
	i := 0
	var wg sync.WaitGroup
	resps := make([]*http.Response, 0)
	respChan := make(chan *http.Response)
	wg.Add(len(ips))
	for i < len(ips) {
		url := strings.Join([]string{"http://", string(ips[i]), ":", string(ports[i]), r.URL.Path}, "")
		response, err := http.NewRequest(r.Method, url, bytes.NewBuffer(json_obj))
		if err != nil {
			os.Exit(2)
		} else {
			go func(response *http.Request) {
				defer wg.Done()
				response.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp_received, err := client.Do(response)
				if err != nil {
					panic(err)
				} else {
					respChan <- resp_received
				}
				time.Sleep(time.Second * 2)
			}(response)
		}
		i++
	}
	go func() {
		for response := range respChan {
			resps = append(resps, response)
		}
	}()
	wg.Wait()
	fmt.Println(resps)
}

func build_queries_object(contents []uint8) Queries {
	ks := Queries{}
	err := json.Unmarshal(contents, &ks)

	if err != nil {
		fmt.Println("Error when extracting json")
		fmt.Println(err)
		os.Exit(1)
	}
	return ks
}

func build_kvdata_array(contents []uint8) []KVData {
	var d []KVData
	err := json.Unmarshal(contents, &d)
	if err != nil {
		fmt.Println("Error when extracting json")
		fmt.Println(err)
		os.Exit(1)
	}
	return d
}

func build_destination_list(d []KVData, mode string) (int, map[string][]KVData) {
	// Get number of servers
	n_server := len(ips)
	// Prepare the seed
	rand.Seed(time.Now().Unix())
	// Init the destinations
	destinations := make(map[string][]KVData)
	status_code := status.SUCCESS

	var address string

	for _, el := range d {
		temp := KVData{
			Key: el.Key,
			Value: Value{
				Encoding: el.Value.Encoding,
				Data:     el.Value.Data,
			},
		}

		if mode == "POST" {
			// Assign data to server randomly
			// Random is one of the method for load balancing
			idx := rand.Int() % n_server
			address = ips[idx] + ":" + ports[idx]
		} else {
			// Check if it has the key already, if doesn't update code to 206
			if key_exists(el.Key) {
				address = addressBook[el.Key]
			} else {
				status_code = status.PARTIAL_SUCCESS
				continue
			}
		}
		destinations[address] = append(destinations[address], temp)
	}
	return status_code, destinations
}

func key_exists(key string) bool {
	if _, ok := addressBook[key]; ok {
		return true
	}
	return false
}

//func format_response(response []*http.Response) (byte[], int) {
//}

func build_addresses(servers []string) ([]string, []string) {
	ips := make([]string, len(servers))
	ports := make([]string, len(servers))

	for i := 0; i < len(servers); i++ {
		temp := strings.Split(servers[i], ":")
		ips[i] = temp[0]
		ports[i] = temp[1]
	}

	return ips, ports
}

func main() {
	arg := os.Args[0:]
	servers := arg[1:]

	ips, ports = build_addresses(servers)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		request_handler(w, r)
	})

	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":5595", nil))
}
