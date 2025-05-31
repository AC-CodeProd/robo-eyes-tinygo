package roboeyestinygo

import (
	"image/color"
	"math/rand"
	"time"
)

// DeviceInterface defines required methods for device control
type DeviceInterface interface {
	ClearBuffer()                      // Clear display buffer
	Display() error                    // Update physical display
	SetPixel(x, y int16, c color.RGBA) // Set individual pixel
	Size() (width, height int16)       // Get display dimensions
}

// Mood constants
type Mood byte

const (
	MoodDefault Mood = iota
	MoodTired
	MoodAngry
	MoodHappy
)

// Eye direction constants
type Direction byte

const (
	DirCenter Direction = iota
	DirN
	DirNE
	DirE
	DirSE
	DirS
	DirSW
	DirW
	DirNW
)

// RoboEyes represents the robot eyes controller
type RoboEyes struct {
	device    DeviceInterface
	startTime time.Time
	eyesColor color.RGBA
	bgColor   color.RGBA

	// Display parameters
	screenWidth   int16
	screenHeight  int16
	frameInterval uint32
	fpsTimer      uint32

	// Eye states
	tired   bool
	angry   bool
	happy   bool
	curious bool

	cyclops   bool
	eyeL_open bool
	eyeR_open bool

	// Left eye geometry
	eyeLwidthDefault        int16
	eyeLheightDefault       int16
	eyeLwidthCurrent        int16
	eyeLheightCurrent       int16
	eyeLwidthNext           int16
	eyeLheightNext          int16
	eyeLheightOffset        int16
	eyeLborderRadiusDefault byte
	eyeLborderRadiusCurrent byte
	eyeLborderRadiusNext    byte
	eyeLxDefault            int16
	eyeLyDefault            int16
	eyeLx                   int16
	eyeLy                   int16
	eyeLxNext               int16
	eyeLyNext               int16

	// Right eye geometry
	eyeRwidthDefault        int16
	eyeRheightDefault       int16
	eyeRwidthCurrent        int16
	eyeRheightCurrent       int16
	eyeRwidthNext           int16
	eyeRheightNext          int16
	eyeRheightOffset        int16
	eyeRborderRadiusDefault byte
	eyeRborderRadiusCurrent byte
	eyeRborderRadiusNext    byte
	eyeRxDefault            int16
	eyeRyDefault            int16
	eyeRx                   int16
	eyeRy                   int16
	eyeRxNext               int16
	eyeRyNext               int16

	// Common parameters
	spaceBetweenDefault          int16
	spaceBetweenCurrent          int16
	spaceBetweenNext             int16
	eyelidsHeightMax             int16
	eyelidsTiredHeight           int16
	eyelidsTiredHeightNext       int16
	eyelidsAngryHeight           int16
	eyelidsAngryHeightNext       int16
	eyelidsHappyBottomOffset     int16
	eyelidsHappyBottomOffsetNext int16
	eyelidsHappyBottomOffsetMax  int16

	// Animation states
	hFlicker                  bool
	hFlickerAlternate         bool
	hFlickerAmplitude         int16
	vFlicker                  bool
	vFlickerAlternate         bool
	vFlickerAmplitude         int16
	autoblinker               bool
	blinkInterval             uint32
	blinkIntervalVariation    uint32
	blinktimer                uint32
	idle                      bool
	idleInterval              uint32
	idleIntervalVariation     uint32
	idleAnimationTimer        uint32
	confused                  bool
	confusedAnimationTimer    uint32
	confusedAnimationDuration uint32
	confusedToggle            bool
	laugh                     bool
	laughAnimationTimer       uint32
	laughAnimationDuration    uint32
	laughToggle               bool
}

