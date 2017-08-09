var mqtt = require('mqtt')
var express = require('express')
var bodyParser = require('body-parser')
var request = require('request-promise');
var os = require('os')

const BROKER = "mqtt://mosquitto:1883"
const BROKER_USERNAME = "guest"
const BROKER_PASSWORD = "password"
const MQTT_TRIGGER_URL = "http://mqtt-trigger:8080"
const APP_URL = "http://" + os.hostname() + ":3000"

var app = express()
app.use(bodyParser.json())
app.use((req, resp, next) => {
  console.log("[test] received : " + req.method + " " + req.originalUrl);
  next()
})

app.listen(3000, function () {
  console.log('[test] Example app listening on port 3000!')
})

function expect(a, b) {
    if (a != b) {
        throw new Error("Expecting " + a + " == " + b)
    }
}

function mqtt_publish(topic, message) {
  return new Promise(function(resolve, reject) {
    const client = mqtt.connect(BROKER, { username: BROKER_USERNAME, password: BROKER_PASSWORD })
    client.on('connect', function () {
        client.publish(topic, message, { qos: 1 }, (err) => {
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
  return new Promise(function(resolve, reject) {
    console.log("[test] subscribe connect", topic)
    const client = mqtt.connect(BROKER, { username: BROKER_USERNAME, password: BROKER_PASSWORD })
    client.on('connect', function () {
        console.log("[test] subscribe connected", topic)
        client.subscribe(topic, { qos: 1 }, (err) => {
          console.log("[test] subscribe subscribed", topic)
          if (err) {
            return reject(err)
          }
          resolve(client)
        })
    })
    client.on('error', function (err, granted) {
      console.log("[test] subscribe error", topic)
      reject(err)
    })
    client.on('close', function (err) {
      console.log("[test] subscribe close", topic)
      reject(err)
    })
  })
}

describe('mqtt-trigger', () => {
    it('test mqtt connectivity', () => {
      return new Promise(async function(resolve, reject) {
        console.log("[test] subscribe to presence")
        const client = await mqtt_subscribe("presence");
        console.log("install message handler and wait")
        client.on('message', function (topic, message) {
          console.log("[test] got a message")
          let ok = false
          try {
            expect(topic, 'presence')
            expect(message.toString(), 'Hello mqtt')
            ok = true
          } catch (e) {
            reject(new Error("error" + e))
          }
          if (ok) resolve("OK")
        })
        return await mqtt_publish("presence", "Hello mqtt");
      })
    })
    it("test webserver loop", () => {
      return new Promise(async (resolve, reject) => {
        app.post("/loop", (req, res) => {
           console.log("[test] /loop call body=", req.body)
           if (req.body.msg == "loop-message") {
             resolve("")
           } else {
             reject(new Error("bad message"))
           }
           res.status(200).send(req.body)
        });
        console.log("[test] call app_url /loop call")
        const resp = await request(APP_URL + "/loop", { method: "POST", body: { msg: "loop-message" }, json: true });
        console.log("[test] call app_url /loop call", resp)
        await mqtt_publish("presence", "Hello mqtt");
        console.log("[test] hello mqtt published");
      })
    })

    it.skip('configure and get message (etcd mode)', () => {
      return new Promise(async (resolve, reject) => {
        var conf = [{
          name: "test01",
          topic: "test01-topic",
          url: APP_URL + "/test01-api",
          headers: [ "content-type: application/json" ],
          broker: BROKER,
          clientId: "test01-clientId",
          username: "guest",
          password: "guest",
        }];
        app.post("/test01-api", (req, res) => {
          console.log("[test] got message xxx", req.body)
           if (req.body.msg == "test01-message") {
             resolve("")
           } else {
             reject(new Error("bad message"))
           }
           res.status(200).send({})
        });
        const resp = await request(MQTT_TRIGGER_URL + "/api/mqtt-triggers", { method: "POST", body: conf[0], json: true });
        console.log("[test] trigger post resp", resp)
        await mqtt_publish("test01-topic", '{"msg": "test01-message"}');
      })
    })

    it('test default config file', () => {
      return new Promise(async (resolve, reject) => {
        app.post("/api/topic/simplest", (req, res) => {
          console.log("[test] got message xxx", req.body)
           if (req.body.msg !== "simplest") {
             reject(new Error("bad message"))
           } else if (req.header("topic") !== "simplest") {
             reject(new Error("bad topic"))
           } else {
             resolve("")
           }
           res.status(200).send({})
        });
        console.log("[test] mqtt trigger")
        await mqtt_publish("simplest", '{"msg": "simplest"}');
      })
    })
})
