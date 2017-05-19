/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package api_test

import (
	"net/http"
	"testing"
	"time"
	"fmt"
	"io/ioutil"
	"encoding/json"

	ic "github.com/influxdata/influxdb/client/v2"
	mfdb "github.com/mainflux/mainflux-influxdb-reader/db"
	"github.com/mainflux/mainflux-influxdb-reader/models"
)

type (
	NatsMsg struct {
		Channel   string `json:"channel"`
		Publisher string `json:"publisher"`
		Protocol  string `json:"protocol"`
		Payload   []byte `json:"payload"`
	}
)

func TestGetMessage(t *testing.T) {
	cases := []struct {
		id string
		code int
	}{
		{"testID", 200},
		// TODO: unvalidID
	}

	tt := time.Now().UTC().Format(time.RFC3339)

	// New InfluxDB point batch
	bp, err := ic.NewBatchPoints(ic.BatchPointsConfig{
		Database:  "mainflux",
		Precision: "s",
	})

	// InfluxDB tags
	tags := map[string]string{"name": "testName"}

	// InfluxDB fields
	fields := make(map[string]interface{})
	fields["unit"] = "v"
	fields["value"] = 77.8
	fields["channel"] = "testID"
	fields["publisher"] = "s"
	fields["protocol"] = "http"
	fields["created"] = tt

	pt, err := ic.NewPoint("testID", tags, fields, time.Now())
	if err != nil {
		t.Errorf("influxDB NewPoint error")
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := mfdb.InfluxClient.Write(bp); err != nil {
		t.Errorf(err.Error())
		return
	}

	url := fmt.Sprintf("%s/msg/testID", ts.URL)

	for i, c := range cases {
		res, err := http.Get(url)

		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
			return
		}
		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d, got %d", i+1, c.code, res.StatusCode)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
			return
		}
		var msgs []models.InfluxMessage
		json.Unmarshal([]byte(body), &msgs)
		for _, v := range msgs {
			if v.Channel != c.id {
				t.Errorf("case %d: %s", i+1, "unvalid ID")
			}
	   }
	}
}
