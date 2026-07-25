from lights import LEDStrip
from mqtt import MQTTClient

def nightlight(strip: LEDStrip):
    strip.speed = 0.1
    strip.setColor(10, 10, 10)

def main():
    print("Starting...")
    strip = LEDStrip()
    strip.startAsyncUpdates()

    mqttClient = MQTTClient()
    mqttClient.register_action("on", lambda: strip.setColor(255, 255, 255))
    mqttClient.register_action("off", lambda: strip.setColor(0, 0, 0))
    mqttClient.register_action("middle_night", nightlight)
    mqttClient.start()

    print("Started!")

    while True:
        try:
            pass
        except KeyboardInterrupt:
            print("Stopping...")
            strip.stop()
            mqttClient.stop()
            break

    print("Exited")

if __name__=="__main__":
    main()