package com.example.farmscurity.service

import android.Manifest
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.Service
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Build
import android.os.IBinder
import android.util.Log
import androidx.core.app.ActivityCompat
import androidx.core.app.NotificationCompat
import androidx.core.app.NotificationManagerCompat
import com.hivemq.client.mqtt.MqttClient
import com.hivemq.client.mqtt.MqttGlobalPublishFilter
import com.hivemq.client.mqtt.datatypes.MqttQos
import kotlinx.coroutines.*
import kotlin.random.Random

class MQTTServiceNotification : Service() {
    private val channelId = "mqtt_channel_id"
    private val notificationId = 1

    private val mqttClient = MqttClient.builder()
        .useMqttVersion3()
        .serverHost("192.168.1.6")
        .serverPort(1883)
        .buildBlocking()

    private var mqttJob: Job? = null

    override fun onBind(intent: Intent?): IBinder? {
        return null
    }

    override fun onCreate() {
        super.onCreate()
        createNotificationChannel()
        startForegroundService()
        connect()
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                channelId,
                "Gerakan Terdeteksi!",
                NotificationManager.IMPORTANCE_HIGH
            )
            val notificationManager = getSystemService(NotificationManager::class.java)
            notificationManager.createNotificationChannel(channel)
        }
    }

    private fun startForegroundService() {
        val notification = NotificationCompat.Builder(this, channelId)
            .setContentTitle("FarmScurity")
            .setContentText("Menunggu deteksi gerakan...")
            .setSmallIcon(android.R.drawable.ic_dialog_info)
            .setPriority(NotificationCompat.PRIORITY_LOW)
            .build()

        startForeground(notificationId, notification)
    }

    private fun connect() {
        try {
            mqttClient.connect()
            Log.d("com.example.farmscurity.service.MQTTService", "✅ Connected to MQTT Broker")
            subscribe("broker/farm-security/notification")
        } catch (e: Exception) {
            Log.e("com.example.farmscurity.service.MQTTService", "❌ Failed to connect: ${e.message}")
        }
    }

    private fun subscribe(topic: String) {
        try {
            mqttClient.subscribeWith()
                .topicFilter(topic)
                .qos(MqttQos.AT_LEAST_ONCE)
                .send()
            Log.d("com.example.farmscurity.service.MQTTService", "✅ Subscribed to topic: $topic")

            mqttJob = CoroutineScope(Dispatchers.IO).launch {
                val publishFlow = mqttClient.publishes(MqttGlobalPublishFilter.ALL)
                while (isActive) {
                    val publish = publishFlow.receive()
                    val message = String(publish.payloadAsBytes)
                    Log.d("com.example.farmscurity.service.MQTTService", "📩 Received MQTT message: $message")

                    withContext(Dispatchers.Main) {
                        showNotification("Gerakan Terdeteksi‼️", message)
                    }
                }
            }
        } catch (e: Exception) {
            Log.e("com.example.farmscurity.service.MQTTService", "❌ Failed to subscribe: ${e.message}")
        }
    }

    private fun showNotification(title: String, message: String) {
        Log.d("com.example.farmscurity.service.MQTTService", "🔔 Showing notification: $title - $message")

        val notification = NotificationCompat.Builder(this, channelId)
            .setSmallIcon(android.R.drawable.ic_dialog_info)
            .setContentTitle(title)
            .setContentText(message)
            .setPriority(NotificationCompat.PRIORITY_HIGH)
            .setAutoCancel(true)
            .build()

        if (ActivityCompat.checkSelfPermission(this, Manifest.permission.POST_NOTIFICATIONS)
            != PackageManager.PERMISSION_GRANTED) {
            Log.e("com.example.farmscurity.service.MQTTService", "🚨 Notification permission NOT granted")
            return
        }

        NotificationManagerCompat.from(this).notify(Random.nextInt(), notification)
    }


    override fun onDestroy() {
        mqttJob?.cancel()
        mqttClient.disconnect()
        super.onDestroy()
    }
}
