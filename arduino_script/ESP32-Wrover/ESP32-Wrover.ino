#include "esp_camera.h"
#include <WiFi.h>
#include <HTTPClient.h>
#include <PubSubClient.h>

// #define FLASH_PIN 4

// PIN GPIO yang digunakan
#define PIR_PIN_1 13    
#define PIR_PIN_2 12     
#define PIR_PIN_3 14

#define BUZZER_PIN 33   
#define LED_GREEN 32

// PIN Kamera
#define PWDN_GPIO_NUM  -1
#define RESET_GPIO_NUM -1
#define XCLK_GPIO_NUM  21
#define SIOD_GPIO_NUM  26
#define SIOC_GPIO_NUM  27

#define Y9_GPIO_NUM    35
#define Y8_GPIO_NUM    34
#define Y7_GPIO_NUM    39
#define Y6_GPIO_NUM    36
#define Y5_GPIO_NUM    19
#define Y4_GPIO_NUM    18
#define Y3_GPIO_NUM    5
#define Y2_GPIO_NUM    4
#define VSYNC_GPIO_NUM 25
#define HREF_GPIO_NUM  23
#define PCLK_GPIO_NUM  22

// ID Sensor
const char* pirId1 = "SENPIR0001";
const char* pirId2 = "SENPIR0002";
const char* pirId3 = "SENPIR0003";

bool isPir1Active = true;
bool isPir2Active = true;
bool isPir3Active = true;

// Konfigurasi WIFI & Endpoint server
const char* ssid = "RYUJIN";
const char* password = "RyuJin1029";
const char* serverUrl = "https://farm.dihara.my.id/api/";

// Konfigurasi Broker
const char* mqttServer = "broker.hivemq.com";
const int mqttPort = 1883;
const char* mqttTopic = "bido_dihara/broker/farm-security";

WiFiClient espClient;
PubSubClient client(espClient);

bool takePhoto = false;

void setup() {
  Serial.begin(115200);
  startCamera();
  connectToWiFi();
  client.setServer(mqttServer, mqttPort);
  client.setCallback(mqttCallback);

  pinMode(PIR_PIN_1, INPUT);
  pinMode(PIR_PIN_2, INPUT);
  pinMode(PIR_PIN_3, INPUT);
  pinMode(BUZZER_PIN, OUTPUT);
  pinMode(LED_GREEN, OUTPUT);

  digitalWrite(LED_GREEN, HIGH);
}

void loop() {
  if (!client.connected()) {
      reconnectMQTT();
  }
  client.loop();

  bool motionDetected1 = isMotionDetected(PIR_PIN_1);
  bool motionDetected2 = isMotionDetected(PIR_PIN_2);
  bool motionDetected3 = isMotionDetected(PIR_PIN_3);

  if (motionDetected1 && isPir1Active) {
      Serial.println("Gerakan terdeteksi oleh PIR 1");
      buzzerActive();
      sendPhoto(pirId1, false);
  }
  if (motionDetected2 && isPir2Active) {
      Serial.println("Gerakan terdeteksi oleh PIR 2");
      buzzerActive();
      sendPhoto(pirId2, false);
  }
  if (motionDetected3 && isPir3Active) {
      Serial.println("Gerakan terdeteksi oleh PIR 3");
      buzzerActive();
      sendPhoto(pirId3, false);
  }
  delay(300);
}

void startCamera() {
  camera_config_t config;
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
  config.grab_mode = CAMERA_GRAB_WHEN_EMPTY;
  config.fb_location = CAMERA_FB_IN_PSRAM;
  config.frame_size = FRAMESIZE_SVGA;  
  config.jpeg_quality = 20;
  config.fb_count = 1;
  
  esp_err_t err = esp_camera_init(&config);
  if (err != ESP_OK) {
      Serial.printf("Kamera Gagal Inisialisasi Dengan Error 0x%x", err);
      return;
  }
}

void connectToWiFi() {
  WiFi.begin(ssid, password);
  Serial.print("Menghubungkan ke WiFi...");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nTersambung ke WiFi");
}

