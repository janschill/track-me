package db

import (
	"math"
	"time"
)

func perpendicularDistance(event, lineStart, lineEnd Event) float64 {
	if lineStart.Latitude == lineEnd.Latitude && lineStart.Longitude == lineEnd.Longitude {
		return math.Sqrt(math.Pow(event.Latitude-lineStart.Latitude, 2) + math.Pow(event.Longitude-lineStart.Longitude, 2))
	}

	numerator := math.Abs((lineEnd.Longitude-lineStart.Longitude)*event.Latitude - (lineEnd.Latitude-lineStart.Latitude)*event.Longitude + lineEnd.Latitude*lineStart.Longitude - lineEnd.Longitude*lineStart.Latitude)
	denominator := math.Sqrt(math.Pow(lineEnd.Longitude-lineStart.Longitude, 2) + math.Pow(lineEnd.Latitude-lineStart.Latitude, 2))
	return numerator / denominator
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

const earthRadiusKm float64 = 6371

// Calculate distance between two coordinates
func haversine(lat1, lon1, lat2, lon2 float64) (km float64) {
	lat1 = degreesToRadians(lat1)
	lon1 = degreesToRadians(lon1)
	lat2 = degreesToRadians(lat2)
	lon2 = degreesToRadians(lon2)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	km = c * earthRadiusKm

	return km
}

// Ramer–Douglas–Peucker algorithm
// Used to decimate a curve composed of line segments to a similar curve with fewer points
func Rdp(events []Event, epsilon float64) []Event {
	if len(events) < 3 {
		return events
	}

	dmax := 0.0
	index := 0
	end := len(events) - 1

	for i := 1; i < end; i++ {
		d := perpendicularDistance(events[i], events[0], events[end])
		if dmax < d {
			index = i
			dmax = d
		}
	}

	if dmax > epsilon {
		recResults1 := Rdp(events[:index+1], epsilon)
		recResults2 := Rdp(events[index:], epsilon)

		return append(recResults1[:len(recResults1)-1], recResults2...)
	}

	return []Event{events[0], events[end]}
}

// m/s
// to get km/h: * 3.6
func calculateSpeed(event1, event2 Event) float64 {
	distance := haversine(event1.Latitude, event1.Longitude, event2.Latitude, event2.Longitude)
	timeDiff := time.Unix(event2.TimeStamp, 0).Sub(time.Unix(event1.TimeStamp, 0)).Seconds()
	if timeDiff == 0 {
		return 0
	}
	return distance / timeDiff
}

// movingTime in Seconds
// averageSpeed in km/h
func CalculateMovingTimeAndAverageSpeed(events []Event, speedThreshold float64) (float64, float64) {
	var totalKms float64
	var totalSeconds float64

	for i := 1; i < len(events); i++ {
		speed := calculateSpeed(events[i-1], events[i])
		if speed > speedThreshold {
			totalKms += haversine(events[i-1].Latitude, events[i-1].Longitude, events[i].Latitude, events[i].Longitude)
			totalSeconds += time.Unix(events[i].TimeStamp, 0).Sub(time.Unix(events[i-1].TimeStamp, 0)).Seconds()
		}
	}

	if totalSeconds == 0 {
		return 0, 0
	}

	averageSpeed := totalKms / totalSeconds
	return totalSeconds, averageSpeed * 3.6 * 1000 // km/h
}

func CalculateDistance(events []Event) float64 {
	if len(events) < 2 {
		return 0
	}

	var kms float64
	for i := 1; i < len(events); i++ {
		kms += haversine(events[i-1].Latitude, events[i-1].Longitude, events[i].Latitude, events[i].Longitude)
	}

	return kms * 1000 // return in meters
}

func CalculateAltitudes(events []Event) (averageAltitude, maxAltitude, minAltitude int64) {
	if len(events) < 1 {
		return 0, 0, 0
	}
	var totalAltitude int64
	minAltitude = math.MaxInt64
	for _, event := range events {
		if event.Altitude > maxAltitude {
			maxAltitude = event.Altitude
		}
		if event.Altitude < minAltitude {
			minAltitude = event.Altitude
		}
		totalAltitude += event.Altitude
	}
	averageAltitude = totalAltitude / int64(len(events))
	return averageAltitude, maxAltitude, minAltitude
}

func CalculateElevationGainAndLoss(events []Event) (elevationGain, elevationLoss int64) {
	if len(events) < 2 {
		return 0, 0
	}

	for i := 1; i < len(events); i++ {
		altitudeDiff := events[i].Altitude - events[i-1].Altitude
		if altitudeDiff > 0 {
			elevationGain += altitudeDiff
		} else {
			elevationLoss -= altitudeDiff
		}
	}
	return elevationGain, elevationLoss
}

// TODO:
func CalculateStops(events []Event) (numberOfStops, totalStopTimeInSeconds int) {
	if len(events) < 2 {
		return 0, 0
	}

	inStop := false
	for i := 1; i < len(events); i++ {
		distance := haversine(events[i-1].Latitude, events[i-1].Longitude, events[i].Latitude, events[i].Longitude)
		timeDiff := time.Unix(events[i].TimeStamp, 0).Sub(time.Unix(events[i-1].TimeStamp, 0)).Seconds()
		speed := (distance / timeDiff) * 3.6 // Convert m/s to km/h

		if speed < 1.0 { // Assuming speed < 1 km/h is considered a stop
			if !inStop {
				numberOfStops++
				inStop = true
			}
			totalStopTimeInSeconds += int(timeDiff)
		} else {
			inStop = false
		}
	}
	return numberOfStops, totalStopTimeInSeconds
}

// TODO:
func CalculateMaxSpeed(events []Event) (maxSpeed float64) {
	for i := 1; i < len(events); i++ {
		distance := haversine(events[i-1].Latitude, events[i-1].Longitude, events[i].Latitude, events[i].Longitude)
		timeDiff := time.Unix(events[i].TimeStamp, 0).Sub(time.Unix(events[i-1].TimeStamp, 0)).Seconds()
		speed := (distance / timeDiff) * 3.6 * 1000 // Convert m/s to km/h

		if speed > maxSpeed {
			maxSpeed = speed
		}
	}
	return maxSpeed
}
