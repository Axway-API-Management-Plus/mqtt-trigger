var mqtt = require('mqtt')
var express = require('express')
var bodyParser = require('body-parser')
var request = require('request-promise');
var os = require('os')

const BROKER="mqtt://mosquitto:1883"
const BROKER_USERNAME="guest"
const BROKER_PASSWORD="password"
const MQTT_TRIGGER_URL="http://mqtt-trigger:8080"
const APP_URL="http://"+os.hostname()+":3000"

var app = express()
app.use(bodyParser.json())
app.listen(3000, function () {
  console.log('Example app listening on port 3000!')
})

function expect(a,b) {
    if (a!=b) {
        throw new Error("Expecting "+a+" == "+ b)
    }
}

function mqtt_publish(topic, message) {
  return new Promise( function(resolve, reject) {
    const client = mqtt.connect(BROKER, { username: BROKER_USERNAME, password: BROKER_PASSWORD})
    client.on('connect', function () {
        client.publish(topic, message, { qos:1 }, (err) => {
          if (err) {
            return reject(err)
          }
          resolve("ok")
          //client.close()
        })
    })
    client.on('error', function (err) {
      reject(err)
    })
  })
}

function mqtt_subscribe(topic) {
  return new Promise( function(resolve, reject) {
    console.log("subscribe connect")
    const client = mqtt.connect(BROKER, { username: BROKER_USERNAME, password: BROKER_PASSWORD})
    client.on('connect', function () {
        console.log("subscribe connected")
        client.subscribe(topic, { qos:1 }, (err) => {
          console.log("subscribe subscribed")
          if (err) {
            return reject(err)
          }
          resolve(client)
        })
    })
    client.on('error', function (err, granted) {
      console.log("subscribe error")
      reject(err)
    })
    client.on('close', function (err) {
      console.log("subscribe close")
      reject(err)
    })
  })
}

describe('mqtt-trigger', () => {
    it('test mqtt connectivity', () => {
      return new Promise( async function(resolve, reject) {
        console.log("subscribe to presence")
        const client=await mqtt_subscribe("presence");
        console.log("install message handler and wait")
        client.on('message', function (topic, message) {
          console.log("got a message")
          let ok=false
          try {
            expect(topic, 'presence')
            expect(message.toString(), 'Hello mqtt')
            ok=true
          } catch(e) {
            reject(new Error("error"+e))
          }
          if(ok) resolve("OK")
        })
        return await mqtt_publish("presence", "Hello mqtt");
      })
    })
    it("test webserver loop", () => {
      return new Promise( async (resolve, reject) => {
        app.post("/loop", (req, res) => {
           console.log("/loop call body=", req.body)
           if(req.body.msg == "loop-message") {
             resolve("")
           } else {
             reject(new Error("bad message"))
           }
           res.status(200).send(req.body)
        });
        console.log("call app_url /loop call")
        const resp=await request( APP_URL+"/loop", { method: "POST", body: { msg:"loop-message"}, json:true  } );
        console.log("call app_url /loop call", resp)
        await mqtt_publish("presence", "Hello mqtt");
        console.log("hello mqtt published");
      })
    })

    it('configure and get message', () => {
      return new Promise( async (resolve, reject) => {
        var conf= [{
          name: "test01",
          topic: "test01-topic",
          url: APP_URL+"/test01-api",
          headers: [ "content-type: application/json" ],
          broker: BROKER,
          clientId: "test01-clientId",
          username: "guest",
          password: "guest",
        }];
        app.post("/test01-api", (req, res) => {
          console.log("got message xxx", req.body)
           if(req.body.msg == "test01-message") {
             resolve("")
           } else {
             reject(new Error("bad message"))
           }
           res.status(200).send({})
        });
        const resp=await request(MQTT_TRIGGER_URL+"/api/mqtt-triggers", { method: "POST", body:conf[0], json:true} );
        console.log("trigger post resp", resp)
        await mqtt_publish("test01-topic", '{"msg": "test01-message"}');
      })
    })
})
