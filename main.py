from lights import LEDStrip
from mqtt import MQTTClient

import patterns

def setSolidColor(strip, r, g, b, speed=-1):
    patterns.solidPattern.setColor(r, g, b)
    strip.setTargetPattern(patterns.solidPattern, speed=speed)

def main():
    print("Starting...")
    strip = LEDStrip()
    strip.startAsyncUpdates()

    mqttClient = MQTTClient()
    mqttClient.register_action("on", lambda: setSolidColor(strip, 255, 255, 255, 0.1))
    mqttClient.register_action("off", lambda: strip.setTargetPattern(patterns.emptyPattern, speed=0.5))
    mqttClient.register_action("middle_night", lambda: setSolidColor(strip, 120, 64, 0, 0.01)) # night light
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