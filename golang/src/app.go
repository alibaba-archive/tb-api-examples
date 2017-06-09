package main

import (
	"fmt"
	"net/http"
	"tb-api-examples/golang/src/api"
)

func main() {

	http.HandleFunc("/auth", api.HandleMain)
	http.HandleFunc("/TBLogin", api.HandleTBLogin)
	http.HandleFunc("/tb/callback", api.HandleTBCallback)
	fmt.Println(http.ListenAndServe(":3000", nil))
}
