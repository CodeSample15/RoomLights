import time

# All light patterns MUST follow this format
class emptyPattern:
    def tick(strip): # Used as an "update" method
        pass
    def reset(): # Used to reset any color changes
        pass
    def at(x):
        return [0, 0, 0] # used to return the color for a specific pixel

class solidPattern:
    r, g, b = 50, 50, 50

    def tick(strip):
        pass
    def reset():
        pass
    def setColor(red, green, blue):
        solidPattern.r = red
        solidPattern.g = green
        solidPattern.b = blue
    def at(x):
        return [solidPattern.r, solidPattern.g, solidPattern.b]

class starsPattern:
    createStarTime = 0.055
    currentTime = time.time()
    starDecay = 0.8
    starChannelNormalValue = 2

    def tick(strip):
        if time.time() - currentTime > createStarTime:
            currentTime = time.time()
            randPos = random.randrange(0, strip.LED_COUNT)

            randC = random.randrange(0, 3)
            currColor = strip.getPixelColorRGB(randPos)

    def reset():
        pass

    def at(x):
        return [0,0,0]

patterns = [starsPattern,]