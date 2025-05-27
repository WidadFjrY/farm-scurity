package broker

import (
	"farm-scurity/domain/web"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MQTTRequest(mqttConf web.MQTTRequest, isNotification bool) (bool, string) {
	broker := "tcp://broker.hivemq.com:1883"
	var payload string
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(mqttConf.ClientId)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Publish(mqttConf.Topic, 1, false, mqttConf.Payload)
	token.Wait()

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if !isNotification {
		done := make(chan bool)

		token = client.Subscribe(mqttConf.Topic, 1, func(client mqtt.Client, msg mqtt.Message) {
			if strings.HasPrefix(string(msg.Payload()), "ok") {
				client.Disconnect(250)
				payload = string(msg.Payload())

				done <- true
			}
			payload = string(msg.Payload())
			done <- false
		})

		if token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		timeOut := 10

		select {
		case <-done:
			return true, payload
		case <-time.After(time.Duration(timeOut) * time.Second):
			client.Disconnect(250)
			return false, payload
		}
	}
	return true, payload
}
