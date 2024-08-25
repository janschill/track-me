package main

type DayEvent struct {
	ID          int64       `json:"ID"`
	TripID      int         `json:"TripID"`
	Imei        string      `json:"Imei"`
	MessageCode int         `json:"MessageCode"`
	FreeText    string      `json:"FreeText"`
	TimeStamp   int64       `json:"TimeStamp"`
	Addresses   interface{} `json:"Addresses"` // Use interface{} if the structure is unknown or varies
	Status      struct {
		Autonomous     int `json:"Autonomous"`
		LowBattery     int `json:"LowBattery"`
		IntervalChange int `json:"IntervalChange"`
		ResetDetected  int `json:"ResetDetected"`
	} `json:"Status"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	Altitude  float64 `json:"Altitude"`
	GpsFix    int     `json:"GpsFix"`
	Course    int     `json:"Course"`
	Speed     int     `json:"Speed"`
}

// var days []service.Day = []service.Day{
// 	service.Day{
// 		Points: `[{"ID":4536,"TripID":0,"Imei":"","MessageCode":0,"FreeText":"","TimeStamp":1726012577,"Addresses":null,"Status":{"Autonomous":0,"LowBattery":0,"IntervalChange":0,"ResetDetected":0},"Latitude":47.499658,"Longitude":-113.783468,"Altitude":1244,"GpsFix":0,"Course":0,"Speed":0},{"ID":4539,"TripID":0,"Imei":"","MessageCode":0,"FreeText":"","TimeStamp":1726012748,"Addresses":null,"Status":{"Autonomous":0,"LowBattery":0,"IntervalChange":0,"ResetDetected":0},"Latitude":47.500389,"Longitude":-113.780864,"Altitude":1232,"GpsFix":0,"Course":0,"Speed":0}]`,
// 	},
// 	service.Day{
// 		Points: `[{"ID":5536,"TripID":0,"Imei":"","MessageCode":0,"FreeText":"","TimeStamp":1726012577,"Addresses":null,"Status":{"Autonomous":0,"LowBattery":0,"IntervalChange":0,"ResetDetected":0},"Latitude":47.499658,"Longitude":-113.783468,"Altitude":1244,"GpsFix":0,"Course":0,"Speed":0},{"ID":5539,"TripID":0,"Imei":"","MessageCode":0,"FreeText":"","TimeStamp":1726012748,"Addresses":null,"Status":{"Autonomous":0,"LowBattery":0,"IntervalChange":0,"ResetDetected":0},"Latitude":47.500389,"Longitude":-113.780864,"Altitude":1232,"GpsFix":0,"Course":0,"Speed":0}]`,
// 	},
// }

func main() {
	// var combinedPoints []string
	// for _, day := range days {
	// 		// Remove the leading '[' and trailing ']' from each Points string
	// 		trimmedPoints := strings.Trim(day.Points, "[]")
	// 		combinedPoints = append(combinedPoints, trimmedPoints)
	// }
	//     finalJSONArrayString := "[" + strings.Join(combinedPoints, ",") + "]"

	// log.Printf("day: %v", template.JS(finalJSONArrayString))
}
