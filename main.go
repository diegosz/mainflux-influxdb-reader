/**
 * Copyright (c) 2017 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/mainflux/mainflux-influxdb-reader/api"
	"github.com/mainflux/mainflux-influxdb-reader/db"
)

const (
	help string = `
Usage: mainflux-influxdb [options]
Options:
	-a, --addr	HTTP host
	-r, --port	HTTP port
	-i, --ihost	InfluxDB host
	-q, --iport	InfluxDB port
	-s, --db	InfluxDB database
	-u, --user	InfluxDB username
	-p, --pass	InfluxDB password
	-p, --precision	InfluxDB time precision
	-h, --help	Show help
`
)

type (
	Opts struct {
		HTTPHost string
		HTTPPort string

		InfluxHost      string
		InfluxPort      string
		InfluxDatabase  string
		InfluxUser      string
		InfluxPass      string
		InfluxPrecision string

		Help bool
	}
)

func main() {
	opts := Opts{}

	flag.StringVar(&opts.HTTPHost, "a", "localhost", "HTTP host.")
	flag.StringVar(&opts.HTTPPort, "r", "7080", " HTTP port.")
	flag.StringVar(&opts.InfluxHost, "i", "localhost", "InfluxDB host.")
	flag.StringVar(&opts.InfluxPort, "q", "8086", "InfluxDB port.")
	flag.StringVar(&opts.InfluxDatabase, "d", "mainflux", "InfluxDB databse name.")
	flag.StringVar(&opts.InfluxUser, "u", "mainflux", "InfluxDB username.")
	flag.StringVar(&opts.InfluxPass, "s", "", "InfluxDB password.")
	flag.StringVar(&opts.InfluxPrecision, "p", "s", "InfluxDB time precision.")
	flag.BoolVar(&opts.Help, "h", false, "Show help.")
	flag.BoolVar(&opts.Help, "help", false, "Show help.")

	flag.Parse()

	if opts.Help {
		fmt.Printf("%s\n", help)
		os.Exit(0)
	}

	// Connect to InfluxDB
	if err := db.InfluxInit(opts.InfluxHost, opts.InfluxPort, opts.InfluxDatabase,
		opts.InfluxUser, opts.InfluxPass, opts.InfluxPrecision); err != nil {

		log.Fatalf("InfluxDB: Can't connect: %v\n", err)
	}

	if _, err := db.InfluxQueryDB(fmt.Sprintf("CREATE DATABASE %s", opts.InfluxDatabase)); err != nil {
		log.Fatal(err)
	}

	// Print banner
	color.Cyan(banner)

	// Serve HTTP
	httpHost := fmt.Sprintf("%s:%s", opts.HTTPHost, opts.HTTPPort)
	http.ListenAndServe(httpHost, api.HTTPServer())
}

var banner = `
╦┌┐┌┌─┐┬  ┬ ┬─┐ ┬╔╦╗╔╗    ┬─┐┌─┐┌─┐┌┬┐┌─┐┬─┐  
║│││├┤ │  │ │┌┴┬┘ ║║╠╩╗───├┬┘├┤ ├─┤ ││├┤ ├┬┘  
╩┘└┘└  ┴─┘└─┘┴ └─═╩╝╚═╝   ┴└─└─┘┴ ┴─┴┘└─┘┴└─  

         == Industrial IoT System ==

        Made with <3 by Mainflux Team
[w] http://mainflux.io
[t] @mainflux

`
