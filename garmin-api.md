# IPC Outbound API Documentation

## Overview

- Version: Indicates the version of the schema being used.
- Events: An array containing one or more event objects, each corresponding to a transmission from an inReach device.

Each event object contains:

- IMEI: A string representing the unique 15-digit identifier of the inReach device.
- messageCode: A number indicating the type of message. This is further detailed in the Message Codes Table.
- freeText: A string containing the text message sent by the inReach device. It can be empty if no text was sent.
- timeStamp: A number representing the time the message was created, in milliseconds since the epoch (January 1, 1970).
- pingbackReceived: The time (in milliseconds since epoch) when the pingback request was received.
- pingbackResponded: The time (in milliseconds since epoch) when the pingback response was constructed.
- addresses: An array of address objects, which could include SMS phone numbers, email addresses, or synchronized contacts.
- point: An object containing the geographic information like latitude, longitude, altitude, GPS fix, course, and speed.
- status: An object indicating various statuses of the device, such as battery level, interval changes, and whether a factory reset has been detected.
- payload: A Base64 encoded string representing any binary data associated with the event​

### Point Object

The Point object contains geographic information related to the location of the inReach device at the time of the event. Here are the fields within the Point object:

| Field Name | Type | Description |
| - | - | - |
| latitude | float | Latitude of the device in decimal degrees. |
| longitude | float | Longitude of the device in decimal degrees. |
| altitude | float | Altitude of the device in meters above sea level. |
| gpsFix | integer | Indicates the quality of the GPS fix. Values typically range from 0 (no fix) to 3 (3D fix). |
| course | float | Course or direction of travel in degrees relative to true north. |
| speed | float | Speed of the device in meters per second. |

### Status Object

The Status object provides information about the current state of the inReach device. It includes details on the device’s battery, settings, and any changes in the operational mode. Below are the fields within the Status object:

| Field Name | Type | Description |
| - | - | - |
| batteryLevel | float | Current battery level as a percentage (0.0 to 100.0). |
| intervalChangeDetected | boolean | Indicates whether a change in reporting interval has been detected. |
| factoryResetDetected | boolean | Indicates whether a factory reset has been detected on the device. |
| trackingPausedDetected | boolean | Indicates whether tracking has been paused. |
| trackingStoppedDetected | boolean | Indicates whether tracking has been stopped. |

### Message Code Object

The messageCode is a key part of the event object that indicates the type of message or event that has occurred. Each code corresponds to a specific event, such as sending a message, updating the location, or triggering an SOS alert. Here's a breakdown of common message codes:

| Code | Name                | Description                                                                                                 |
|------|---------------------|-------------------------------------------------------------------------------------------------------------|
| 0    | Position Report     | Drops a breadcrumb while tracking.                                                                          |
| 1    | Reserved            | Reserved for later use.                                                                                     |
| 2    | Locate Response     | Position for a locate request.                                                                              |
| 3    | Free Text Message   | Message containing a free-text block.                                                                       |
| 4    | Declare SOS         | Declares an emergency state.                                                                                |
| 5    | Reserved            | Reserved for later use.                                                                                     |
| 6    | Confirm SOS         | Confirms an unconfirmed SOS.                                                                                |
| 7    | Cancel SOS          | Stops a SOS event.                                                                                          |
| 8    | Reference Point     | Shares a non-GPS location.                                                                                  |
| 10   | Start Track         | Begins a tracking process on the server.                                                                    |
| 11   | Track Interval      | Indicates changes in tracking interval.                                                                     |
| 12   | Stop Track          | Ends a tracking process on the server.                                                                      |
| 13   | Unknown Index       | Used when the device receives a message from the server addressed to a synced contact identifier not on the device. |
| 14   | Puck Message 1      | Sends the first of three inReach message button events.                                                     |
| 15   | Puck Message 2      | Sends the second of three inReach message button events.                                                    |
| 16   | Puck Message 3      | Sends the third of three inReach message button events.                                                     |
| 17   | Map Share           | Sends a message to the shared map.                                                                          |
| 20   | Mail Check          | Sent to determine if any messages are queued for the device.                                                |
| 21   | Am I Alive          | Sent when the device needs to determine if it is active. Automatically replied to by the Garmin server.     |
| 24-63| Pre-defined Message | The index for a text message synchronized with the server.                                                  |
| 64   | Encrypted Binary    | An encrypted binary Earthmate message.                                                                      |
| 65   | Pingback Message    | A pingback response message (initiated through IPCInbound).                                                 |
| 66   | Generic Binary      | An uninterpreted binary message.                                                                            |
| 67   | EncryptedPinpoint   | A fully-encrypted inReach message.                                                                          |
| 3099 | Canned Message      | A Quicktext message, potentially edited by the user.                                                        |

```go
type IPCOutbound struct {
  Version string  `json:"Version"`
  Events  []Event `json:"Events"`
}

// Event object containing data from an inReach device
type Event struct {
  IMEI        string     `json:"imei"`
  MessageCode int        `json:"messageCode"`
  FreeText    string     `json:"freeText"`
  TimeStamp   int64      `json:"timeStamp"`
  Addresses   []Address  `json:"addresses"`
  Point       Point      `json:"point"`
  Status      Status     `json:"status"`
}

// Address object within an Event
type Address struct {
  Address string `json:"address"`
}

// Point object within an Event
type Point struct {
  Latitude  float64 `json:"latitude"`
  Longitude float64 `json:"longitude"`
  Altitude  int `json:"altitude"`
  GPSFix    int     `json:"gpsFix"`
  Course    float64 `json:"course"`
  Speed     float64 `json:"speed"`
}

// Status object within an Event
type Status struct {
  Autonomous    int `json:"autonomous"`
  LowBattery    int `json:"lowBattery"`
  IntervalChange int `json:"intervalChange"`
  ResetDetected int `json:"resetDetected"`
}
```
