# MQTT-Trigger

`mqtt-trigger` subscribe to MQTT topic and call an REST API

## Start
Command Line
```txt
mqtt-trigger [OPTIONS]
Usage of ./mqtt-trigger:
  -conf string
    	conffile
  -etcd-urls string
    	urls to etcd (default "http://localhost:2379")
  -http-headers string
    	headers ( 'key1 : value1, key2: value2' )
  -http-url string
    	Default prefix for url to forward messages ( http-url + trigger-name)
  -mqtt-broker string
    	Default MQTT broker url (mqtt://)
  -mqtt-client-id string
    	Default prefix to MQTT Client ID (client-id + '-' + trigger-name)
  -mqtt-password string
    	Default MQTT password
  -mqtt-username string
    	Default MQTT username
  -port int
    	api port (default 8080)
```

Test full environment with docker
```
docker-compose -f docker-compose.yml up
```

## Configuration

### By Configuration file
```yaml

defaults:
  url: http://localhost:8080/api/topic/
  headers:
  - "Content-Type: application/json"
  broker: mqtt://localhost:1883
  clientid: mqtt-trigger
  username: mqtt-trigger
  password : goodpass

triggers:
- name: simplest
- name: all
  topic: "/#-*"
- name: activemq
  topic: activemq/#
  broker: mqtt://localhost:8883

```


### By API (with etcd enabled for persistence)
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

## Trigger conf
```javascript
{
  name:    "string",
  topic:   "string",     // Default: <name>
  url:     "string",     // Default (conffile.Defaults.url || -mqtt-url ) <name>
  headers: [ "string" ], // Default: "content-type: application/json" || conffile.Defaults.headers || -mqtt-headers

  broker:   "string",    // Default: conffile.Defaults.Broker || -mqtt-broker
  clientid: "string",    // Default: (conffile.Defaults.ClientID || -mqttclient-id) + <name>)
  username: "string",    // Default: conffile.Defaults.Username || -mqtt-username
  password: "string",    // Default: conffile.Defaults.Password || -mqtt-password
}
```

## API
#### GET /mqtt-triggers
List all triggers

#### GET /mqtt-triggers/{trigger-name}
Get a single trigger

#### POST /mqtt-triggers
Create a trigger

#### DELETE /mqtt-triggers
Remove all triggers

#### DELETE /mqtt-triggers/{trigger-name}
Remove a specific trigger

#### GET /status
Get node status

## Todo
- [Done]Add default configuration
- [Done] Add a static configuration file
- Add a web based configuration file, with configurable refresh (30s)
- Add a disable REST Api hook
- Add a disable etcd : store config in file (json-like)
- [Done] Coalesce Broker configuration (Broker section ?)
- Distribute work across nodes

## Limitations
- One Node only !!!! : No distribution across nodes
- ClientID MUST be unique !
