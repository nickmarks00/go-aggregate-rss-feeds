package main

import (
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{}) // we need to respond with some payload so use empty struct - only care about 200 response
}
