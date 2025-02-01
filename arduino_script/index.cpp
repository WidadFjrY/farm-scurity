#include <WiFi.h>
#include <PubSubClient.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include "esp_camera.h"

// Sesuaikan nanti
const char *ssid = "NAMA_WIFI";
const char *password = "PASSWORD_WIFI";

const char *mqtt_server = ""; // Server Mosquitto
const int mqtt_port = 1883;
const char *mqtt_user = "user";
const char *mqtt_pass = "password";
const char *subscribe_topic = "broker/farm-scurity";
const char *publish_topic = "esp32-cam/status";
const char *api_url = "http://localhost:8080/api/motion-detected";

const int pirPin = 13;
const int speakerPin = 4;

// Konfigurasi Kamera
#define PWDN_GPIO_NUM 32
#define RESET_GPIO_NUM -1
#define XCLK_GPIO_NUM 0
#define SIOD_GPIO_NUM 26
#define SIOC_GPIO_NUM 27
#define Y9_GPIO_NUM 35
#define Y8_GPIO_NUM 34
#define Y7_GPIO_NUM 39
#define Y6_GPIO_NUM 36
#define Y5_GPIO_NUM 21
#define Y4_GPIO_NUM 19
#define Y3_GPIO_NUM 18
#define Y2_GPIO_NUM 5
#define VSYNC_GPIO_NUM 25
#define HREF_GPIO_NUM 23
#define PCLK_GPIO_NUM 22

WiFiClient espClient;
PubSubClient client(espClient);

void setup_camera()
{
    camera_config_t config;
    config.ledc_channel = LEDC_CHANNEL_0;
    config.ledc_timer = LEDC_TIMER_0;
    config.pin_d0 = Y2_GPIO_NUM;
    config.pin_d1 = Y3_GPIO_NUM;
    config.pin_d2 = Y4_GPIO_NUM;
    config.pin_d3 = Y5_GPIO_NUM;
    config.pin_d4 = Y6_GPIO_NUM;
    config.pin_d5 = Y7_GPIO_NUM;
    config.pin_d6 = Y8_GPIO_NUM;
    config.pin_d7 = Y9_GPIO_NUM;
    config.pin_xclk = XCLK_GPIO_NUM;
    config.pin_pclk = PCLK_GPIO_NUM;
    config.pin_vsync = VSYNC_GPIO_NUM;
    config.pin_href = HREF_GPIO_NUM;
    config.pin_sscb_sda = SIOD_GPIO_NUM;
    config.pin_sscb_scl = SIOC_GPIO_NUM;
    config.pin_pwdn = PWDN_GPIO_NUM;
    config.pin_reset = RESET_GPIO_NUM;
    config.xclk_freq_hz = 20000000;
    config.pixel_format = PIXFORMAT_JPEG;
    config.frame_size = FRAMESIZE_SVGA;
    config.jpeg_quality = 12;
    config.fb_count = 1;

    esp_err_t err = esp_camera_init(&config);
    if (err != ESP_OK)
    {
        Serial.printf("Camera init failed with error 0x%x", err);
        return;
    }
}

