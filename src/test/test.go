package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonStr := `{"a":1,"b":{"a":1},"c":{"b":{"a":1}}}`
	var d interface{}
	_ = json.Unmarshal([]byte(jsonStr), &d)
	//fmt.Println(d)
	for k, v := range d.(map[string]interface{}) {
		fmt.Println(k)
		//fmt.Println(k, v)
		if v,ok := v.(map[string]interface{}) ;ok {
			fmt.Println(v)
		}

	}
}

func calDepth(m map[string]interface{}){

}
