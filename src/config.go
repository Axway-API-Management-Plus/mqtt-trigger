package main

import log "github.com/Sirupsen/logrus"

type TriggerConf struct {
	Name string `jsonapi:"primary,mqtt"`

	Topic   string   `jsonapi:"attr,topic,omitempty"`
	URL     string   `jsonapi:"attr,url,omitempty"`
	Headers []string `jsonapi:"attr,headers,omitempty"`

	Broker   string `jsonapi:"attr,broker,omitempty"`
	ClientId string `jsonapi:"attr,clientid,omitempty"`
	Username string `jsonapi:"attr,username,omitempty"`
	Password string `jsonapi:"attr,password,omitempty"`
}

var triggerConfPath = "/config/mqtt-triggers"
var triggerApiPath = "/api/mqtt-triggers"
var triggerLogPrefix = "Trigger MQTT"

func TriggerConfInit(server *Server) {
	var triggers []*TriggerConf

	err := server.Config.GetAllCollection(triggerConfPath, &triggers)
	if err != nil {
		log.Errorln(triggerLogPrefix+" Conf - Cannot get trigger configuration :", err)
		//panic(triggerLogPrefix + " Conf - Cannot get trigger configuration")
	}
	for _, trigger := range triggers {
		runtimeTriggerSet(trigger)
	}
}

func TriggerConfWatch(server *Server) {
	var index uint64

	for {
		var trigger TriggerConf
		action, key, err := server.Config.WatchCollection(triggerConfPath, &index, &trigger)
		if err != nil {
			log.Errorln(triggerLogPrefix+" Conf watch error", triggerConfPath, err)
		} else {
			if action == "delete" {
				log.Println(triggerLogPrefix+" Conf watch - delete -", triggerConfPath, key, trigger)
				runtimeTriggerDelete(key)
			} else if action == "set" {
				log.Println(triggerLogPrefix+" Conf watch - set -", triggerConfPath, key, trigger.Name, trigger)
				runtimeTriggerSet(&trigger)
			} else {
				log.Errorln(triggerLogPrefix+" Conf watch - unknown action", triggerConfPath, key, action)
			}
		}
	}
}
