# Track Me

- Professional Garmin Explore Plan needed
- inReach Mini 2

## Stack

- go
- SQLite
- leaflet
  - [leaflet-gpx](https://github.com/mpetazzoni/leaflet-gpx)
  - [leaflet-elevation](https://github.com/Raruto/leaflet-elevation)

## IPC Inbound & Outbound

The Inbound services can be used to interact with the InReach device. For example messages can be sent, location can be requested etc.
The Outbound services send periodically (10 minutes) HTTP POST requests to your service.

<https://explore.garmin.com/IPCInbound/docs/>

### API

>The Garmin data push service requires end users to setup a web service to handle incoming HTTP-POST requests from the Garmin gateway.

## Development

Use [Air](https://github.com/air-verse/air) for server hot reloading.

```sh
go install github.com/air-verse/air@latest
# add go/bin/air binary as an alias
echo "alias air=$(go env GOPATH)" >> ~/.zshrc
```
