# MQTT-Trigger

`mqtt-trigger` subscribe to MQTT topic and call an REST API

## Start
Command Line
```txt
mqtt-trigger [OPTIONS]
Usage of ./mqtt-trigger:
  -conf string
    	conffile
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
docker build -t mqtt-trigger .
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

## Todo
- Add TLS for the trigger
- Add TLS for the MQTT trigger
- Add a web based configuration file, with configurable refresh (30s)
- Add auto refresh of configuration (inotify?)

## Limitations
- One Node only !!!! : No distribution across nodes
- ClientID MUST be unique !

## Changelog
- 0.0.2
  - configuration file support with default
  - disabled etcd support
  - use docker 17.05 for compact build

## Contributing

Please read [Contributing.md](https://github.com/Axway-API-Management-Plus/Common/blob/master/Contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Team

![alt text][Axwaylogo] Axway Team

[Axwaylogo]: https://github.com/Axway-API-Management/Common/blob/master/img/AxwayLogoSmall.png  "Axway logo"


## License
[Apache License 2.0](/LICENSE)
