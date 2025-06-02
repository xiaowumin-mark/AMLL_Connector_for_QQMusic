package main

import (
	"time"

	"github.com/creasty/go-easing"
)

// Animation represents a single animation instance
type Animation struct {
	from       []float64
	to         []float64
	duration   time.Duration
	easeFunc   func(float64) float64
	onUpdate   func(values []float64)
	onComplete func()
}

// Animate creates a new animation instance
func Animate(from, to []float64, duration time.Duration, easeFunc func(float64) float64) *Animation {
	if len(from) != len(to) {
		panic("from and to arrays must have the same length")
	}
	return &Animation{
		from:     from,
		to:       to,
		duration: duration,
		easeFunc: easeFunc,
	}
}

// OnUpdate sets the update callback
func (a *Animation) OnUpdate(callback func(values []float64)) *Animation {
	a.onUpdate = callback
	return a
}

// OnComplete sets the complete callback
func (a *Animation) OnComplete(callback func()) *Animation {
	a.onComplete = callback
	return a
}

// Start begins the animation
func (a *Animation) Start() {
	startTime := time.Now()
	go func() {
		ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
		defer ticker.Stop()

		for {
			<-ticker.C
			elapsed := time.Since(startTime)
			t := float64(elapsed) / float64(a.duration)
			if t > 1 {
				t = 1
			}

			// Calculate eased values for all parameters
			values := make([]float64, len(a.from))
			for i, fromValue := range a.from {
				values[i] = easing.Transition(fromValue, a.to[i], a.easeFunc(t))
			}

			// Call the update callback
			if a.onUpdate != nil {
				a.onUpdate(values)
			}

			// End the animation if time is up
			if t == 1 {
				if a.onComplete != nil {
					a.onComplete()
				}
				break
			}
		}
	}()
}
