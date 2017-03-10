package main

import (
	"strings"
	"sync"

	"github.com/namsral/flag"

	tools "./tools"
	log "github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var conffile string

func LogInit() {
	//log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	formatter := new(prefixed.TextFormatter)
	formatter.DisableTimestamp = false
	formatter.ShortTimestamp = false
	formatter.TimestampFormat = "2006-01-02 15:04:05.000000000"
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)
}

func main() {
	var port int
	var etcdURLs string
	flag.IntVar(&port, "port", 8080, "api port")
	flag.StringVar(&etcdURLs, "etcd-urls", "http://localhost:2379", "urls to etcd")
	flag.Parse()

	LogInit()

	etcds := strings.Split(etcdURLs, ",")

	log.Println(triggerLogPrefix+" Using etcd %v", etcds)
	server := &Server{}
	server.Config = tools.NewEtcdConfig(etcds)

	var wg sync.WaitGroup
	wg.Add(1)

	log.Println(triggerLogPrefix+" initializing etcd", triggerConfPath)
	server.Config.Init()
	server.Config.Wait()

	log.Println(triggerLogPrefix+" initializing runtime", triggerConfPath)
	TriggerRuntimeInit(server)

	log.Println(triggerLogPrefix+" loading shared configuration (etcd)", triggerConfPath)
	TriggerConfInit(server)

	log.Println(triggerLogPrefix+" watching conf changes", triggerConfPath)
	go TriggerConfWatch(server)

	log.Println(triggerLogPrefix+" setting up API", triggerApiPath)
	TriggerApiInit(server)

	log.Println(triggerLogPrefix+" listening to http port", port)
	ServerListenAndServe(server, port)

	wg.Wait()
}
