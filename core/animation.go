// Package core provides animation capabilities for the Spectrex framework.
package core

// AnimationType defines the type of property being animated.
type AnimationType int

const (
	AnimationTypeRotation AnimationType = iota
	AnimationTypePosition
	AnimationTypeColor
	AnimationTypeScale
)

// EaseType defines the easing function type.
type EaseType int

const (
	EaseLinear EaseType = iota
	EaseInOut
	EaseIn
	EaseOut
)

// Animation represents an animation on an object property.
type Animation struct {
	Type         AnimationType
	Target       interface{}
	StartValue   interface{}
	EndValue     interface{}
	CurrentValue interface{}
	Duration     float32
	Timer        float32
	Completed    bool
	EaseType     EaseType
}

// AnimationManager handles all active animations.
type AnimationManager struct {
	Animations []*Animation
}

// NewAnimationManager creates a new animation manager.
func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		Animations: make([]*Animation, 0),
	}
}

// Update updates all animations.
func (am *AnimationManager) Update(deltaTime float32) {
	for i := 0; i < len(am.Animations); i++ {
		anim := am.Animations[i]

		if anim.Completed {
			continue
		}

		anim.Timer += deltaTime

		progress := anim.Timer / anim.Duration
		if progress >= 1.0 {
			progress = 1.0
			anim.Completed = true
		}

		easedProgress := applyEasing(progress, anim.EaseType)

		switch anim.Type {
		case AnimationTypeRotation, AnimationTypePosition, AnimationTypeScale:
			startVal := anim.StartValue.(Vec3)
			endVal := anim.EndValue.(Vec3)
			anim.CurrentValue = Vec3{
				X: startVal.X + (endVal.X-startVal.X)*easedProgress,
				Y: startVal.Y + (endVal.Y-startVal.Y)*easedProgress,
				Z: startVal.Z + (endVal.Z-startVal.Z)*easedProgress,
			}

		case AnimationTypeColor:
			startVal := anim.StartValue.(Color)
			endVal := anim.EndValue.(Color)
			anim.CurrentValue = Color{
				R: uint8(float32(startVal.R) + (float32(endVal.R)-float32(startVal.R))*easedProgress),
				G: uint8(float32(startVal.G) + (float32(endVal.G)-float32(startVal.G))*easedProgress),
				B: uint8(float32(startVal.B) + (float32(endVal.B)-float32(startVal.B))*easedProgress),
				A: uint8(float32(startVal.A) + (float32(endVal.A)-float32(startVal.A))*easedProgress),
			}
		}
	}

	// Remove completed animations
	i := 0
	for i < len(am.Animations) {
		if am.Animations[i].Completed {
			lastIdx := len(am.Animations) - 1
			am.Animations[i] = am.Animations[lastIdx]
			am.Animations = am.Animations[:lastIdx]
		} else {
			i++
		}
	}
}

// AddAnimation adds a new animation.
func (am *AnimationManager) AddAnimation(anim *Animation) {
	am.Animations = append(am.Animations, anim)
}

// SimpleRotation creates a simple rotation animation.
func (am *AnimationManager) SimpleRotation(target interface{}, axis string, startAngle, endAngle, duration float32) *Animation {
	startVal := Vec3{}
	endVal := Vec3{}

	switch axis {
	case "x":
		startVal.X = startAngle
		endVal.X = endAngle
	case "y":
		startVal.Y = startAngle
		endVal.Y = endAngle
	case "z":
		startVal.Z = startAngle
		endVal.Z = endAngle
	}

	anim := &Animation{
		Type:         AnimationTypeRotation,
		Target:       target,
		StartValue:   startVal,
		EndValue:     endVal,
		CurrentValue: startVal,
		Duration:     duration,
		Timer:        0,
		Completed:    false,
		EaseType:     EaseLinear,
	}

	am.AddAnimation(anim)
	return anim
}

// applyEasing applies the easing function to a progress value.
func applyEasing(progress float32, easeType EaseType) float32 {
	switch easeType {
	case EaseInOut:
		if progress < 0.5 {
			return progress * progress * 2
		}
		return 1 - ((1 - progress) * (1 - progress) * 2)
	case EaseIn:
		return progress * progress
	case EaseOut:
		return 1 - ((1 - progress) * (1 - progress))
	default:
		return progress
	}
}
