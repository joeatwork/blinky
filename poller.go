package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Sample struct {
	date  time.Time
	datum float64 // Might be NaN
}

type PollError string

type sampleSlice []Sample

func (self sampleSlice) Len() int {
	return len(self)
}

func (self PollError) Error() string {
	return string(self)
}

func (self sampleSlice) Less(i, j int) bool {
	i_date := self[i].date
	j_date := self[j].date
	return i_date.Before(j_date)
}

func (self sampleSlice) Swap(i, j int) {
	tmp := self[i]
	self[i] = self[j]
	self[j] = tmp
}

func poll(event, where string, client *ReportClient) (samples []Sample, err error) {
	err = nil
	now := time.Now()
	ago := now.Add(time.Duration(-time.Hour * 24))
	todayString := now.Format("2006-01-02")
	agoString := ago.Format("2006-01-02")

	params := map[string]string{
		"event":     event,
		"type":      "general",
		"from_date": agoString,
		"to_date":   todayString,
		"unit":      "hour",
	}
	if where != "" {
		params["where"] = where
	}

	response_i, err := client.Request("segmentation", params)
	if err != nil {
		return
	}
	response := response_i.(map[string]interface{})
	data_i, ok := response["data"]
	if !ok {
		err = PollError("Can't find {data:}")
		return
	}
	data := data_i.(map[string]interface{})
	values_i, ok := data["values"]
	if !ok {
		err = PollError("Can't find {data:{values:}}")
		return
	}
	values := values_i.(map[string]interface{})
	segment_i := values[event]
	segment := segment_i.(map[string]interface{})
	samples = make([]Sample, len(segment))
	i := 0
	for datestring, datum_i := range segment {
		date, terr := time.Parse("2006-01-02 15:04:05", datestring)
		if terr != nil {
			fmt.Printf(
				"Can't parse %s as date: %s\n",
				datestring,
				terr.Error(),
			)
			err = terr
			return
		}
		var datum float64
		if datum_i == nil {
			datum = math.NaN()
		} else {
			datum = datum_i.(float64)
		}
		samples[i] = Sample{
			date:  date,
			datum: datum,
		}
		i++
	}

	sort.Sort(sampleSlice(samples))
	return
}
