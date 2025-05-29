#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <PubSubClient.h>
#include "esp_camera.h"
#include "soc/soc.h"
#include "soc/rtc_cntl_reg.h"

// PIN GPIO yang digunakan
#define PIR_PIN_1 13
#define PIR_PIN_2 15
#define PIR_PIN_3 14

#define BUZZER_PIN 33
#define LED_GREEN 32

// PIN Kamera
#define PWDN_GPIO_NUM -1
#define RESET_GPIO_NUM -1
#define XCLK_GPIO_NUM 21
#define SIOD_GPIO_NUM 26
#define SIOC_GPIO_NUM 27

#define Y9_GPIO_NUM 35
#define Y8_GPIO_NUM 34
#define Y7_GPIO_NUM 39
#define Y6_GPIO_NUM 36
#define Y5_GPIO_NUM 19
#define Y4_GPIO_NUM 18
#define Y3_GPIO_NUM 5
#define Y2_GPIO_NUM 4
#define VSYNC_GPIO_NUM 25
#define HREF_GPIO_NUM 23
#define PCLK_GPIO_NUM 22

// ID Sensor
const char *pirId1 = "SENPIR0001";
const char *pirId2 = "SENPIR0002";
const char *pirId3 = "SENPIR0003";

bool isPir1Active = false;
bool isPir2Active = false;
bool isPir3Active = false;

// Konfigurasi WIFI & Endpoint server
const char *ssid = "POCO F5";
const char *password = "RyuJin1029";
const int serverPort = 80;
String serverName = "farm.dihara.my.id";
// String serverName = "farm.test.dihara.my.id";

// Konfigurasi Broker
const char *mqttServer = "broker.hivemq.com";
const int mqttPort = 1883;
const char *mqttTopic = "bido_dihara/broker/farm-security";

String picture = "";

WiFiClient wifiClient;
PubSubClient client(wifiClient);

bool takePhoto = false;

void setup()
{
  WRITE_PERI_REG(RTC_CNTL_BROWN_OUT_REG, 0);
  Serial.begin(115200);
  WiFi.mode(WIFI_STA);

  startCamera();
  connectToWiFi();

  client.setServer(mqttServer, mqttPort);
  client.setCallback(mqttCallback);

  pinMode(PIR_PIN_1, INPUT);
  pinMode(PIR_PIN_2, INPUT);
  pinMode(PIR_PIN_3, INPUT);
  pinMode(BUZZER_PIN, OUTPUT);
  pinMode(LED_GREEN, OUTPUT);
}

void loop()
{
  if (!client.connected())
  {
    digitalWrite(LED_GREEN, LOW);
    reconnectMQTT();
  }
  else
  {
    digitalWrite(LED_GREEN, HIGH);
  }
  client.loop();

  bool motionDetected1 = isMotionDetected(PIR_PIN_1);
  bool motionDetected2 = isMotionDetected(PIR_PIN_2);
  bool motionDetected3 = isMotionDetected(PIR_PIN_3);

  if (motionDetected1 && isPir1Active)
  {
    Serial.println("Gerakan terdeteksi oleh PIR 1");
    buzzerActive();
    sendPhoto(pirId1, false);
  }
  if (motionDetected2 && isPir2Active)
  {
    Serial.println("Gerakan terdeteksi oleh PIR 2");
    buzzerActive();
    sendPhoto(pirId2, false);
  }
  if (motionDetected3 && isPir3Active)
  {
    Serial.println("Gerakan terdeteksi oleh PIR 3");
    buzzerActive();
    sendMotionEvent("", pirId3);
  }
  delay(1000);
}

void startCamera()
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
  config.pin_sccb_sda = SIOD_GPIO_NUM;
  config.pin_sccb_scl = SIOC_GPIO_NUM;
  config.pin_pwdn = PWDN_GPIO_NUM;
  config.pin_reset = RESET_GPIO_NUM;
  config.xclk_freq_hz = 20000000;
  config.pixel_format = PIXFORMAT_JPEG;

  config.frame_size = FRAMESIZE_QVGA;
  config.jpeg_quality = 10;
  config.fb_count = 1;

  esp_err_t err = esp_camera_init(&config);
  if (err != ESP_OK)
  {
    Serial.printf("Camera init failed with error 0x%x", err);
    delay(1000);
    ESP.restart();
  }
}

void connectToWiFi()
{
  WiFi.begin(ssid, password);
  Serial.print("Menghubungkan ke WiFi...");
  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nTersambung ke WiFi");
}

void reconnectMQTT()
{
  Serial.print("Menyambung ke MQTT...");
  while (!client.connected())
  {
    String clientId = "ESP32Cam-" + String(random(0xffff), HEX);
    if (client.connect(clientId.c_str()))
    {
      Serial.println("Tersambung ke MQTT");
      client.subscribe(mqttTopic);
    }
    else
    {
      Serial.print("Gagal. rc=");
      Serial.print(client.state());
    }
  }
}

