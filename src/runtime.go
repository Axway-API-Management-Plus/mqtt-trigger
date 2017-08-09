package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type TriggerRuntime struct {
	TriggerConf
	origin string
	status string
	client MQTT.Client
	hash   string
	state  string
}

var triggers map[string]*TriggerRuntime

func TriggerRuntimeInit(server *Server) {
	triggers = make(map[string]*TriggerRuntime)
}

func TriggerInit(trigger *TriggerRuntime, origin string, triggerConf *TriggerConf, triggerDefault *TriggerDefaults) {
	trigger.TriggerConf = *triggerConf
	trigger.origin = origin

	if trigger.Username == "" {
		trigger.Username = triggerDefault.Username
	}

	if trigger.Password == "" {
		trigger.Password = triggerDefault.Password
	}

	if trigger.Broker == "" {
		trigger.Broker = triggerDefault.Broker
	}

	if trigger.URL == "" {
		trigger.URL = triggerDefault.URL + "/" + trigger.Name
	} else if trigger.URL[0] == '/' {
		trigger.URL = triggerDefault.URL + trigger.URL
	}

	if trigger.ClientId == "" {
		trigger.ClientId = triggerDefault.ClientID + "-" + trigger.Name
	}

	if trigger.Topic == "" {
		trigger.Topic = trigger.Name
	}

	if len(trigger.Headers) == 0 {
		trigger.Headers = triggerDefault.Headers
	}

	trigger.hash = trigger.URL +
		strings.Join(trigger.Headers, "") +
		trigger.Broker +
		trigger.Username +
		trigger.Password +
		trigger.ClientId +
		trigger.Topic
	log.Println(triggerLogPrefix+" Initialized - ", trigger)
}

func runtimeTriggerSet(origin string, triggerConf *TriggerConf, triggerDefaults *TriggerDefaults) {
	trigger := new(TriggerRuntime)

	TriggerInit(trigger, origin, triggerConf, triggerDefaults)

	if triggers[trigger.Name] == nil || triggers[trigger.Name].hash != trigger.hash {
		if triggers[trigger.Name] == nil {
			runtimeTriggerDelete(trigger.Name)
		}
		triggers[trigger.Name] = trigger
		log.Println(triggerLogPrefix+" Runtime create - ", triggerConfPath, trigger.Name, trigger)
		trigger.state = "installed"
		go trigger.installWithRetry()
	} else {
		log.Println(triggerLogPrefix+" Runtime not changed - ", triggerConfPath, trigger.Name, trigger)
	}
}

func runtimeTriggerDelete(name string) {
	trigger := triggers[name]
	if trigger != nil {
		trigger.uninstall()
	}
	delete(triggers, name)

	log.Println(triggerLogPrefix+" Runtime delete - ", triggerConfPath, name, trigger)
}

func (t *TriggerRuntime) processMessage(mqttclient MQTT.Client, msg MQTT.Message) {
	log.Println(triggerLogPrefix, "Process Message", t.Name, t.ClientId, msg.MessageID(), msg.Topic(), t.URL, string(msg.Payload()))

	client := &http.Client{}
	req, err := http.NewRequest("POST", t.URL, bytes.NewReader(msg.Payload()))
	for _, header := range t.Headers {
		splits := strings.SplitN(header, ":", 2)
		key := strings.TrimSpace(splits[0])
		value := ""
		if len(splits) > 1 {
			value = strings.TrimSpace(splits[1])
		}
		req.Header.Add(key, value)
	}
	/*if t.Headers != nil {
		for key, value := range *t.Headers {
			req.Header.Add(key, value)
		}
	}*/
	req.Header.Add("TOPIC", msg.Topic())

	resp, err := client.Do(req)
	//defer resp.Body.Close()

	if err != nil {
		log.Errorln(triggerLogPrefix+" Error sending message", msg.MessageID(), msg.Topic(), err)
		return
	}

	log.Println(triggerLogPrefix, "Response", t.Name, t.ClientId, msg.MessageID(), msg.Topic(), resp.StatusCode, resp.Status, "(", t.URL, ")")
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		log.Errorln(triggerLogPrefix, "Response Body", string(body))
	}
}

func (t *TriggerRuntime) installWithRetry() {
	delay := 1 * time.Second
	for t.state == "installed" {
		err := t.install()
		if err == nil {
			break
		}
		time.Sleep(delay)
		delay = delay * 2
		if delay > 60*time.Second {
			delay = 60 * time.Second
		}
	}
}

func (t *TriggerRuntime) install() error {
	if strings.HasPrefix(t.Broker, "mqtt://") {
		t.Broker = "tcp://" + strings.TrimPrefix(t.Broker, "mqtt://")
	}
	opts := MQTT.NewClientOptions().AddBroker(t.Broker)
	opts.SetClientID(t.ClientId)
	opts.SetUsername(t.Username)
	opts.SetPassword(t.Password)
	opts.SetDefaultPublishHandler(t.processMessage)
	log.Println(triggerLogPrefix+" Runtime - Connecting to MQTT Broker",
		"broker=", t.Broker,
		"clientId=", t.ClientId,
		"username=", t.Username,
		"password=", t.Password,
		"topic=", t.Topic,
	)
	//create and start a client using the above ClientOptions
	t.client = MQTT.NewClient(opts)
	if token := t.client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorln(triggerLogPrefix+" Runtime - Connect", t.Broker, "as", t.ClientId, token.Error())
		t.client = nil
		return token.Error()
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := t.client.Subscribe(t.Topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Errorln(triggerLogPrefix+"Runtime - Subscribe", t.ClientId, t.Topic, token.Error())
		t.client.Disconnect(0)
		t.client = nil
		return token.Error()
	}

	log.Println(triggerLogPrefix+" Runtime - Connected to MQTT Broker",
		"clientId=", t.ClientId,
		"username=", t.Username,
		"password=", t.Password,
		"topic=", t.Topic,
	)
	return nil
}

func (t *TriggerRuntime) uninstall() {
	t.state = "uninstalled"
	if t.client != nil {
		if token := t.client.Unsubscribe(t.Topic); token.Wait() && token.Error() != nil {
			log.Errorln(triggerLogPrefix+"Runtime - Unsubscribe", t.ClientId, t.Topic, token.Error())
		}

		t.client.Disconnect(0)
	}
}
