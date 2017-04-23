/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package api

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
)

// HTTPServer function
func HTTPServer() http.Handler {
	mux := bone.New()

	// Status
	mux.Get("/status", http.HandlerFunc(getStatus))

	// Messages
	mux.Get("/msg/:channel_id", http.HandlerFunc(getMessage))

	n := negroni.Classic()
	n.UseHandler(mux)
	return n
}
