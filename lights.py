from rpi_ws281x import *
from patterns import emptyPattern
import threading
import time

class LEDStrip:
    def __init__(self, led_count=300, led_pin=12, freq=800000, dma=10, brightness=255, invert=False, channel=0):
        self.LED_COUNT = led_count
        self.LED_PIN = led_pin
        self.LED_FREQ_HZ = freq
        self.LED_DMA = dma
        self.LED_BRIGHTNESS = brightness
        self.LED_INVERT = invert
        self.LED_CHANNEL = channel

        # Speed is how fast the lights will adjust to the target color.
        # Set this to zero for instant change.
        self.speed = 0

        self.thread = threading.Thread(target=self._asyncUpdates)
        self.running = False

        self.targetPattern = emptyPattern

        self.strip = Adafruit_NeoPixel(self.LED_COUNT, self.LED_PIN, self.LED_FREQ_HZ, self.LED_DMA, self.LED_INVERT, self.LED_BRIGHTNESS, self.LED_CHANNEL)
        self.strip.begin()

        #start off as off
        for p in range(led_count):
            self.strip.setPixelColor(p, Color(0,0,0))
        self.strip.show()

    def startAsyncUpdates(self):
        self.thread.start()
        print("LED strip started!")

    def stop(self):
        self.running = False
        self.thread.join()
        print("Stopped LED strip")

    def setTargetPattern(self, pattern, reset=False, speed=-1):
        print("LED strip: Changing target pattern")
        if speed >= 0:
            self.speed = speed
        if reset:
            pattern.reset()
        self.targetPattern = pattern

    def _asyncUpdates(self):
        self.running = True
        while self.running:
            for p in range(self.LED_COUNT):
                r, g, b = self.targetPattern.at(p)

                if self.speed != 0:
                    # calculate color change
                    currColor = self.strip.getPixelColorRGB(p)
                    currColor = [currColor.r, currColor.g, currColor.b]
                    targetColor = [r, g, b]

                    for i in range(3):
                        if targetColor[i] - currColor[i] <= 3:
                            currColor[i] = targetColor[i]
                        else:
                            currColor[i] += float(targetColor[i] - currColor[i]) * self.speed

                    self.strip.setPixelColor(p, Color(int(currColor[0]), int(currColor[1]), int(currColor[2])))
                else:
                    self.strip.setPixelColor(p, Color(int(r), int(g), int(b)))

            self.targetPattern.tick(self.strip)
            self.strip.show()
            time.sleep(0.01)

if __name__=="__main__":
    test=LEDStrip()