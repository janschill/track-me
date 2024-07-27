package db

import (
	"math"
)

func Speed() {}

func MovingTimeInSeconds() {}

func perpendicularDistance(event, lineStart, lineEnd Event) float64 {
	if lineStart.Latitude == lineEnd.Latitude && lineStart.Longitude == lineEnd.Longitude {
		return math.Sqrt(math.Pow(event.Latitude-lineStart.Latitude, 2) + math.Pow(event.Longitude-lineStart.Longitude, 2))
	}

	numerator := math.Abs((lineEnd.Longitude-lineStart.Longitude)*event.Latitude - (lineEnd.Latitude-lineStart.Latitude)*event.Longitude + lineEnd.Latitude*lineStart.Longitude - lineEnd.Longitude*lineStart.Latitude)
	denominator := math.Sqrt(math.Pow(lineEnd.Longitude-lineStart.Longitude, 2) + math.Pow(lineEnd.Latitude-lineStart.Latitude, 2))
	return numerator / denominator
}

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

func CalculateSpeed(events []Event) int64 {
	return 0
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

const earthRadiusKm float64 = 6371

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

// Calculate the total distance for a slice of events
func CalculateDistance(events []Event) float64 {
	if len(events) < 2 {
		return 0
	}

	var totalDistance float64
	for i := 1; i < len(events); i++ {
		totalDistance += haversine(events[i-1].Latitude, events[i-1].Longitude, events[i].Latitude, events[i].Longitude)
	}

	return totalDistance
}
