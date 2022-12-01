package main

import (
	"fmt"
	"net/http"

	"github.com/justintout/protohackers"
)

func main() {
	fmt.Printf("started listener 0 at: %s\n", protohackers.Addr0)
	http.ListenAndServe(protohackers.Addr0, protohackers.Mux0())
}
