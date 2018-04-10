#!/bin/bash
curl -X POST -d '[{"key":"hello5","value":{"encoding":"string","data":"value4"}},{"key":"hello13","value":{"encoding":"binary","data":"010101"}}, {"key":"hello124","value":{"encoding":"string","data":"value4"}}]' -i http://localhost:5595/set
curl -X PUT -d '[{"key":"hello13","value":{"encoding":"string","data":"value5"}},{"key":"hello9","value":{"encoding":"binary","data":"010101"}}]' -i http://localhost:5595/set 
curl -X GET -i  http://localhost:5595/get
curl -X POST -d '{"keys":["hello12"]}' -i http://localhost:5595/get
curl -X POST -d '{"keys":["hello5"]}' -i http://localhost:5595/get

