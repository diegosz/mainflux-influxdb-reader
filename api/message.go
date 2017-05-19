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
	var st string
	var et string
	var err error
	st = r.URL.Query().Get("start_time")
	if len(st) == 0 {
		// Invoking an empty time.Time struct literal will return Go's zero date.
		// fmt.Println(time.Time{}) -> 0001-01-01 00:00:00 +0000 UTC
		// st = time.Time{}.Format(time.RFC3339)
		// However, InfluxDB currently does not support times before Unix epoch.
		// The Unix epoch is the time 00:00:00 UTC on 1 January 1970
		st = "1970-01-01T00:00:00Z"
	}

	et = r.URL.Query().Get("end_time")
	if len(et) == 0 {
		et = time.Now().Format(time.RFC3339)
	}

	//TODO check cid
	q := fmt.Sprintf("SELECT * FROM \"%s\" WHERE time >= '%s' AND time <= '%s'", cid, st, et)

	results, e := db.InfluxQueryDB(q)
	if e != nil {
		println(e.Error())
		return
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
			case "created":
				m.Created = x.(string)
			case "channel":
				m.Channel = x.(string)
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
