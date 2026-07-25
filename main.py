from lights import LEDStrip
from mqtt import MQTTClient

def main():
    print("Starting...")
    strip = LEDStrip()
    strip.startAsyncUpdates()

    mqttClient = MQTTClient()
    mqttClient.register_action("on", lambda: strip.setColor(255, 255, 255))
    mqttClient.register_action("on", lambda: strip.setColor(0, 0, 0))
    mqttClient.start()

    while True:
        try:
            pass
        except KeyboardInterrupt:
            strip.stop()
            mqttClient.stop()
            break

    print("Exited")

if __name__=="__main__":
    main()