func (r *RoboEyes) setDefault(screenWidth, screenHeight int16) {
	r.startTime = time.Now()

	r.eyesColor = color.RGBA{255, 255, 255, 255}
	r.bgColor = color.RGBA{0, 0, 0, 255}

	// For general setup - screen size and max. frame rate
	r.screenWidth = screenWidth   // OLED display width, in pixels
	r.screenHeight = screenHeight // OLED display height, in pixels
	r.frameInterval = 20          // default value for 50 frames per second (1000/50 = 20 milliseconds)
	r.fpsTimer = 0                // for timing the frames per second

	// For controlling mood types and expressions
	r.tired = false
	r.angry = false
	r.happy = false
	r.curious = false   // if true, draw the outer eye larger when looking left or right
	r.cyclops = false   // if true, draw only one eye
	r.eyeL_open = false // left eye opened or closed?
	r.eyeR_open = false // right eye opened or closed?

	//*********************************************************************************************
	//  Eyes Geometry
	//*********************************************************************************************

	// EYE LEFT - size and border radius
	r.eyeLwidthDefault = 36
	r.eyeLheightDefault = 36
	r.eyeLwidthCurrent = r.eyeLwidthDefault
	r.eyeLheightCurrent = 1 // start with closed eye, otherwise set to eyeLheightDefault
	r.eyeLwidthNext = r.eyeLwidthDefault
	r.eyeLheightNext = r.eyeLheightDefault
	r.eyeLheightOffset = 0
	// Border Radius
	r.eyeLborderRadiusDefault = 8
	r.eyeLborderRadiusCurrent = r.eyeLborderRadiusDefault
	r.eyeLborderRadiusNext = r.eyeLborderRadiusDefault

	// EYE RIGHT - size and border radius
	r.eyeRwidthDefault = r.eyeLwidthDefault
	r.eyeRheightDefault = r.eyeLheightDefault
	r.eyeRwidthCurrent = r.eyeRwidthDefault
	r.eyeRheightCurrent = 1 // start with closed eye, otherwise set to eyeRheightDefault
	r.eyeRwidthNext = r.eyeRwidthDefault
	r.eyeRheightNext = r.eyeRheightDefault
	r.eyeRheightOffset = 0
	// Border Radius
	r.eyeRborderRadiusDefault = 8
	r.eyeRborderRadiusCurrent = r.eyeRborderRadiusDefault
	r.eyeRborderRadiusNext = r.eyeRborderRadiusDefault

	// EYE LEFT - Coordinates
	r.eyeLxDefault = ((r.screenWidth) - (r.eyeLwidthDefault + r.spaceBetweenDefault + r.eyeRwidthDefault)) / 2
	r.eyeLyDefault = ((r.screenHeight - r.eyeLheightDefault) / 2)
	r.eyeLx = r.eyeLxDefault
	r.eyeLy = r.eyeLyDefault
	r.eyeLxNext = r.eyeLx
	r.eyeLyNext = r.eyeLy

	// EYE RIGHT - Coordinates
	r.eyeRxDefault = r.eyeLx + r.eyeLwidthCurrent + r.spaceBetweenDefault
	r.eyeRyDefault = r.eyeLy
	r.eyeRx = r.eyeRxDefault
	r.eyeRy = r.eyeRyDefault
	r.eyeRxNext = r.eyeRx
	r.eyeRyNext = r.eyeRy

	// BOTH EYES
	// Eyelid top size
	r.eyelidsHeightMax = r.eyeLheightDefault / 2 // top eyelids max height
	r.eyelidsTiredHeight = 0
	r.eyelidsTiredHeightNext = r.eyelidsTiredHeight
	r.eyelidsAngryHeight = 0
	r.eyelidsAngryHeightNext = r.eyelidsAngryHeight
	// Bottom happy eyelids offset
	r.eyelidsHappyBottomOffsetMax = (r.eyeLheightDefault / 2) + 3
	r.eyelidsHappyBottomOffset = 0
	r.eyelidsHappyBottomOffsetNext = 0
	// Space between eyes
	r.spaceBetweenDefault = 10
	r.spaceBetweenCurrent = r.spaceBetweenDefault
	r.spaceBetweenNext = 10

	//*********************************************************************************************
	//  Macro Animations
	//*********************************************************************************************

	// Animation - horizontal flicker/shiver
	r.hFlicker = false
	r.hFlickerAlternate = false
	r.hFlickerAmplitude = 2

	// Animation - vertical flicker/shiver
	r.vFlicker = false
	r.vFlickerAlternate = false
	r.vFlickerAmplitude = 10

	// Animation - auto blinking
	r.autoblinker = false               // activate auto blink animation
	r.blinkInterval = 1 * 1000          // basic interval between each blink in full seconds
	r.blinkIntervalVariation = 4 * 1000 // interval variaton range in full seconds, random number inside of given range will be add to the basic blinkInterval, set to 0 for no variation
	r.blinktimer = 0                    // for organising eyeblink timing

	// Animation - idle mode: eyes looking in random directions
	r.idle = false
	r.idleInterval = 1 * 1000          // basic interval between each eye repositioning in full seconds
	r.idleIntervalVariation = 3 * 1000 // interval variaton range in full seconds, random number inside of given range will be add to the basic idleInterval, set to 0 for no variation
	r.idleAnimationTimer = 0           // for organising eyeblink timing

	// Animation - eyes confused: eyes shaking left and right
	r.confused = false
	r.confusedAnimationTimer = 0
	r.confusedAnimationDuration = 500
	r.confusedToggle = true

	// Animation - eyes laughing: eyes shaking up and down
	r.laugh = false
	r.laughAnimationTimer = 0
	r.laughAnimationDuration = 500
	r.laughToggle = true

}

