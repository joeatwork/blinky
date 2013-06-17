package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Sample struct {
	date time.Time
	datum float64 // Might be NaN
}

type sampleSlice []Sample

func (self sampleSlice) Len() int {
	return len(self)
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

func poll(event, where string, client *ReportClient) (float64, error) {
	now := time.Now()
	lastHourAdjustment := 60.0 / float64(now.Minute()) // Typically > 1
	yesterday := now.Add(time.Duration(-time.Hour * 48))
	todayString := now.Format("2006-01-02")
	yesterdayString := yesterday.Format("2006-01-02")

	params := map[string] string {
		"event": event,
		"type": "general",
		"from_date": yesterdayString,
		"to_date": todayString,
		"unit": "hour",
	}
	if where != "" {
		params["where"] = where
	}

	response_i, err := client.Request("segmentation", params)
	if err != nil {
		panic("Couldn't understand the segmentation report")
	}
	response := response_i.(map[string] interface{})
	data_i, ok := response["data"]
	if ! ok {
		panic("Can't find {data:}")
	}
	data := data_i.(map[string] interface{})
	values_i, ok := data["values"]
	if ! ok {
		panic("Can't find {data:{values:}}")
	}
	values := values_i.(map[string] interface{})
	segment_i := values[event]
	segment := segment_i.(map[string] interface{})
	samples := make([]Sample, len(segment))
	i := 0
	for datestring, datum_i := range segment {
		date, err := time.Parse("2006-01-02 15:04:05", datestring)
		if err != nil {
			fmt.Printf(
				"Can't parse %s as date: %s\n",
				datestring,
				err.Error(),
			)
			panic("Can't parse segment key as date")
		}
		var datum float64
		if datum_i == nil {
			datum = math.NaN()
		} else {
			datum = datum_i.(float64)
		}
		samples[i] = Sample{
			date: date,
			datum: datum,
		}
		i++
	}

	sort.Sort(sampleSlice(samples))
	var maxSampleVal float64 = 0.0
	var minSampleVal float64 = math.Inf(1)
	var currentSample float64 = 0.0
	for _, sample := range samples {
		if ! math.IsNaN(sample.datum) {
			if sample.datum > 0.0 &&
				currentSample > 0.0 &&
				minSampleVal > currentSample {
				// One behind current Sample
				minSampleVal = sample.datum
			}
			if sample.datum > 0.0 {
				currentSample = sample.datum
			}
			if maxSampleVal < sample.datum {
				maxSampleVal = sample.datum
			}
		}
	}

	scale := 0.0
	if maxSampleVal > 0.0 {
		scale = lastHourAdjustment / maxSampleVal
	}

	scaledSample := 0.0
	if currentSample > 0.0 {
		fmt.Printf("Current %v Last Hour %v Maxval %v Minval %v Scale %v\n",
			currentSample, lastHourAdjustment,
			maxSampleVal, minSampleVal, scale)
		fmt.Printf("    Minval %v * %v == %v\n",
			minSampleVal, 1 / maxSampleVal, minSampleVal / maxSampleVal)
		scaledSample = currentSample * scale
	}
	if scaledSample > 1.0 {
		scaledSample = 1.0
	}

	return scaledSample, nil
}

func RunPollingService(pollingRate time.Duration, event, where string, client *ReportClient) (value chan float64, kill chan bool) {
	// Closes the channel on error
	value = make(chan float64)
	kill = make(chan bool, 1)
	die := false
	var err *error = nil
	go func () {
		time.Sleep(time.Second * time.Duration(rand.Intn(60)))
		for ! die && err == nil {
			val, err := poll(event, where, client)
			if err != nil {
				fmt.Printf("ERROR IN POLLER")
				fmt.Printf("  %v, %v, %v", event, where, client)
				fmt.Printf(err.Error())
			} else {
				value <- val
				time.Sleep(pollingRate)
				select {
				case die = <- kill:
				default:
				}
			}
		}
		fmt.Printf("Poller exiting")
		close(value)
	}()

	return
}
