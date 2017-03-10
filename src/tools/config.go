package tools

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/coreos/etcd/client"
)

type EtcdConfig struct {
	kapi  client.KeysAPI
	mapi  client.MembersAPI
	index uint64
	urls  []string
}

func NewEtcdConfig(urls []string) *EtcdConfig {
	conf := new(EtcdConfig)
	conf.urls = urls
	return conf
}

func (e *EtcdConfig) Init() {
	log.Printf("Using etcd : %v", e.urls)
	cfg := client.Config{
		Endpoints: e.urls[:],
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: 3 * time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	e.kapi = client.NewKeysAPI(c)
	e.mapi = client.NewMembersAPI(c)

}

func (e *EtcdConfig) Wait() {
	ctx := context.Background()
	ms := 10 * time.Millisecond
	for {
		_, err := e.mapi.List(ctx)
		if err == nil {
			break
		}
		log.Errorln("Error getting etcd members, retrying....", ms/time.Millisecond, err)
		ms = ms * 2
		time.Sleep(ms)
	}
}

func (e *EtcdConfig) GetAllCollection(collectionPath string, collection interface{}) error {
	options := client.GetOptions{Recursive: true}
	resp, err := e.kapi.Get(context.Background(), collectionPath, &options)
	if err != nil {
		if client.IsKeyNotFound(err) {
			log.Errorln("Config  - Getting all collection", collectionPath, err)
			return nil
		}
		return err
	}

	s := "["
	for _, node := range resp.Node.Nodes {
		if s == "[" {
			s += node.Value
		} else {
			s += "," + node.Value
		}
	}
	s += "]"

	err = json.Unmarshal([]byte(s), collection)
	log.Println("Config - Getting all collection", collectionPath, s, err)
	return err
}

func (e *EtcdConfig) GetCollectionItem(collectionPath string, id string, item interface{}) error {
	resp, err := e.kapi.Get(context.Background(), collectionPath+"/"+id, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			log.Errorln("Config  - GetCollectionItem", collectionPath, id, err)
			return err
		}
		return nil
	}
	json.Unmarshal([]byte(resp.Node.Value), item)
	return nil
}

func (e *EtcdConfig) SetCollectionItem(collectionPath string, id string, item interface{}) error {
	js, err := json.Marshal(item)
	resp, err := e.kapi.Set(context.Background(), collectionPath+"/"+id, string(js), nil)
	if err != nil {
		log.Errorln("Config - SetCollectionItem : ", collectionPath, id, err)
		return err
	}
	log.Printf("Config - SetCollectionItem : %q key has %q value\n", resp.Node.Key, resp.Node.Value)
	return nil
}

func (e *EtcdConfig) DeleteCollectionItem(collectionPath string, id string) error {
	_, err := e.kapi.Delete(context.Background(), collectionPath+"/"+id, nil)
	if err != nil {
		log.Errorln("Config - DeleteCollectionItem : %q key has %q value\n", collectionPath, id, err)
		return err
	}
	log.Printf("Config - DeleteCollectionItem : %q key has %q value\n", collectionPath, id)
	return nil
}

func (e *EtcdConfig) DeleteCollection(collectionPath string) error {
	_, err := e.kapi.Delete(context.Background(), collectionPath, &client.DeleteOptions{Dir: true, Recursive: true})
	if err != nil {
		log.Errorln("Config - DeleteCollection : %s\n", collectionPath, err)
		return err
	}
	log.Printf("Config - DeleteCollectionItem : %s\n", collectionPath)
	return nil
}

func (e *EtcdConfig) WatchCollection(collectionPath string, index *uint64, item interface{}) (string, string, error) {
	retry := 1
	for {
		options := client.WatcherOptions{Recursive: true, AfterIndex: *index}
		w := e.kapi.Watcher(collectionPath, &options)
		resp, err := w.Next(context.Background())
		if err != nil {
			retry = retry * 2
			if retry > 20 {
				log.Fatalln("Config Watch - Next Error", collectionPath, retry, err)
				panic("Config Watch")
			}
			log.Errorln("Config Watch - Next Error", collectionPath, retry, "ms", err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		retry = 1
		*index = resp.Index + 1

		if resp.Action == "set" {
			err = json.Unmarshal([]byte(resp.Node.Value), item)
			if err != nil {
				log.Errorln("Config Watch - Error unmarchalling collection : ", collectionPath, resp.Action, resp.Node.Key, resp.Node.Value, err)
			} else {
				log.Println("Config Watch - collection updated/created", collectionPath, resp.Action, resp.Node.Key, resp.Node.Value)
				/*for _, value := range e.notifications {
					value <- item
				}*/
				return resp.Action, resp.Node.Key, nil
				//log.Println("Config Watch - Service updated/created", resp.Node.Key, resp.Node.Value)
			}
		} else {
			log.Println("Config Watch - collection other", collectionPath, resp.Action, resp.Node.Key, resp.Node.Value)
			return resp.Action, resp.Node.Key, nil
		}
	}
}
