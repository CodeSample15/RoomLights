import paho.mqtt.client as mqtt

# HAH enjoy the hardcoded IP nerds, you can only access it through my tailnet
MQTT_SERVER_IP = "100.93.66.64"
MQTT_SUBSCRIBE_TOPIC = "pico_commands"

def mqtt_on_connect(client, userdata, flags, reason_code, properties):
    print(f"Connected with result code {reason_code}")
    client.subscribe(MQTT_SUBSCRIBE_TOPIC)

def mqtt_on_message(client, userdata, msg):
    print(msg.topic+" "+str(msg.payload))

def MQTTClientStart():
    mqttc = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
    mqttc.on_connect = mqtt_on_connect
    mqttc.on_message = mqtt_on_message

    mqttc.connect(MQTT_SERVER_IP, 1883, 60)

    mqttc.loop_forever()