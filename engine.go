package main

import "math"

//BarnesHut is our highest level function.
//Input: initial Universe object, a number of generations, and a time interval.
//Output: collection of Universe objects corresponding to updating the system
//over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens+1)
	timePoints[0] = initialUniverse
	
	for gen := 0; gen < numGens; gen++ {
		timePoints[gen+1] = UpdateUniverse(timePoints[gen], time)
	}

	return timePoints
}

// UpdateUniverse returns a new universe after time t.
func UpdateUniverse(univ *Universe, t float64) *Universe {

	newUniverse := CopyUniverse(univ)
	for b := range univ.stars {
		// update pos, vel and accel
		newUniverse.stars[b].Update(univ, t)
	}

	return newUniverse
}

func (b *Star) Update(univ *Universe, t float64) {
	acc := b.NewAccel(univ)
	vel := b.NewVelocity(t)
	pos := b.NewPosition(t)
	b.acceleration, b.velocity, b.position = acc, vel, pos
}


// NewVelocity makes the velocity of this object consistent with the acceleration.
func (b *Star) NewVelocity(t float64) OrderedPair {
	return OrderedPair{
		x: b.velocity.x + b.acceleration.x*t,
		y: b.velocity.y + b.acceleration.y*t,
	}
}

// NewPosition computes the new poosition given the updated acc and velocity.
//
// Assumputions: constant acceleration over a time step.
// => DeltaX = v_avg * t
//    DeltaX = (v_start + v_final)*t/ 2
// because v_final = v_start + acc*t:
//	  DeltaX = (v_start + v_start + acc*t)t/2
// Simplify:
//	DeltaX = v_start*t + 0.5acc*t*t
// =>
//  NewX = v_start*t + 0.5acc*t*t + OldX
//
func (b *Star) NewPosition(t float64) OrderedPair {
	return OrderedPair{
		x: b.position.x + b.velocity.x*t + 0.5*b.acceleration.x*t*t,
		y: b.position.y + b.velocity.y*t + 0.5*b.acceleration.y*t*t,
	}
}

// UpdateAccel computes the new accerlation vector for b
func (b *Star) NewAccel(univ *Universe) OrderedPair {
	F := ComputeNetForce(*univ, *b)
	return OrderedPair{
		x: F.x / b.mass,
		y: F.y / b.mass,
	}
}

// ComputeNetForce sums the forces of all bodies in the universe
// acting on b.
func ComputeNetForce(univ Universe, b Star) OrderedPair {
	var netForce OrderedPair
	for _, body := range univ.stars {
		if *body != b {
			f := ComputeGravityForce(b, *body)
			netForce.Add(f)
		}
	}
	return netForce
}

// ComputeGravityForce computes the gravity force between body 1 and body 2.
func ComputeGravityForce(b1, b2 Star) OrderedPair {
	d := Dist(b1, b2)
	deltaX := b2.position.x - b1.position.x
	deltaY := b2.position.y - b1.position.y
	F := G * b1.mass * b2.mass / (d * d)

	return OrderedPair{
		x: F * deltaX / d,
		y: F * deltaY / d,
	}
}
// Compute the Euclidian Distance between two stars
func Dist(b1, b2 Star) float64 {
	dx := b1.position.x - b2.position.x
	dy := b1.position.y - b2.position.y
	return math.Sqrt(dx*dx + dy*dy)
}