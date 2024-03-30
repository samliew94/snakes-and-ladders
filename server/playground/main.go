package main

import (
	"encoding/json"
	"log"
)

func main() {

	mockGen(true)

}

func mockGen(b bool) {

	res := map[string]interface{}{"isAuthenticate": b}
	jbytes, _ := json.Marshal(res)
	log.Println(string(jbytes))
	
}

