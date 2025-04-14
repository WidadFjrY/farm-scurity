package com.example.farmscurity.service

import android.content.Context
import android.content.pm.PackageManager
import androidx.core.app.NotificationCompat
import androidx.core.app.NotificationManagerCompat
import androidx.core.content.ContextCompat
import com.hivemq.client.mqtt.MqttClient
import com.hivemq.client.mqtt.MqttGlobalPublishFilter
import com.hivemq.client.mqtt.datatypes.MqttQos
import kotlinx.coroutines.*

class MQTTService(private val context: Context) {

    private val mqttClient = MqttClient.builder()
        .useMqttVersion3()
        .serverHost("192.168.1.6")
        .serverPort(1883)
        .buildBlocking()

    fun connect() {
        mqttClient.connect()
    }

    fun subscribe(topic: String) {
        mqttClient.subscribeWith()
            .topicFilter(topic)
            .qos(MqttQos.AT_LEAST_ONCE)
            .send()

        GlobalScope.launch(Dispatchers.IO) {
            val publishFlow = mqttClient.publishes(MqttGlobalPublishFilter.ALL)
            while (true) {
                val publish = publishFlow.receive()
                val message = String(publish.payloadAsBytes)

                withContext(Dispatchers.Main) {
                    showNotification("Pesan MQTT", message)
                }
            }
        }

    }

    private fun showNotification(title: String, message: String) {
        val channelId = "mqtt_channel_id"
        val notificationId = (System.currentTimeMillis() % 10000).toInt()

        // Cek apakah permission sudah diberikan
        if (ContextCompat.checkSelfPermission(context, android.Manifest.permission.POST_NOTIFICATIONS)
            == PackageManager.PERMISSION_GRANTED) {

            val notification = NotificationCompat.Builder(context, channelId)
                .setSmallIcon(android.R.drawable.ic_dialog_info)
                .setContentTitle(title)
                .setContentText(message)
                .setPriority(NotificationCompat.PRIORITY_HIGH)
                .setAutoCancel(true)
                .build()

            NotificationManagerCompat.from(context).notify(notificationId, notification)
        }
    }



    fun disconnect() {
        mqttClient.disconnect()
    }
}
