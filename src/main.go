package main

import (
	"strings"
	"sync"

	"github.com/namsral/flag"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func LogInit() {
	//log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	formatter := new(prefixed.TextFormatter)
	formatter.DisableTimestamp = false
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05.000000000"
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)
}

var Version string
var Build string
var Date string

func main() {
	//var port int
	var conffile string

	var HTTPUrl string
	var HTTPHeaders string
	var MQTTBroker string
	var MQTTClientID string
	var MQTTUsername string
	var MQTTPassword string

	//flag.IntVar(&port, "port", 8080, "api port")
	//flag.StringVar(&etcdURLs, "etcd-urls", "", "urls to etcd ")
	flag.StringVar(&conffile, "conf", "./mqtt-trigger.yml", "conffile")

	flag.StringVar(&HTTPUrl, "http-url", "", "Default prefix for url to forward messages ( http-url + trigger-name)")
	flag.StringVar(&HTTPHeaders, "http-headers", "", "headers ( 'key1 : value1, key2: value2' )")
	flag.StringVar(&MQTTBroker, "mqtt-broker", "", "Default MQTT broker url (mqtt://)")
	flag.StringVar(&MQTTClientID, "mqtt-client-id", "", "Default prefix to MQTT Client ID (client-id + '-' + trigger-name)")
	flag.StringVar(&MQTTUsername, "mqtt-username", "", "Default MQTT username")
	flag.StringVar(&MQTTPassword, "mqtt-password", "", "Default MQTT password")

	flag.Parse()

	LogInit()
	log.Println("Starting mqtt-proxy - version:", Version, "build:", Build, " date:", Date)
	//etcds := strings.Split(etcdURLs, ",")
	var wg sync.WaitGroup
	wg.Add(1)
	server := &Server{}
	TriggerRuntimeInit(server)

	log.Println(triggerLogPrefix+" loading configuration", conffile)
	server.TriggerConfFileInit(conffile)
	server.TriggerConfWatch(conffile)

	if HTTPUrl != "" {
		server.TriggerDefault.URL = HTTPUrl
	}

	if HTTPHeaders != "" {
		server.TriggerDefault.Headers = strings.Split(HTTPHeaders, ",")
	}

	if MQTTBroker != "" {
		server.TriggerDefault.Broker = MQTTBroker
	}

	if MQTTClientID != "" {
		server.TriggerDefault.ClientID = MQTTClientID
	}

	if MQTTUsername != "" {
		server.TriggerDefault.Username = MQTTUsername
	}

	if MQTTPassword != "" {
		server.TriggerDefault.Password = MQTTPassword
	}

	/*if etcdURLs != "" {
		log.Println(triggerLogPrefix+" initializing etcd link %s %v", triggerConfPath, etcds)
		server.Config = tools.NewEtcdConfig(etcds)
		server.Config.Init()
		server.Config.Wait()

		log.Println(triggerLogPrefix+" initializing runtime", triggerConfPath)

		log.Println(triggerLogPrefix+" loading shared configuration (etcd)", triggerConfPath)
		TriggerConfInit(server)

		log.Println(triggerLogPrefix+" watching conf changes", triggerConfPath)
		go TriggerConfWatch(server)
	}*/

	//log.Println(triggerLogPrefix+" setting up API", triggerApiPath)
	//TriggerApiInit(server)

	//log.Println(triggerLogPrefix+" listening to http port", port)
	//ServerListenAndServe(server, port)

	wg.Wait()
}
