from rpi_ws281x import *
import asyncio

#LED mode enum values
LED_MODE_SOLID = 0
LED_MODE_GRADIENT = 1

class LEDStrip:
    def __init__(self, led_count=300, led_pin=18, freq=800000, dma=10, brightness=255, invert=False, channel=0):
        self.LED_COUNT = led_count
        self.LED_PIN = led_pin
        self.LED_FREQ_HZ = freq
        self.LED_DMA = dma
        self.LED_BRIGHTNESS = brightness
        self.LED_INVERT = invert
        self.LED_CHANNEL = channel

        #speed is how fast the lights will adjust to the target color
        #set this to zero for instant change
        self.speed = 0
        self.mode = LED_MODE_GRADIENT

        self.running = False
        self.currentColor = (0,0,0)
        self.targetColor = (0,0,0)
        self.updateTask = None

        self.strip = Adafruit_NeoPixel(self.LED_COUNT, self.LED_PIN, self.LED_FREQ_HZ, self.LED_DMA, self.LED_INVERT, self.LED_BRIGHTNESS, self.LED_CHANNEL)
        self.strip.begin()

        #start off as off
        for p in range(led_count):
            self.strip.setPixelColor(p, Color(0,0,0))

    def startAsyncUpdates(self):
        asyncio.create_task(self._asyncUpdates())
        print("LED strip started!")

    def stop(self):
        self.running = False
        print("Stopped LED strip")

    def setColor(self, r, g, b):
        self.targetColor = (r, g, b)

    async def _asyncUpdates(self):
        self.running = True
        while self.running:
            # calculate color deltas
            if self.speed != 0:
                for i in range(3):
                    self.currentColor[i] = float(self.targetColor[i] - self.currentColor[i]) * self.speed
            else:
                self.currentColor = self.targetColor

            for p in range(self.LED_COUNT):
                self.strip.setPixelColor(p, Color(self.currentColor[0], self.currentColor[1], self.currentColor[2]))
            await asyncio.sleep(0.1)