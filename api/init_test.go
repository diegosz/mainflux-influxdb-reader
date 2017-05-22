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
	"log"
	"fmt"

	"github.com/mainflux/mainflux-influxdb-reader/api"
	mfdb "github.com/mainflux/mainflux-influxdb-reader/db"
	"gopkg.in/ory-am/dockertest.v3"
)

var ts *httptest.Server

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
	    log.Fatalf("Could not connect to docker: %s", err)
	}

	options := &dockertest.RunOptions{
	    Repository: "influxdb",
	    Tag:        "latest",
	    Mounts:     []string{"/tmp/influxdb:/etc/influxdb"},
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
	    log.Fatalf("Could not start resource: %s", err)
	}


	// host, port, databse, username, password, precision
	if err := mfdb.InfluxInit("0.0.0.0", "8086", "mainflux", "mainflux",
	                          "", "s"); err != nil {
		log.Println(err)
	}

	if _, err = mfdb.InfluxQueryDB(fmt.Sprintf("CREATE DATABASE %s", "mainflux")); err != nil {
		log.Println(err)
	}

	// Start the HTTP server
	ts = httptest.NewServer(api.HTTPServer())
	defer ts.Close()

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// Exit tests
	os.Exit(code)
}