// Begin initializes the RoboEyes controller
func (r *RoboEyes) Begin(device DeviceInterface, width, height int16, frameRate uint32) {
	r.device = device
	r.setDefault(width, height)
	r.SetFramerate(frameRate)
}

// millis returns milliseconds since initialization
func (r *RoboEyes) millis() uint32 {
	return uint32(time.Since(r.startTime).Milliseconds())
}

// Update handles timed updates and animations
// Should be called in the main loop
func (r *RoboEyes) Update() {
	currentTime := r.millis()

	// Limit updates to defined frame rate
	if currentTime-r.fpsTimer >= r.frameInterval {
		r.DrawEyes()
		r.fpsTimer = currentTime
	}
}

// SetFramerate sets the maximum frame rate
func (r *RoboEyes) SetFramerate(fps uint32) {
	if fps > 0 {
		r.frameInterval = 1000 / fps
	}
}

// SetSize sets default eye dimensions
func (r *RoboEyes) SetSize(left, right int16) {
	r.eyeLwidthNext = left
	r.eyeRwidthNext = right
	r.eyeLwidthDefault = left
	r.eyeRwidthDefault = right
}

// SetBorderRadius sets eye corner rounding
func (r *RoboEyes) SetBorderRadius(left, right byte) {
	r.eyeLborderRadiusNext = left
	r.eyeRborderRadiusNext = right
	r.eyeLborderRadiusDefault = left
	r.eyeRborderRadiusDefault = right
}

// SetSpaceBetween sets distance between eyes
func (r *RoboEyes) SetSpaceBetween(space int16) {
	r.spaceBetweenNext = space
	r.spaceBetweenDefault = space
}

// SetMood configures eye expression
func (r *RoboEyes) SetMood(mood Mood) {
	r.tired, r.angry, r.happy = false, false, false
	switch mood {
	case MoodTired:
		r.tired = true
	case MoodAngry:
		r.angry = true
	case MoodHappy:
		r.happy = true
	default:

	}
}

// SetDirection moves eyes to predefined location
func (r *RoboEyes) SetDirection(direction Direction) {
	maxX := r.GetScreenConstraintX()
	maxY := r.GetScreenConstraintY()
	switch direction {
	case DirN:
		// North, top center
		r.eyeLxNext = maxX / 2
		r.eyeLyNext = 0
	case DirNE:
		// North-east, top right
		r.eyeLxNext = maxX
		r.eyeLyNext = 0
	case DirE:
		// East, middle right
		r.eyeLxNext = maxX
		r.eyeLyNext = maxY / 2
	case DirSE:
		// South-east, bottom right
		r.eyeLxNext = maxX
		r.eyeLyNext = maxY
	case DirS:
		// South, bottom center
		r.eyeLxNext = maxX / 2
		r.eyeLyNext = maxY
	case DirSW:
		// South-west, bottom left
		r.eyeLxNext = 0
		r.eyeLyNext = maxY
	case DirW:
		// West, middle left
		r.eyeLxNext = 0
		r.eyeLyNext = maxY / 2
	case DirNW:
		// North-west, top left
		r.eyeLxNext = 0
		r.eyeLyNext = 0
	default:
		// Middle center
		r.eyeLxNext = maxX / 2
		r.eyeLyNext = maxY / 2
	}
}

// GetScreenConstraintX returns maximum X position for left eye
func (r *RoboEyes) GetScreenConstraintX() int16 {
	return r.screenWidth - r.eyeLwidthCurrent - r.spaceBetweenCurrent - r.eyeRwidthCurrent
}

// GetScreenConstraintY returns maximum Y position for eyes
func (r *RoboEyes) GetScreenConstraintY() int16 {
	return r.screenHeight - r.eyeLheightDefault
}

