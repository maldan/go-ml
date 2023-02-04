package ms_core

import "net/http"

func DisableCors(rw http.ResponseWriter) {
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Header().Add("Access-Control-Allow-Methods", "*")
	rw.Header().Add("Access-Control-allow-Headers", "*")
}
