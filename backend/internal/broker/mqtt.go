package broker

import (
	"farm-scurity/domain/web"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MQTTRequest(mqttConf web.MQTTRequest) bool {
	broker := "tcp://localhost:1883"
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

	done := make(chan bool)

	token = client.Subscribe(mqttConf.Topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		if string(msg.Payload()) == mqttConf.MsgResp {
			client.Disconnect(250)
			done <- true
		}
		done <- false
	})

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	timeOut := 5

	select {
	case <-done:
		return true
	case <-time.After(time.Duration(timeOut) * time.Second):
		client.Disconnect(250)
		return false
	}
}