// SetAutoBlinker configures automatic blinking
func (r *RoboEyes) SetAutoBlinkerWithInterval(active bool, interval, variation uint32) {
	r.SetAutoBlinker(active)
	r.blinkInterval = interval * 1000
	r.blinkIntervalVariation = variation * 1000
}

func (r *RoboEyes) SetAutoBlinker(active bool) {
	r.autoblinker = active
}

// SetIdleMode configures random eye movements
func (r *RoboEyes) SetIdleModeWithInterval(active bool, interval, variation uint32) {
	r.SetIdleMode(active)
	r.idleInterval = interval * 1000
	r.idleIntervalVariation = variation * 1000
}

func (r *RoboEyes) SetIdleMode(active bool) {
	r.idle = active
}

// SetCuriosity enables/disables curious gaze effect
func (r *RoboEyes) SetCuriosity(active bool) {
	r.curious = active
}

// SetCyclops enables/disables single eye mode
func (r *RoboEyes) SetCyclops(active bool) {
	r.cyclops = active
}

// SetHFlicker configures horizontal flicker effect
func (r *RoboEyes) SetHFlicker(active bool, amplitude int16) {
	r.hFlicker = active
	r.hFlickerAmplitude = amplitude
}

// SetVFlicker configures vertical flicker effect
func (r *RoboEyes) SetVFlicker(active bool, amplitude int16) {
	r.vFlicker = active
	r.vFlickerAmplitude = amplitude
}

// Close closes both eyes
func (r *RoboEyes) Close() {
	r.eyeLheightNext = 1
	r.eyeRheightNext = 1
	r.eyeL_open = false
	r.eyeR_open = false
}

// Open opens both eyes
func (r *RoboEyes) Open() {
	r.eyeL_open = true
	r.eyeR_open = true
}

// Blink performs a blink animation
func (r *RoboEyes) Blink() {
	r.Close()
	r.Open()
}

// CloseEyes closes specified eyes
func (r *RoboEyes) CloseEyes(left, right bool) {
	if left {
		r.eyeLheightNext = 1
		r.eyeL_open = false
	}
	if right {
		r.eyeRheightNext = 1
		r.eyeR_open = false
	}
}

// OpenEyes opens specified eyes
func (r *RoboEyes) OpenEyes(left, right bool) {
	if left {
		r.eyeL_open = true
	}
	if right {
		r.eyeR_open = true
	}
}

// BlinkEyes blinks specified eyes
func (r *RoboEyes) BlinkEyes(left, right bool) {
	r.CloseEyes(left, right)
	r.OpenEyes(left, right)
}

// AnimConfused triggers confused animation
func (r *RoboEyes) AnimConfused() {
	r.confused = true
}

// AnimLaugh triggers laugh animation
func (r *RoboEyes) AnimLaugh() {
	r.laugh = true
}

// DrawEyes renders the eyes on the display
func (r *RoboEyes) DrawEyes() {
	currentTime := r.millis()

	// Calculate eye geometry with smoothing
	r.calculateGeometry()

	// Handle automatic animations
	r.handleAnimations(currentTime)

	// Prepare display
	r.device.ClearBuffer()

	// Draw eyes
	r.drawEyeShapes()

	// Draw eyelids based on mood
	r.drawEyelids()

	// Update physical display
	r.device.Display()
}

