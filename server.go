package main

import (
  "encoding/json"
  "fmt"
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
    handle_set(r, contents)
  case "/get":
    fmt.Println("/get")
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

func handle_set(r *http.Request, contents []uint8) {
  switch r.Method {
  case "POST":
    d := json_to_object_post(contents)
    save(d)
    fmt.Println(dataStore)
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
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    request_handler(w, r)
  })

  fmt.Println("Starting server...")
  log.Fatal(http.ListenAndServe(":8181", nil))
}
