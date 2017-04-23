/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package models

import "encoding/json"

type (
	InfluxMessage struct {
		////
		// SenML stuff
		////
		Link string `json:"l,omitempty"  xml:"l,attr,omitempty"`

		Name       string      `json:"n,omitempty"  xml:"n,attr,omitempty"`
		Unit       string      `json:"u,omitempty"  xml:"u,attr,omitempty"`
		Time       string      `json:"t,omitempty"  xml:"t,attr,omitempty"`
		UpdateTime json.Number `json:"ut,omitempty"  xml:"ut,attr,omitempty"`

		Value       json.Number `json:"v,omitempty"  xml:"v,attr,omitempty"`
		StringValue string      `json:"vs,omitempty"  xml:"vs,attr,omitempty"`
		DataValue   string      `json:"vd,omitempty"  xml:"vd,attr,omitempty"`
		BoolValue   bool        `json:"vb,omitempty"  xml:"vb,attr,omitempty"`

		Sum json.Number `json:"s,omitempty"  xml:"sum,,attr,omitempty"`

		////
		// Mainflux stuff
		////
		Publisher string `json:"publisher"`
		Protocol  string `json:"protocol"`
		Timestamp string `json:"timestamp"`

		// Channel to which this message belongs
		Channel string `json:"channel"`
	}
)