// calculateGeometry updates eye positions and sizes with smoothing
func (r *RoboEyes) calculateGeometry() {
	// Apply curious effect (enlarge outer eye)
	if r.curious {
		if r.eyeLxNext <= 10 {
			r.eyeLheightOffset = 8
		} else if r.eyeLxNext >= (r.GetScreenConstraintX()-10) && r.cyclops {
			r.eyeLheightOffset = 8
		} else {
			r.eyeLheightOffset = 0
		}
		if r.eyeRxNext >= r.screenWidth-r.eyeRwidthCurrent-10 {
			r.eyeRheightOffset = 8
		} else {
			r.eyeRheightOffset = 0
		}
	} else {
		r.eyeLheightOffset = 0
		r.eyeRheightOffset = 0
	}

	// Left eye height with smoothing
	r.eyeLheightCurrent = (r.eyeLheightCurrent + r.eyeLheightNext + r.eyeLheightOffset) / 2
	r.eyeLy += (r.eyeLheightDefault - r.eyeLheightCurrent) / 2
	r.eyeLy -= r.eyeLheightOffset / 2

	// Right eye height with smoothing
	r.eyeRheightCurrent = (r.eyeRheightCurrent + r.eyeRheightNext + r.eyeRheightOffset) / 2
	r.eyeRy += (r.eyeRheightDefault - r.eyeRheightCurrent) / 2
	r.eyeRy -= r.eyeRheightOffset / 2

	// Reopen eyes after closing
	if r.eyeL_open && r.eyeLheightCurrent <= 1+r.eyeLheightOffset {
		r.eyeLheightNext = r.eyeLheightDefault
	}
	if r.eyeR_open && r.eyeRheightCurrent <= 1+r.eyeRheightOffset {
		r.eyeRheightNext = r.eyeRheightDefault
	}

	// Width smoothing
	r.eyeLwidthCurrent = (r.eyeLwidthCurrent + r.eyeLwidthNext) / 2
	r.eyeRwidthCurrent = (r.eyeRwidthCurrent + r.eyeRwidthNext) / 2

	// Space between eyes smoothing
	r.spaceBetweenCurrent = (r.spaceBetweenCurrent + r.spaceBetweenNext) / 2

	// Position smoothing
	r.eyeLx = (r.eyeLx + r.eyeLxNext) / 2
	r.eyeLy = (r.eyeLy + r.eyeLyNext) / 2
	r.eyeRxNext = r.eyeLxNext + r.eyeLwidthCurrent + r.spaceBetweenCurrent
	r.eyeRyNext = r.eyeLyNext
	r.eyeRx = (r.eyeRx + r.eyeRxNext) / 2
	r.eyeRy = (r.eyeRy + r.eyeRyNext) / 2

	// Border radius smoothing
	r.eyeLborderRadiusCurrent = (r.eyeLborderRadiusCurrent + r.eyeLborderRadiusNext) / 2
	r.eyeRborderRadiusCurrent = (r.eyeRborderRadiusCurrent + r.eyeRborderRadiusNext) / 2
}

// handleAnimations processes automatic and triggered animations
func (r *RoboEyes) handleAnimations(currentTime uint32) {
	// Automatic blinking
	if r.autoblinker && currentTime >= r.blinktimer {
		r.Blink()
		r.blinktimer = currentTime + r.blinkInterval + uint32(rand.Intn(int(r.blinkIntervalVariation)))
	}

	// Laugh animation (vertical shaking)
	if r.laugh {
		if r.laughToggle {
			r.SetVFlicker(true, 5)
			r.laughAnimationTimer = currentTime
			r.laughToggle = false
		} else if currentTime >= r.laughAnimationTimer+r.laughAnimationDuration {
			r.SetVFlicker(false, 0)
			r.laughToggle = true
			r.laugh = false
		}
	}

	// Confused animation (horizontal shaking)
	if r.confused {
		if r.confusedToggle {
			r.SetHFlicker(true, 20)
			r.confusedAnimationTimer = currentTime
			r.confusedToggle = false
		} else if currentTime >= r.confusedAnimationTimer+r.confusedAnimationDuration {
			r.SetHFlicker(false, 0)
			r.confusedToggle = true
			r.confused = false
		}
	}

	// Idle mode (random eye movements)
	if r.idle && currentTime >= r.idleAnimationTimer {
		r.eyeLxNext = int16(rand.Intn(int(r.GetScreenConstraintX())))
		r.eyeLyNext = int16(rand.Intn(int(r.GetScreenConstraintY())))
		r.idleAnimationTimer = currentTime + r.idleInterval + uint32(rand.Intn(int(r.idleIntervalVariation)))
	}

	// Apply horizontal flicker
	if r.hFlicker {
		if r.hFlickerAlternate {
			r.eyeLx += r.hFlickerAmplitude
			r.eyeRx += r.hFlickerAmplitude
		} else {
			r.eyeLx -= r.hFlickerAmplitude
			r.eyeRx -= r.hFlickerAmplitude
		}
		r.hFlickerAlternate = !r.hFlickerAlternate
	}

	// Apply vertical flicker
	if r.vFlicker {
		if r.vFlickerAlternate {
			r.eyeLy += r.vFlickerAmplitude
			r.eyeRy += r.vFlickerAmplitude
		} else {
			r.eyeLy -= r.vFlickerAmplitude
			r.eyeRy -= r.vFlickerAmplitude
		}
		r.vFlickerAlternate = !r.vFlickerAlternate
	}

	// Cyclops mode (hide right eye)
	if r.cyclops {
		r.eyeRwidthCurrent = 0
		r.eyeRheightCurrent = 0
		r.spaceBetweenCurrent = 0
	}
}

