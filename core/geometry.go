// Package core provides geometry utilities for the Spectrex framework.
package core

import "math"

// MakePoly generates a polygon coordinate list given number of sides, radius & start angle.
func MakePoly(sides int, radius float32, start float32) []Vec3 {
	theta := float64(start)
	delta := (2.0 * math.Pi) / float64(sides)

	poly := make([]Vec3, sides)
	for vertex := 0; vertex < sides; vertex++ {
		poly[vertex] = Vec3{
			X: radius * float32(math.Cos(theta)),
			Y: radius * float32(math.Sin(theta)),
			Z: 0,
		}
		theta += delta
	}

	return poly
}

// TransformPoly applies rotation and translation to a polygon.
func TransformPoly(poly []Vec3, position Vec3, rotation Vec3) []Vec3 {
	transformed := make([]Vec3, len(poly))

	rotX := DegToRad(rotation.X)
	rotY := DegToRad(rotation.Y)
	rotZ := DegToRad(rotation.Z)

	for i, v := range poly {
		p := v

		// Apply X-axis rotation
		if rotX != 0 {
			cosX := float32(math.Cos(float64(rotX)))
			sinX := float32(math.Sin(float64(rotX)))
			y := p.Y*cosX - p.Z*sinX
			z := p.Y*sinX + p.Z*cosX
			p.Y = y
			p.Z = z
		}

		// Apply Y-axis rotation
		if rotY != 0 {
			cosY := float32(math.Cos(float64(rotY)))
			sinY := float32(math.Sin(float64(rotY)))
			x := p.X*cosY + p.Z*sinY
			z := -p.X*sinY + p.Z*cosY
			p.X = x
			p.Z = z
		}

		// Apply Z-axis rotation
		if rotZ != 0 {
			cosZ := float32(math.Cos(float64(rotZ)))
			sinZ := float32(math.Sin(float64(rotZ)))
			x := p.X*cosZ - p.Y*sinZ
			y := p.X*sinZ + p.Y*cosZ
			p.X = x
			p.Y = y
		}

		// Apply translation
		transformed[i] = Vec3{
			X: p.X + position.X,
			Y: p.Y + position.Y,
			Z: p.Z + position.Z,
		}
	}

	return transformed
}
