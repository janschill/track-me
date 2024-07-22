# Track Me

- Professional Garmin Explore Plan needed
- inReach Mini 2

## Stack

- go
- Sqlite
- Docker
- DigitalOcean
- _BetterStack_
- _Tailwind_

## IPC Inbound & Outbound

The Inbound services can be used to interact with the InReach device. For example messages can be sent, location can be requested etc.
The Outbound services send periodically HTTP POST requests to your service.

<https://explore.garmin.com/IPCInbound/docs/>

### API

>The Garmin data push service requires end users to setup a web service to handle incoming HTTP-POST requests from the Garmin gateway.

## Development

```sh
go install github.com/air-verse/air@latest
go env GOPATH
# update air.toml to your go/bin/air location
```