// drawEyeShapes renders the main eye shapes
func (r *RoboEyes) drawEyeShapes() {
	// Convert border radius to int16 for drawing
	borderL := int16(r.eyeLborderRadiusCurrent)
	borderR := int16(r.eyeRborderRadiusCurrent)

	// Draw left eye
	r.fillRoundRect(
		r.eyeLx, r.eyeLy,
		r.eyeLwidthCurrent, r.eyeLheightCurrent,
		borderL, r.eyesColor,
	)

	// Draw right eye unless in cyclops mode
	if !r.cyclops {
		r.fillRoundRect(
			r.eyeRx, r.eyeRy,
			r.eyeRwidthCurrent, r.eyeRheightCurrent,
			borderR, r.eyesColor,
		)
	}
}

// fillCircle draws a filled circle (optimized Bresenham's algorithm)
func (r *RoboEyes) fillCircle(x0, y0, radius int16, corners int, c color.RGBA) {
	if radius <= 0 {
		return
	}

	f := int16(1 - radius)
	ddF_x := int16(1)
	ddF_y := int16(-2 * radius)
	x := int16(0)
	y := radius
	for x <= y {

		if corners&0x1 != 0 { // Top-left quadrant
			r.drawFastHLine(x0-y, y0-x, 2*y+1, c)
			r.drawFastHLine(x0-x, y0-y, 2*x+1, c)
		}
		if corners&0x2 != 0 { // Top-right quadrant
			r.drawFastHLine(x0, y0-y, x, c)
			r.drawFastHLine(x0, y0-x, y, c)
		}

		if corners&0x4 != 0 { // Bottom-left
			r.drawFastHLine(x0-y, y0+x, 2*y+1, c)
			r.drawFastHLine(x0-x, y0+y, 2*x+1, c)
		}
		if corners&0x8 != 0 { // Bottom-right
			r.drawFastHLine(x0, y0+y, x+1, c)
			r.drawFastHLine(x0, y0+x, y+1, c)
		}

		if f >= 0 {
			y--
			ddF_y += 2
			f += ddF_y
		}
		x++
		ddF_x += 2
		f += ddF_x
	}
}

func (r *RoboEyes) fillRoundRect(x, y, width, height, radius int16, c color.RGBA) {
	if height <= 2 || width <= 2 {
		r.fillRect(x, y, width, height, c)
		return
	}

	maxRadius := min(radius, width/2, height/2)
	if maxRadius < 1 {
		maxRadius = 0
	}

	// Draw center rectangle
	r.fillRect(x+maxRadius, y, width-2*maxRadius, height, c)

	// Draw side rectangles
	r.fillRect(x, y+maxRadius, maxRadius, height-2*maxRadius, c)
	r.fillRect(x+width-maxRadius, y+maxRadius, maxRadius, height-2*maxRadius, c)

	// Draw rounded corners only if radius > 0
	if maxRadius > 0 {
		r.fillCircle(x+maxRadius, y+maxRadius, maxRadius, 1, c)                  // Top-left
		r.fillCircle(x+width-maxRadius-1, y+maxRadius, maxRadius, 2, c)          // Top-right
		r.fillCircle(x+maxRadius, y+height-maxRadius-1, maxRadius, 4, c)         // Bottom-left
		r.fillCircle(x+width-maxRadius-1, y+height-maxRadius-1, maxRadius, 8, c) // Bottom-right
	}
}

func min(vals ...int16) int16 {
	minVal := vals[0]
	for _, v := range vals {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func (r *RoboEyes) drawFastHLine(x, y, length int16, c color.RGBA) {
	if length <= 0 {
		return
	}

	// Fallback to manual pixel drawing
	end := x + length
	if end > r.screenWidth {
		end = r.screenWidth
	}
	for ; x < end; x++ {
		r.device.SetPixel(x, y, c)
	}
}

// fillRect fills a rectangle with color
func (r *RoboEyes) fillRect(x, y, width, height int16, c color.RGBA) {
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			if i >= 0 && i < r.screenWidth && j >= 0 && j < r.screenHeight {
				r.device.SetPixel(i, j, c)
			}
		}
	}
}

