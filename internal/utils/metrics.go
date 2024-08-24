package utils

import (
	"math"
	"time"

	"github.com/janschill/track-me/internal/repository"
)

func Movement(events []repository.Event) (float64, bool) {
	eventCount := len(events)
	if eventCount < 2 {
		return 0, false
	}

	secondToLastEvent := events[eventCount-2]
	lastEvent := events[eventCount-1]
	speed := CalculateSpeed(secondToLastEvent, lastEvent) * 3.6 * 1000

	return OneDecimal(speed), lastEvent.Latitude != secondToLastEvent.Latitude || lastEvent.Longitude != secondToLastEvent.Longitude
}

func Progress(distanceSoFar float64) float64 {
	distanceInTotal := 3891.0
	percentage := (distanceSoFar / 1000) / distanceInTotal * 100

	return OneDecimal(percentage)
}

func RestingTime(elapsedDays int, movingTime int64) int64 {
	return int64(elapsedDays)*24*60*60 - movingTime
}

func perpendicularDistance(event, lineStart, lineEnd repository.Event) float64 {
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
func Rdp(events []repository.Event, epsilon float64) []repository.Event {
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

	return []repository.Event{events[0], events[end]}
}

// m/s
// to get km/h: * 3.6
func CalculateSpeed(event1, event2 repository.Event) float64 {
	distance := haversine(event1.Latitude, event1.Longitude, event2.Latitude, event2.Longitude)
	timeDiff := time.Unix(event2.TimeStamp, 0).Sub(time.Unix(event1.TimeStamp, 0)).Seconds()
	if timeDiff == 0 {
		return 0
	}
	return distance / timeDiff
}

// movingTime in Seconds
// averageSpeed in km/h
func CalculateMovingTimeAndAverageSpeed(events []repository.Event, speedThreshold float64) (float64, float64, float64) {
	var totalKms float64
	var totalSeconds float64
	maxSpeed := math.Inf(-1)

	for i := 1; i < len(events); i++ {
		speed := CalculateSpeed(events[i-1], events[i])
		if speed > speedThreshold {
			totalKms += haversine(events[i-1].Latitude, events[i-1].Longitude, events[i].Latitude, events[i].Longitude)
			totalSeconds += time.Unix(events[i].TimeStamp, 0).Sub(time.Unix(events[i-1].TimeStamp, 0)).Seconds()
		}
		if speed > maxSpeed {
			maxSpeed = speed
		}
	}

	if totalSeconds == 0 {
		return 0, 0, 0
	}

	averageSpeed := totalKms / totalSeconds
	return totalSeconds, averageSpeed * 3.6 * 1000, maxSpeed // km/h
}

func DistanceInMeters(events []repository.Event) float64 {
	if len(events) < 2 {
		return 0
	}

	var kms float64
	for i := 1; i < len(events); i++ {
		kms += haversine(events[i-1].Latitude, events[i-1].Longitude, events[i].Latitude, events[i].Longitude)
	}

	return kms * 1000 // return in meters
}

func CalculateAltitudes(events []repository.Event) (averageAltitude, maxAltitude, minAltitude float64) {
	ignoredEventsCount := 0

	if len(events) < 1 {
		return 0, 0, 0
	}
	var totalAltitude float64
	minAltitude = math.MaxInt64
	for _, event := range events {
		// Message Code 10 announces a tracking start
		// Usually the altitude on these events are quite off
		if event.Altitude < 0 || event.MessageCode == 10 {
			ignoredEventsCount++
			continue
		}
		if event.Altitude > maxAltitude {
			maxAltitude = event.Altitude
		}
		if event.Altitude < minAltitude {
			minAltitude = event.Altitude
		}
		totalAltitude += event.Altitude
	}
	averageAltitude = totalAltitude / float64(len(events)-ignoredEventsCount)
	return averageAltitude, maxAltitude, minAltitude
}

func CalculateElevationGainAndLoss(events []repository.Event) (elevationGain, elevationLoss int64) {
	if len(events) < 2 {
		return 0, 0
	}

	i := 0
	for j := 1; j < len(events); j++ {
		// Message Code 10 announces a tracking start
		// Usually the altitude on these events are quite off
		if events[j].Altitude < 0 || events[j].MessageCode == 10 {
			continue
		}

		if events[i].Altitude >= 0 && events[i].MessageCode != 10 {
			altitudeDiff := events[j].Altitude - events[i].Altitude
			if altitudeDiff > 0 {
				elevationGain += int64(altitudeDiff)
			} else {
				elevationLoss -= int64(altitudeDiff)
			}
		}

		i = j
	}

	return elevationGain, elevationLoss
}
