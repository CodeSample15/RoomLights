import paho.mqtt.client as mqtt
import json

# HAH enjoy the hardcoded IP nerds, you can only access it through my tailnet
MQTT_SERVER_IP = "100.93.66.64"
MQTT_SUBSCRIBE_TOPIC = "pico_commands"

class MQTTClient:
    def __init__(self):
            self.mqttc = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
            self.mqttc.on_connect = self.mqtt_on_connect
            self.mqttc.on_message = self.mqtt_on_message

    def mqtt_on_connect(self, client, userdata, flags, reason_code, properties):
        print(f"Connected with result code {reason_code}")
        client.subscribe(MQTT_SUBSCRIBE_TOPIC)

    def mqtt_on_message(self, client, userdata, msg):
        print(msg.topic+" "+str(msg.payload))
        body = json.loads(str(msg.payload))

        self.actions.get(body['type'], lambda:None)()

        self.actions = {}

    def register_action(self, messageType, action):
        self.actions[messageType] = action

    def start(self):
        self.mqttc.connect(MQTT_SERVER_IP, 1883, 60)
        self.mqttc.loop_start()

    def stop(self):
        self.mqttc.loop_stop()
        print("Stopped MQTT Client")