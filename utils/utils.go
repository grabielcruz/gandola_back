package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	rand.Seed(time.Now().UTC().Unix())
	for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func SendInternalServerError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Error del servidor")
	log.Fatal(err)
}