String sendPhoto(const char *deviceId, bool isFromUser)
{
  String getAll;
  String getBody;

  camera_fb_t *fb = esp_camera_fb_get();
  if (!fb)
  {
    Serial.println("Camera capture failed");
    delay(1000);
    ESP.restart();
  }

  String filename = "esp32-" + String(millis()) + "-" + String(random(1000, 9999)) + ".jpg";
  picture = filename;

  if (wifiClient.connect(serverName.c_str(), serverPort))
  {
    String boundary = "SecurityFarm";
    String head = "--" + boundary + "\r\n"
                                    "Content-Disposition: form-data; name=\"imageFile\"; filename=\"" +
                  filename + "\"\r\n"
                             "Content-Type: image/jpeg\r\n\r\n";
    String tail = "\r\n--" + boundary + "--\r\n";

    uint32_t imageLen = fb->len;
    uint32_t extraLen = head.length() + tail.length();
    uint32_t totalLen = imageLen + extraLen;

    wifiClient.println("POST /api/upload/" + filename + " HTTP/1.1");
    wifiClient.println("Host: " + serverName);
    wifiClient.println("Content-Length: " + String(totalLen));
    wifiClient.println("Content-Type: multipart/form-data; boundary=" + boundary);
    wifiClient.println();

    wifiClient.print(head);

    uint8_t *fbBuf = fb->buf;
    size_t fbLen = fb->len;
    for (size_t n = 0; n < fbLen; n += 1024)
    {
      if (n + 1024 < fbLen)
      {
        wifiClient.write(fbBuf, 1024);
        fbBuf += 1024;
      }
      else if (fbLen % 1024 > 0)
      {
        size_t remainder = fbLen % 1024;
        wifiClient.write(fbBuf, remainder);
      }
    }

    wifiClient.print(tail);
    esp_camera_fb_return(fb);

    int timoutTimer = 10000;
    long startTimer = millis();
    bool state = false;

    while ((startTimer + timoutTimer) > millis())
    {
      Serial.print(".");
      delay(100);
      while (wifiClient.available())
      {
        char c = wifiClient.read();
        if (c == '\n')
        {
          if (getAll.length() == 0)
            state = true;
          getAll = "";
        }
        else if (c != '\r')
        {
          getAll += String(c);
        }
        if (state == true)
        {
          getBody += String(c);
        }
        startTimer = millis();
      }
      if (getBody.length() > 0)
        break;
    }

    Serial.println();
    wifiClient.stop();
    Serial.println(getBody);

    if (!isFromUser)
    {
      sendMotionEvent(filename, deviceId);
    } 
  }
  else
  {
    getBody = "Connection to " + serverName + " failed.";
    Serial.println(getBody);
  }

  return getBody;
}

String sendMotionEvent(String filename, String deviceId)
{
  String getAll;
  String getBody;

  if (filename == "")
  {
    filename = "motion-" + String(millis()) + "-" + String(random(1000, 9999));
  }

  String endpoint = "/api/motion-detected/" + filename;

  if (wifiClient.connect(serverName.c_str(), serverPort))
  {
    Serial.println("Connected to server!");

    String jsonBody = "{\"device_id\": \"" + deviceId + "\", \"motion_detected\": true}";
    int contentLength = jsonBody.length();

    wifiClient.println("POST " + endpoint + " HTTP/1.1");
    wifiClient.println("Host: " + serverName);
    wifiClient.println("Content-Type: application/json");
    wifiClient.println("Content-Length: " + String(contentLength));
    wifiClient.println();
    wifiClient.print(jsonBody);

    long startTimer = millis();
    int timeout = 5000;
    while (wifiClient.available() == 0 && millis() - startTimer < timeout)
    {
      delay(10);
    }

    while (wifiClient.available())
    {
      char c = wifiClient.read();
      getBody += c;
    }

    wifiClient.stop();
    Serial.println("Server response:");
    Serial.println(getBody);
  }
  else
  {
    Serial.println("Connection failed!");
  }

  return getBody;
}

void mqttCallback(char *topic, byte *payload, unsigned int length)
{
  Serial.print("Pesan Diterima [");
  Serial.print(topic);
  Serial.print("]: ");
  String message;
  for (int i = 0; i < length; i++)
  {
    message += (char)payload[i];
  }
  Serial.println(message);

  if (message == "TAKE_PHOTO")
  {
    sendPhoto(pirId1, true);
    reconnectMQTT();
    String message = "ok, pictureId: " + picture;
    client.publish(mqttTopic, message.c_str());

  }
  else if (message == "BUZZER_ON")
  {
    buzzerActive();
    client.publish(mqttTopic, "ok");
  }
  else if (message == "ALARM_ON")
  {
    alarmActive();
    client.publish(mqttTopic, "ok");
  }
  else if (message == "ALARM_OFF")
  {
    alarmNonActive();
    client.publish(mqttTopic, "ok");
  }
  else if (message.startsWith("DISABLE SENSOR"))
  {
    int idStart = message.indexOf("SENSOR ID: ") + strlen("SENSOR ID: ");
    int idEnd = message.indexOf(",", idStart);
    String sensorId = message.substring(idStart, idEnd);
    int isActiveStart = message.indexOf("IsActive: ") + strlen("IsActive: ");
    String isActiveStr = message.substring(isActiveStart);

    bool isActive = isActiveStr == "true";

    if (sensorId == pirId1)
    {
      isPir1Active = isActive;
    }
    else if (sensorId == pirId2)
    {
      isPir2Active = isActive;
    }
    else if (sensorId == pirId3)
    {
      isPir3Active = isActive;
    }
    client.publish(mqttTopic, "ok");
  }
}

void buzzerActive()
{
  digitalWrite(BUZZER_PIN, HIGH);
  delay(500);
  digitalWrite(BUZZER_PIN, LOW);
}

void alarmActive()
{
  digitalWrite(BUZZER_PIN, HIGH);
  client.publish(mqttTopic, "ok");
}

void alarmNonActive()
{
  digitalWrite(BUZZER_PIN, LOW);
  client.publish(mqttTopic, "ok");
}

bool isMotionDetected(int pin)
{
  return digitalRead(pin) == HIGH;
}