// drawEyelids renders animated eyelids based on current emotional state
func (r *RoboEyes) drawEyelids() {
	// 1. Determine target eyelid positions based on active emotions
	// Reset all states first to ensure clean transitions
	r.eyelidsTiredHeightNext = 0
	r.eyelidsAngryHeightNext = 0
	r.eyelidsHappyBottomOffsetNext = 0

	// Set next positions based on active emotions (with priority handling)
	if r.tired {
		r.eyelidsTiredHeightNext = r.eyeLheightCurrent / 2
	} else if r.angry {
		r.eyelidsAngryHeightNext = r.eyeLheightCurrent / 2
	}

	if r.happy {
		r.eyelidsHappyBottomOffsetNext = r.eyeLheightCurrent / 2
	}

	// 2. Apply smooth transitions using linear interpolation
	r.eyelidsTiredHeight = (r.eyelidsTiredHeight + r.eyelidsTiredHeightNext) / 2
	r.eyelidsAngryHeight = (r.eyelidsAngryHeight + r.eyelidsAngryHeightNext) / 2
	r.eyelidsHappyBottomOffset = (r.eyelidsHappyBottomOffset + r.eyelidsHappyBottomOffsetNext) / 2

	// Precompute common Y positions for efficiency
	eyeTopY := r.eyeLy - 1 // Top edge of eyes
	eyeBottomY := r.eyeLy + r.eyeLheightCurrent

	// 3. Render tired eyelids - droopy triangles from top corners
	if r.eyelidsTiredHeight > 0 {
		r.drawEyelidTriangles(
			r.eyeLx, eyeTopY, r.eyeLwidthCurrent,
			r.eyeRx, eyeTopY, r.eyeRwidthCurrent,
			r.eyelidsTiredHeight,
			true, // Left eye: triangle points left (from top-left corner)
			true, // Right eye: triangle points right (from top-right corner)
		)
	}

	// 4. Render angry eyelids - inward slanting triangles
	if r.eyelidsAngryHeight > 0 {
		r.drawEyelidTriangles(
			r.eyeLx, eyeTopY, r.eyeLwidthCurrent,
			r.eyeRx, eyeTopY, r.eyeRwidthCurrent,
			r.eyelidsAngryHeight,
			false, // Left eye: triangle points right (from top-right corner)
			false, // Right eye: triangle points left (from top-left corner)
		)
	}

	// 5. Render happy eyelids - bottom curved covers
	if r.eyelidsHappyBottomOffset > 0 {
		// Left eye bottom cover
		r.fillRoundRect(
			r.eyeLx-1, eyeBottomY-r.eyelidsHappyBottomOffset+1,
			r.eyeLwidthCurrent+2, r.eyelidsHappyBottomOffset,
			int16(r.eyeLborderRadiusCurrent), r.bgColor,
		)

		// Right eye bottom cover (only in two-eye mode)
		if !r.cyclops {
			r.fillRoundRect(
				r.eyeRx-1, eyeBottomY-r.eyelidsHappyBottomOffset+1,
				r.eyeRwidthCurrent+2, r.eyelidsHappyBottomOffset,
				int16(r.eyeRborderRadiusCurrent), r.bgColor,
			)
		}
	}
}

