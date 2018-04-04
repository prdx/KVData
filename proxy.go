package main

import (
  "encoding/json"
  "fmt"
  "os"
  "io/ioutil"
  "math/rand"
  "net/http"
  "log"
  _ "./status"
  "strings"
  "time"
)

var (
  ips, ports []string
)

type Value struct {
  Encoding string `json:"encoding"`
  Data string `json:"data"`
}

type KVData struct {
  Key string `json:"key"`
  Value `json:"value"`
}

// To maintain where each value is stored in each server
var addressBook map[string]string

func request_handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
  contents, _ := ioutil.ReadAll(r.Body)
  switch r.URL.Path {
  case "/set":
    handle_set(r, contents)
  case "/get":
    fmt.Println("/get")
  }
}

func handle_set(r *http.Request, contents []uint8) {
  switch r.Method {
  case "POST":
    json_to_object(contents)
  }
}

func json_to_object(contents []uint8) (map[int][]KVData) {
  var d []KVData
  err := json.Unmarshal(contents, &d)
  if err != nil {
    fmt.Println("Error when extracting json")
    fmt.Println(err)
    os.Exit(1)
  }

  // Get number of servers
  n_server := len(ips)
  fmt.Println(n_server)

  // Prepare the seed
  rand.Seed(time.Now().Unix())

  // Init the destinations
  destinations := make(map[int][]KVData)

  for _, el := range d {
    temp := KVData{
      Key: el.Key,
      Value: Value{
        Encoding: el.Value.Encoding,
        Data: el.Value.Data,
      },
    }
    // Assign data to server randomly
    // Random is one of the method for load balancing
    idx := rand.Int() % n_server
    destinations[idx] = append(destinations[idx], temp)
  }
  return destinations
}

func is_duplicate(key string) bool {
  if _, ok := addressBook[key]; ok {
    return true
  }
  return false
}

func handle_get(r *http.Request) {}

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
