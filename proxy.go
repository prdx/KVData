package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "os"
  "io/ioutil"
  "math/rand"
  "net/http"
  "log"
  _ "./status"
  "strings"
  "sync"
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

type Query struct {
  Key string `json:"key"`
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
    fmt.Println(addressBook)
  case "/get":
    //handle_get(r, contents)
  }
}

func handle_get(r *http.Request, contents []uint8) {
}

func handle_set(r *http.Request, contents []uint8) {
  // Get number of servers
  //n_server := len(ips)

  switch r.Method {
  case "POST":
    destinations := json_to_object(contents)
    var wg sync.WaitGroup
    wg.Add(len(destinations))
    respChan := make(chan *http.Response)
    resps := make([]*http.Response, 0)

    for i, destination := range destinations {
      json_obj, _ := json.Marshal(destination)
      url := strings.Join([]string{"http://", ips[i], ":", ports[i], r.URL.Path}, "")
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

        for _, d := range destination {
          addressBook[d.Key] = ips[i] + ":" + ports[i]
        }
        //addressBook[destination = ips[i] + ":" + ports[i]
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