void reconnectMQTT() {
  Serial.print("Menyambung ke MQTT...");
  while (!client.connected()) {
    String clientId = "ESP32Cam-" + String(random(0xffff), HEX);
    if (client.connect(clientId.c_str())) {
      Serial.println("Tersambung ke MQTT");
      client.publish(mqttTopic, "Device online");
      client.subscribe(mqttTopic);
    } else {
      Serial.print("Gagal. rc=");
      Serial.print(client.state());
      Serial.println(" coba lagi dalam 2 detik.");
      delay(2000);
    }
  }
}

void sendPhoto(const char* deviceId, bool isFromUser) {
  delay(200); 
  camera_fb_t* fb2 = esp_camera_fb_get();
  if (!fb2) {
      Serial.println("Gagal Mengambil Gambar");
      return;
  }
  delay(300);
  esp_camera_fb_return(fb2);
  
  camera_fb_t* fb = esp_camera_fb_get();
  if (!fb) {
      Serial.println("Gagal Mengambil Gambar");
      return;
  }

  delay(100);
  randomSeed(millis());

  String pictureId = generateRandomString(15);
  if (WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    String fullUrlUpload = String(serverUrl) + "upload/" + pictureId;
    http.begin(fullUrlUpload);
    http.addHeader("Content-Type", "image/jpeg");
    int httpResponseCode = http.POST(fb->buf, fb->len);
    if (httpResponseCode > 0) {
      String payload = "ok pictureId: " + pictureId;
      Serial.printf("Gambar Terkirim! Kode Respon: %d\n", httpResponseCode);
      client.publish(mqttTopic, payload.c_str());
    } else {
      Serial.printf("Gagal Mengirim Gambar. Error: %s\n", http.errorToString(httpResponseCode).c_str());
    }
    http.end();
    delay(300);
    if (!isFromUser){
      HTTPClient http2;
      String fullUrlMotion = String(serverUrl) + "motion-detected/" + pictureId;
      http2.setTimeout(10000);
      http2.begin(fullUrlMotion);
      http2.addHeader("Content-Type", "application/json");
      String jsonBody = "{\"device_id\": \"" + String(deviceId) + "\", \"motion_detected\": true}";
      http2.POST(jsonBody);
      http2.end();
    }
  } else {
    Serial.println("WiFi Tidak Terhubung.");
  }
  esp_camera_fb_return(fb);
}

void mqttCallback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Pesan Diterima [");
  Serial.print(topic);
  Serial.print("]: ");
  String message;
  for (int i = 0; i < length; i++) {
      message += (char)payload[i];
  }
  Serial.println(message);

  if (message == "TAKE_PHOTO") {
    sendPhoto(pirId1, true); 
  } else if (message == "BUZZER_ON") {
    buzzerActive();
    client.publish(mqttTopic, "ok");
  } else if (message == "ALARM_ON") {
    alarmActive();
    client.publish(mqttTopic, "ok");
  } else if (message == "ALARM_OFF") {
    alarmNonActive();
    client.publish(mqttTopic, "ok");
  } else if (message.startsWith("DISABLE SENSOR")){
    int idStart = message.indexOf("SENSOR ID: ") + strlen("SENSOR ID: ");
    int idEnd = message.indexOf(",", idStart);
    String sensorId = message.substring(idStart, idEnd);
    int isActiveStart = message.indexOf("IsActive: ") + strlen("IsActive: ");
    String isActiveStr = message.substring(isActiveStart);

    bool isActive = isActiveStr == "true";

    if (sensorId == pirId1){
      isPir1Active = isActive;
    } else if (sensorId == pirId2){
      isPir2Active = isActive;
    } else if (sensorId == pirId3){
      isPir3Active = isActive;
    }
    client.publish(mqttTopic, "ok");
  }
}

String generateRandomString(int length) {
  const char charset[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  String randomString;
  for (int i = 0; i < length; i++) {
    randomString += charset[random(sizeof(charset) - 2)];
  }
  return randomString;
}

void buzzerActive(){
  digitalWrite(BUZZER_PIN, HIGH);
  delay(500);
  digitalWrite(BUZZER_PIN, LOW);
}

void alarmActive(){
  digitalWrite(BUZZER_PIN, HIGH);
  client.publish(mqttTopic, "ok");
}

void alarmNonActive(){
  digitalWrite(BUZZER_PIN, LOW);
  client.publish(mqttTopic, "ok");
}

bool isMotionDetected(int pin) {
  return digitalRead(pin) == HIGH;
}