// drawEyelidTriangles renders eyelid triangles with consistent parameters
// leftEyeX, leftEyeY: Position of left eye's top-left corner
// leftWidth: Width of left eye
// rightEyeX, rightEyeY: Position of right eye's top-left corner
// rightWidth: Width of right eye
// height: Eyelid height
// leftPointsLeft: Direction for left eye triangle (true = points left, false = points right)
// rightPointsRight: Direction for right eye triangle (true = points right, false = points left)
func (r *RoboEyes) drawEyelidTriangles(
	leftEyeX, leftEyeY int16, leftWidth int16,
	rightEyeX, rightEyeY int16, rightWidth int16,
	height int16,
	leftPointsLeft bool,
	rightPointsRight bool,
) {
	// Cyclops mode uses single-eye rendering
	if r.cyclops {
		midX := leftEyeX + leftWidth/2

		// LEFT PART - Use leftPointsLeft parameter
		if leftPointsLeft {
			// Left-pointing triangle (Tired style)
			r.fillTriangle(
				leftEyeX, leftEyeY,
				midX, leftEyeY,
				leftEyeX, leftEyeY+height,
				r.bgColor,
			)
		} else {
			// Right-pointing triangle (Angry style)
			r.fillTriangle(
				leftEyeX, leftEyeY,
				midX, leftEyeY,
				midX, leftEyeY+height,
				r.bgColor,
			)
		}

		// RIGHT PART - Use rightPointsRight parameter
		if rightPointsRight {
			// Right-pointing triangle (Tired style)
			r.fillTriangle(
				midX, leftEyeY,
				leftEyeX+leftWidth, leftEyeY,
				leftEyeX+leftWidth, leftEyeY+height,
				r.bgColor,
			)
		} else {
			// Left-pointing triangle (Angry style)
			r.fillTriangle(
				midX, leftEyeY,
				leftEyeX+leftWidth, leftEyeY,
				midX, leftEyeY+height,
				r.bgColor,
			)
		}
		return
	}

	// Normal two-eye rendering
	// Left eye triangle
	if leftPointsLeft {
		r.fillTriangle(
			leftEyeX, leftEyeY,
			leftEyeX+leftWidth, leftEyeY,
			leftEyeX, leftEyeY+height,
			r.bgColor,
		)
	} else {
		r.fillTriangle(
			leftEyeX, leftEyeY,
			leftEyeX+leftWidth, leftEyeY,
			leftEyeX+leftWidth, leftEyeY+height,
			r.bgColor,
		)
	}

	// Right eye triangle
	if rightPointsRight {
		r.fillTriangle(
			rightEyeX, rightEyeY,
			rightEyeX+rightWidth, rightEyeY,
			rightEyeX+rightWidth, rightEyeY+height,
			r.bgColor,
		)
	} else {
		r.fillTriangle(
			rightEyeX, rightEyeY,
			rightEyeX+rightWidth, rightEyeY,
			rightEyeX, rightEyeY+height,
			r.bgColor,
		)
	}
}

// fillTriangle fills a triangle with the specified color using scanline rasterization
func (r *RoboEyes) fillTriangle(x0, y0, x1, y1, x2, y2 int16, c color.RGBA) {
	// Sort vertices by ascending y-coordinate (y0 <= y1 <= y2)
	if y0 > y1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if y0 > y1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	totalHeight := y2 - y0
	if totalHeight == 0 {
		return // Degenerate triangle (zero area)
	}

	topHeight := y1 - y0    // Height of top segment
	bottomHeight := y2 - y1 // Height of bottom segment

	// Precompute inverse heights once (avoid division in loops)
	invTotalHeight := 1.0 / float32(totalHeight)
	var invTopHeight, invBottomHeight float32

	// Render top segment (from y0 to y1-1)
	if topHeight > 0 {
		invTopHeight = 1.0 / float32(topHeight)
		for y := y0; y < y1; y++ {
			// Calculate x-coordinates on left and right edges
			t := float32(y-y0) * invTotalHeight
			ax := x0 + int16(float32(x2-x0)*t) // Long edge (x0 to x2)

			t = float32(y-y0) * invTopHeight
			bx := x0 + int16(float32(x1-x0)*t) // Top edge (x0 to x1)

			// Draw horizontal scanline between ax and bx
			r.drawHorizontalLine(ax, bx, y, c)
		}
	}

	// Render bottom segment (from y1 to y2)
	if bottomHeight > 0 {
		invBottomHeight = 1.0 / float32(bottomHeight)
		for y := y1; y <= y2; y++ {
			// Calculate x-coordinates on left and right edges
			t := float32(y-y0) * invTotalHeight
			ax := x0 + int16(float32(x2-x0)*t) // Long edge (x0 to x2)

			t = float32(y-y1) * invBottomHeight
			bx := x1 + int16(float32(x2-x1)*t) // Bottom edge (x1 to x2)

			// Draw horizontal scanline between ax and bx
			r.drawHorizontalLine(ax, bx, y, c)
		}
	}
	// Handle flat-bottom triangles (topHeight=0) where only bottom part renders
}

// drawHorizontalLine draws a clipped horizontal line efficiently
// Arguments are assumed to have: y within [0, screenHeight-1]
func (r *RoboEyes) drawHorizontalLine(xA, xB, y int16, c color.RGBA) {
	// Ensure we draw from left to right
	if xA > xB {
		xA, xB = xB, xA
	}

	// Early exit if scanline is completely off-screen vertically
	if y < 0 || y >= r.screenHeight {
		return
	}

	// Clip x-coordinates to screen bounds
	if xA < 0 {
		xA = 0
	}
	if xB >= r.screenWidth {
		xB = r.screenWidth - 1
	}

	// Draw visible portion of the scanline
	for x := xA; x <= xB; x++ {
		r.device.SetPixel(x, y, c)
	}
}
