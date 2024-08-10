# Track Me

[![Better Stack Badge](https://uptime.betterstack.com/status-badges/v1/monitor/1gy4z.svg)](https://uptime.betterstack.com/?utm_source=status_badge)

- Professional Garmin Explore Plan needed
- inReach Mini 2

## Architecture

Everything is served on the `/` path. The Garmin Outbound webhook is continously pushing new events from the Garmin InReach Mini 2 to `/garmin-outbound`, which will save all incoming events to a SQLite database.
When visiting `/` all events are queried from the DB and used to plot a traveled path on a Leaflet map. The home page also shows overall Ride stats and a breakdown of days.
The Ride and Days show different stats such as distance traveled, elevation, time moving etc. these are calculated from the events.
All past days are cached in memory. The current day is always computed newly. The Ride stats will use the cached days and the events for the current day.
Messages from visitors are stored in a messages table. Events from the Garmin that have a message will be stored in events, but also parsed to the message struct and stored in the messages table. This allows to render all messages in a list. The messages from Garmin are shown as "Automated Messages".

### Stack

- go
- SQLite
- leaflet
  - [leaflet-gpx](https://github.com/mpetazzoni/leaflet-gpx)
  - [leaflet-elevation](https://github.com/Raruto/leaflet-elevation)
- BetterStack for uptime monitor
- TODO: Grafana for logs, metrics, traces

## IPC Inbound & Outbound

The Inbound services can be used to interact with the InReach device. For example messages can be sent, location can be requested etc.
The Outbound services send periodically (10 minutes) HTTP POST requests to your service.

<https://explore.garmin.com/IPCInbound/docs/>

### API

>The Garmin data push service requires end users to setup a web service to handle incoming HTTP-POST requests from the Garmin gateway.

## Deployment

1. GitHub Actions will build the binary using Docker
2. The SCP action will copy the binary and assets from web/ to a tmp directory on the server
3. The SSH action will connect with the server and do a ping test on the binary and then replace the binary and assets in the track-me directory
4. The systemctl will then restart and use the new binary and assets to start the application

## Development

Use [Air](https://github.com/air-verse/air) for server hot reloading.

```sh
go install github.com/air-verse/air@latest
# add go/bin/air binary as an alias
echo "alias air=$(go env GOPATH)" >> ~/.zshrc
```

### Debugging

```sh
go install -v github.com/go-delve/delve/cmd/dlv@latest
```

Example launch config for running a debugger in VSCode.

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "./cmd/server/main.go",
      "env": {},
      "args": [],
      "cwd": "${workspaceRoot}"
    }
  ]
}
```

## References

- [Organising Database Access in Go](https://www.alexedwards.net/blog/organising-database-access)
- [Better Stack](https://uptime.betterstack.com/team/245141/monitors/2470355)
