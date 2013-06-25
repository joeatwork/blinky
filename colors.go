package main

import (
	"math"
	"time"
)

func colorForCurrentSample(samples []Sample) float64 {
	var maxSampleVal float64 = 0.0
	var minSampleVal float64 = math.Inf(1)
	var previousSample float64 = 0.0
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
				previousSample = currentSample
				currentSample = sample.datum
			}
			if maxSampleVal < sample.datum {
				maxSampleVal = sample.datum
			}
		}
	}

	minutes := float64(time.Now().Minute())
	var scaledSample float64
	if (minutes < 2) {
		// For the first minute, use the last value
		// instead of introducing a lot of artificial
		// detail.
		scaledSample = previousSample
	} else {
		lastHourAdjustment := 60.0 / minutes // Typically > 1
		scale := 0.0
		if maxSampleVal > 0.0 {
			scale = lastHourAdjustment / maxSampleVal
		}

		scaledSample = 0.0
		if currentSample > 0.0 {
			scaledSample = currentSample * scale
		}
	}

	if scaledSample > 1.0 {
		scaledSample = 1.0
	}

	return scaledSample
}
