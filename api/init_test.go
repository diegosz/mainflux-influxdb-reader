/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package api_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mainflux/mainflux-influxdb-reader/api"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	// Start the HTTP server
	ts = httptest.NewServer(api.HTTPServer())
	defer ts.Close()

	code := m.Run()

	// Exit tests
	os.Exit(code)
}
