package main

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

type Conffile struct {
	Defaults TriggerDefaults
	Triggers []*TriggerConf
}

func (server *Server) TriggerConfFileInit(filename string) error {
	log.Println("Loading conffile", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorln("Error opening conf file", err)
		return err
	}
	conffile := Conffile{}
	err = yaml.Unmarshal(data, &conffile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("--- t:\n%v\n\n", conffile)

	for _, trigger := range conffile.Triggers {
		runtimeTriggerSet(trigger, &conffile.Defaults)
	}

	return nil
}

func (server *Server) TriggerConfWatch(filename string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					server.TriggerConfFileInit(filename)
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		return err
	}
	<-done
	return nil
}
