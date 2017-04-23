/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-zoo/bone"
	"github.com/mainflux/mainflux-influxdb-reader/db"
	"github.com/mainflux/mainflux-influxdb-reader/models"
)

// getMessage function
func getMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	cid := bone.GetValue(r, "channel_id")

	// Get fileter values from parameters:
	// - start_time = messages from this moment. UNIX time format.
	// - end_time = messages to this moment. UNIX time format.
	var st float64
	var et float64
	var err error
	var s string
	s = r.URL.Query().Get("start_time")
	if len(s) == 0 {
		st = 0
	} else {
		st, err = strconv.ParseFloat(s, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			str := `{"response": "wrong start_time format"}`
			io.WriteString(w, str)
			return
		}
	}
	s = r.URL.Query().Get("end_time")
	if len(s) == 0 {
		et = float64(time.Now().Unix())
	} else {
		et, err = strconv.ParseFloat(s, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			str := `{"response": "wrong end_time format"}`
			io.WriteString(w, str)
			return
		}
	}

	println(st, et)
	q := fmt.Sprintf("SELECT * FROM \"%s\"", cid)
	results, e := db.InfluxQueryDB(q)
	if e != nil {
		log.Fatal(e)
	}

	a := []models.InfluxMessage{}

	// Put each InfluxDB return value in different SenML message
	for _, ve := range results[0].Series[0].Values {
		m := models.InfluxMessage{}

		// For each InfluxDB column we have to set SenML message field
		for ci, ce := range results[0].Series[0].Columns {
			x := ve[ci]
			if x == nil {
				continue
			}

			switch ce {
			case "time":
				m.Time = x.(string)
			case "unit":
				m.Unit = x.(string)
			case "value":
				m.Value = x.(json.Number)
			case "string_value":
				m.StringValue = x.(string)
			case "bool_value":
				m.BoolValue = x.(bool)
			case "data_value":
				m.DataValue = x.(string)
			case "sum":
				m.Sum = x.(json.Number)
			case "update_time":
				m.UpdateTime = x.(json.Number)
			case "publisher":
				m.Publisher = x.(string)
			case "protocol":
				m.Protocol = x.(string)
			case "timestamp":
				m.Timestamp = x.(string)
			}
		}

		a = append(a, m)
	}

	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(a)
	if err != nil {
		log.Print(err)
	}
	io.WriteString(w, string(res))
}
