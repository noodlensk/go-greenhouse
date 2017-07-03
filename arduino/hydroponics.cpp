#include <Wire.h>
#include "RTClib.h"
#include <ArduinoJson.h>
#include <HTU21D.h>

RTC_DS3231 rtc;
HTU21D myHTU21D(HTU21D_RES_RH12_TEMP14);


typedef struct {
  String Name;
  int StartTime;
  int EndTime;
  bool IsOn;
  int Pin;
} Rele;

Rele ReleList[3] = {
    {
      "Rele1",
      9,
      21,
      false,
      12
    },
    {
      "Rele2",
      6,
      24,
      false,
      13
    },
    {
      "Rele3",
      6,
      24,
      false,
      6
    }
};

bool ManualHandling = false;
DateTime now;
int Hours;
String serialResponse = "";

void setup()
{
  Serial.begin(9600);
  rtc.begin();
  myHTU21D.begin();
  for (int i = 0; i < sizeof(ReleList)/sizeof(ReleList[0]) ; i++) {
    if (ReleList[i].Pin > 0) {
      pinMode(ReleList[i].Pin, OUTPUT);
    }
  }
  if (rtc.lostPower()) {
    Serial.println("RTC lost power, lets set the time!");
    rtc.adjust(DateTime(F(__DATE__), F(__TIME__)));
    // This line sets the RTC with an explicit date & time, for example to set
    // January 21, 2014 at 3am you would call:
    // rtc.adjust(DateTime(2014, 1, 21, 3, 0, 0));
  }
}

void loop()
{
    now = rtc.now();
    Hours = now.hour();
    manageRele();
    if ( Serial.available()) {
      serialResponse = Serial.readStringUntil('\r\n');
      // Convert from String Object to String.
      char buf[100];
      serialResponse.toCharArray(buf, sizeof(buf));
      char *p = buf;
      char *str;
      String data[20];
      int i = 0;
      while ((str = strtok_r(p, ";", &p)) != NULL) {// delimiter is the semicolon
        data[i++] = str;
      }
      String action = data[0];
      if(action == "get_data") {
        Serial.println("OK");
        printData();
      } else if (action == "set_manual") {
          ManualHandling = (bool) data[1].toInt();
          Serial.println("OK");
          Serial.print("ManualHandling  = ");Serial.println(ManualHandling);
      } else if (action == "switch") {
          String releName = data[1];
          bool state = (bool) data[2].toInt();
          for (int i = 0; i < sizeof(ReleList)/sizeof(ReleList[0]) ; i++) {
            if (ReleList[i].Name == releName) {
                ReleList[i].IsOn = state;
                Serial.println("OK");
            }
          }
      }
    }
    //printData();
    //delay(1000);
}

void manageRele() {
    for (int i = 0; i < sizeof(ReleList)/sizeof(ReleList[0]) ; i++) {
    if (!ManualHandling ) {
        if(Hours >= ReleList[i].StartTime && Hours <= (ReleList[i].EndTime - 1)) {
            ReleList[i].IsOn = true;
        } else {
            ReleList[i].IsOn = false;
        }
    }
    if (ReleList[i].IsOn) {
        digitalWrite(ReleList[i].Pin, HIGH);
    } else {
        digitalWrite(ReleList[i].Pin, LOW);
    }
  }
}
void printData() {
    char tempt[10], humd[10];
    String delimiter = ";";

    Serial.print(now.unixtime());

    Serial.print(delimiter);
    dtostrf(myHTU21D.readTemperature(), 2, 2, tempt);
    dtostrf(myHTU21D.readHumidity(), 2, 2, humd);
    Serial.print(tempt);
    Serial.print(delimiter);
    Serial.print(humd);
    Serial.print(delimiter);

    Serial.print(ManualHandling);
    Serial.print(delimiter);
    for (int i = 0; i < sizeof(ReleList)/sizeof(ReleList[0]) ; i++) {
      Serial.print(ReleList[i].Name);Serial.print("#");
      Serial.print(ReleList[i].StartTime);Serial.print("#");
      Serial.print(ReleList[i].EndTime);Serial.print("#");
      Serial.print(ReleList[i].IsOn);Serial.print("#");
      Serial.print(ReleList[i].Pin);
      Serial.print(delimiter);
    }
    Serial.println("");
}