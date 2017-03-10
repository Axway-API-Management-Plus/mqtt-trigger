# MQTT-Trigger

`mqtt-trigger` subscribe to MQTT topic and call an REST API

## Start
Command Line
```
mqtt-trigger --port <PORT:8080> --etcd-urls <ETCD_URLS>
```

Test full environment with docker
```
docker-compose -f docker-compose.yml up
```

## Configure

Use API to configure client...

## Build standalone binary:
Prerequisites : `golang`
```sh
make install-deps
make
```

## Build docker image

```
docker build -t mqtt-trigger:dev .
```

```
make docker
  -or-
docker build -t mqtt-trigger:dev .
docker run --rm mqtt-trigger:dev tar cz mqtt-trigger >mqtt-trigger.tar.gz
docker build -t mqtt-trigger -f Dockerfile.small .
```

## Test

```
make docker-test
```

### Trigger conf
```javascript
{
  name:    "string",
  topic:   "string",
  url:     "string",     // Todo : add default base
  headers: [ "string" ], // Todo : add default "content-type: application/json"

  broker:   "string",    // Todo : add default broker
  clientid: "string",    // Todo : add default (trigger-name)
  username: "string",    // Todo : add default broker
  password: "string",    // Todo : add default broker
}
```

## API
### GET /mqtt-triggers
List all triggers

### GET /mqtt-triggers/{trigger-name}
Get a single trigger

### POST /mqtt-triggers
Create a trigger

### DELETE /mqtt-triggers
Remove all triggers

### DELETE /mqtt-triggers/{trigger-name}
Remove a specific trigger

### GET /status
Get node status

## Todo
- Add default configuration
- Add a static configuration file
- Add a web based configuration file, with configurable refresh (30s)
- Add a disable REST Api hook
- Add a disable etcd : store config in file (json-like)
- Coalesce Broker configuration (Broker section ?)
- Distribute work across nodes

## Limitations
- One Node only !!!! : No distribution across nodes
- ClientID MUST be unique !