void setup_wifi()
{
    delay(10);
    Serial.println();
    Serial.print("Connecting to ");
    Serial.println(ssid);

    WiFi.begin(ssid, password);

    while (WiFi.status() != WL_CONNECTED)
    {
        delay(500);
        Serial.print(".");
    }

    Serial.println("");
    Serial.println("WiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
}

void callback(char *topic, byte *payload, unsigned int length)
{
    // Buat Debugging

    // Serial.print("Message arrived [");
    // Serial.print(topic);
    // Serial.print("]: ");

    String message;
    for (int i = 0; i < length; i++)
    {
        message += (char)payload[i];
    }

    if (String(topic) == subscribe_topic)
    {
        if (message == "capture")
        {
            captureAndSend();
        }
        else if (message == "speaker_on")
        {
            digitalWrite(speakerPin, HIGH);
        }
        else if (message == "speaker_off")
        {
            digitalWrite(speakerPin, LOW);
        }
    }
}

void reconnect()
{
    while (!client.connected())
    {
        Serial.print("Attempting MQTT connection...");
        if (client.connect("ESP32-CAM", mqtt_user, mqtt_pass))
        {
            Serial.println("connected");
            client.subscribe(subscribe_topic);
        }
        else
        {
            Serial.print("failed, rc=");
            Serial.print(client.state());
            Serial.println(" try again in 5 seconds");
            delay(5000);
        }
    }
}

void captureAndSend()
{
    fb = esp_camera_fb_get();
    if (!fb)
    {
        Serial.println("Camera capture failed");
        return;
    }

    HTTPClient http;
    http.begin("http:/localhost:8080/upload");
    http.addHeader("Content-Type", "image/jpeg");

    int httpResponseCode = http.POST(fb->buf, fb->len);

    if (httpResponseCode > 0)
    {
        if (httpResponseCode == 200)
        {
            client.publish(publish_topic, "ok");
            Serial.println("Pesan 'ok' dikirim ke MQTT!");
        }

        Serial.print("Image uploaded, response code: ");
        Serial.println(httpResponseCode);
    }
    else
    {
        Serial.print("Error uploading image: ");
        Serial.println(httpResponseCode);
    }

    http.end();
    esp_camera_fb_return(fb);
}

void sendMotionInfo()
{
    camera_fb_t *fb = esp_camera_fb_get();
    if (!fb)
    {
        Serial.println("Gagal mengambil gambar");
        return;
    }

    StaticJsonDocument<256> doc;
    doc["device_id"] = "ESP32-CAM";
    doc["timestamp"] = millis();
    doc["motion_detected"] = true;

    String jsonData;
    serializeJson(doc, jsonData);

    String boundary = "Boundary-" + String(millis(), HEX);

    HTTPClient http;
    if (!http.begin(api_url))
    {
        Serial.println("Gagal menghubungkan ke server");
        esp_camera_fb_return(fb);
        return;
    }

    http.addHeader("Content-Type", "multipart/form-data; boundary=" + boundary);

    std::string body = "--" + boundary + "\r\n"
                                         "Content-Disposition: form-data; name=\"metadata\"\r\n"
                                         "Content-Type: application/json\r\n\r\n" +
                       jsonData + "\r\n" +
                       "--" + boundary + "\r\n"
                                         "Content-Disposition: form-data; name=\"image\"; filename=\"motion.jpg\"\r\n"
                                         "Content-Type: image/jpeg\r\n\r\n";

    std::string footer = "\r\n--" + boundary + "--\r\n";

    size_t contentLength = body.length() + fb->len + footer.length();

    int httpResponseCode = http.POST(
        [&, sent = size_t(0)](uint8_t *buffer, size_t maxLen) mutable -> size_t
        {
            size_t toSend = 0;

            if (sent < body.length())
            {
                toSend = min(body.length() - sent, maxLen);
                memcpy(buffer, body.c_str() + sent, toSend);
            }
            else if (sent < (body.length() + fb->len))
            {
                size_t imgOffset = sent - body.length();
                toSend = min(fb->len - imgOffset, maxLen);
                memcpy(buffer, fb->buf + imgOffset, toSend);
            }
            else
            {
                size_t footerOffset = sent - (body.length() + fb->len);
                toSend = min(footer.length() - footerOffset, maxLen);
                memcpy(buffer, footer.c_str() + footerOffset, toSend);
            }

            sent += toSend;
            return toSend;
        },
        contentLength);

    if (httpResponseCode > 0)
    {
        Serial.print("Motion info + gambar terkirim. Response code: ");
        Serial.println(httpResponseCode);
    }
    else
    {
        Serial.print("Gagal mengirim: ");
        Serial.println(httpResponseCode);
    }

    http.end();
    esp_camera_fb_return(fb);
}

void setup()
{
    Serial.begin(115200);
    pinMode(pirPin, INPUT);
    pinMode(speakerPin, OUTPUT);
    setup_camera();
    setup_wifi();
    client.setServer(mqtt_server, mqtt_port);
    client.setCallback(callback);
}

void loop()
{
    if (!client.connected())
    {
        reconnect();
    }
    client.loop();

    // Baca sensor PIR
    int pirState = digitalRead(pirPin);

    if (pirState == HIGH)
    {
        Serial.println("Gerakan terdeteksi!");
        sendMotionInfo();
        client.publish(publish_topic, "Motion detected via PIR");
        delay(5000);
    }
}