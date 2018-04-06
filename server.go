package main

import (
  "encoding/json"
  "fmt"
  "flag"
  "io/ioutil"
  "net/http"
  "os"
  "log"
  _ "./status"
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
  Data string `json:"data"`
}

type KVData struct {
  Key string `json:"key"`
  Value `json:"value"`
}

var dataStore = map[string]Value{}

func handle_set(w http.ResponseWriter, r *http.Request, contents []uint8) {
  switch r.Method {
  case "POST":
    d := json_to_object_post(contents)
    save(d)
    fmt.Println(dataStore)
  }
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
  }
}

func json_to_object_post(contents []uint8) ([]KVData) {
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
  log.Fatal(http.ListenAndServe(*ip + ":" + *port, nil))
}
