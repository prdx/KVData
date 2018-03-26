package main

import (
	"encoding/json"
	"fmt"
)

type Key struct {
	Encoding string `json:"encoding"`
	Data     string `json:"data"`
}

type Value struct {
	Encoding string `json:"encoding"`
	Data     string `json:"data"`
}

type Data struct {
	Key   `json:"key"`
	Value `json:"value"`
}
