package db

type Table struct {
	Name       string
	Definition string
}

type Schema struct {
	Tables []Table
}

var schema = Schema{
	Tables: []Table{
		{
			Name: "trips",
			Definition: `CREATE TABLE IF NOT EXISTS trips (
				"id" INTEGER PRIMARY KEY AUTOINCREMENT,
				"startTime" DATETIME,
				"endTime" DATETIME,
				"description" TEXT
			);`,
		},
		{
			Name: "events",
			Definition: `CREATE TABLE IF NOT EXISTS events (
			 "id" INTEGER PRIMARY KEY AUTOINCREMENT,
				"tripId" INTEGER NOT NULL,
				"imei" TEXT NOT NULL,
				"messageCode" INTEGER NOT NULL,
				"freeText" TEXT,
				"timeStamp" INTEGER NOT NULL,
				"latitude" REAL,
				"longitude" REAL,
				"altitude" INTEGER,
				"gpsFix" INTEGER,
				"course" REAL,
				"speed" REAL,
				"autonomous" INTEGER,
				"lowBattery" INTEGER,
				"intervalChange" INTEGER,
				"resetDetected" INTEGER,
				FOREIGN KEY(tripId) REFERENCES trips(id)
			);`,
		},
		{
			Name: "addresses",
			Definition: `CREATE TABLE IF NOT EXISTS addresses (
				"id" INTEGER PRIMARY KEY AUTOINCREMENT,
				"eventId" INTEGER NOT NULL,
				"address" TEXT NOT NULL,
				FOREIGN KEY (eventId) REFERENCES events(id)
			);`,
		},
		{
			Name: "events_cache",
			Definition: `CREATE TABLE IF NOT EXISTS events_cache (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			points TEXT NOT NULL,
			tripId INTEGER NOT NULL,
			averageSpeed REAL,
			maxSpeed REAL,
			minSpeed REAL,
			totalDistanceInMeters INTEGER,
			elevationGain INTEGER,
			elevationLoss INTEGER,
			averageAltitude REAL,
			maxAltitude INTEGER,
			minAltitude INTEGER,
			movingTimeInSeconds INTEGER,
			numberOfStops INTEGER,
			totalStopTimeInSeconds INTEGER,
			timeStamp INTEGER NOT NULL
			);`,
		},
		{
			Name: "messages",
			Definition: `CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tripId INTEGER NOT NULL,
			timeStamp INTEGER NOT NULL,
			name TEXT,
			message TEXT,
			sentToGarmin INTEGER,
			FOREIGN KEY(tripId) REFERENCES trips(id)
			);`,
		},
	},
}
