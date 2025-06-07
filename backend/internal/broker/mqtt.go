package broker

import (
	"farm-scurity/domain/web"
	"fmt"
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
		fmt.Println("MQTT connect error:", token.Error())
		return false, ""
	}

	pubToken := client.Publish(mqttConf.Topic, 1, false, mqttConf.Payload)
	if pubToken.Wait() && pubToken.Error() != nil {
		fmt.Println("MQTT publish error:", pubToken.Error())
		client.Disconnect(250)
		return false, ""
	}

	if !isNotification {
		done := make(chan bool, 1)

		subToken := client.Subscribe(mqttConf.Topic, 1, func(client mqtt.Client, msg mqtt.Message) {
			payload = string(msg.Payload())
			if strings.HasPrefix(payload, "ok") {
				done <- true
			} else {
				done <- false
			}
		})

		if subToken.Wait() && subToken.Error() != nil {
			fmt.Println("MQTT subscribe error:", subToken.Error())
			client.Disconnect(250)
			return false, ""
		}

		select {
		case success := <-done:
			client.Disconnect(250)
			return success, payload
		case <-time.After(10 * time.Second):
			fmt.Println("MQTT timeout waiting for message")
			client.Disconnect(250)
			return false, ""
		}
	}

	client.Disconnect(250)
	return true, payload
}
