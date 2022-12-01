package protohackers

import (
	"io"
	"net/http"
)

const Addr0 = ":10000"

func Mux0() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	return mux
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if _, err := io.Copy(w, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
