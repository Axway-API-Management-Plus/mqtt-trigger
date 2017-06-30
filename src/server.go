package main

import "github.com/gorilla/mux"

type TriggerDefaults struct {
	URL      string   //base url (url+triggerName)
	Headers  []string //default headers ("Content-Type: application/json")
	Broker   string   //default broker url
	ClientID string   //default base for ClientID (clientID+"-"+triggerName)
	Username string   //default username for broker
	Password string   //default password for broker
}

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

type Server struct {
	Mux            *mux.Router
	TriggerDefault TriggerDefaults
}
