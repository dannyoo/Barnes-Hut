package main

import (
	"math"
	"math/rand"
)

const G = 6.67408e-11 // gravitational constant -- don't change this!

const solarMass = 1.989e30 // mass of sun -- don't change this!

const blackHoleMass = 8e36 // mass of black hole -- don't change!

// InitializeUniverse() sets an initial universe given a collection of galaxies and a width.
// It returns a pointer to the resulting universe.
func InitializeUniverse(galaxies []Galaxy, w float64) *Universe {
	var u Universe
	u.width = w
	u.stars = make([]*Star, 0, len(galaxies)*len(galaxies[0]))
	for i := range galaxies {
		for _, b := range galaxies[i] {
			u.stars = append(u.stars, b)
		}
	}
	return &u
}

// InitializeGalaxy takes number of stars in the galaxy, radius of the galaxy to be constructed,
// and center of galaxy to be constructed. Returns a spinning Galaxy object -- which is just a slice of Star pointers
func InitializeGalaxy(numOfStars int, r, x, y float64) Galaxy {
	g := make(Galaxy, numOfStars)

	for i := range g {
		var s Star

		// First choose distance to center of galaxy
		dist := (rand.Float64() + 1.0) / 2.0

		// multiply by factor of r
		dist *= r

		// Next choose the angle in radians to represent the rotation
		angle := rand.Float64() * 2 * math.Pi

		// convert polar coordinates to Cartesian
		s.position.x = x + dist*math.Cos(angle)
		s.position.y = y + dist*math.Sin(angle)

		// set the mass = mass of sun by default
		s.mass = solarMass

		// set the radius equal to radius of sun in m
		s.radius = 696340000

		//set the colors
		s.red = 255
		s.green = 255
		s.blue = 255

		// now spin the galaxy

		// the following is orbital velocity equation
		//dist := Distance(pos, g[i].position)
		speed := 0.5 * math.Sqrt(G*blackHoleMass/dist) // approximation of orbital velocity equation: half of true speed to prevent instability

		s.velocity.x = speed * math.Cos(angle+math.Pi/2.0)
		s.velocity.y = speed * math.Sin(angle+math.Pi/2.0)

		//point g[i] at s
		g[i] = &s

	}

	//add a blackhole to the center of the galaxy

	var blackhole Star
	blackhole.mass = blackHoleMass
	blackhole.position.x = x
	blackhole.position.y = y
	blackhole.blue = 255
	blackhole.radius = 6963400000 // ten times that of a normal star (to make it visible as large)

	g = append(g, &blackhole)

	return g
}

// Returns the pointer to the Jupiter Universe
func Jupiter() *Universe {
	g := make(Galaxy, 5)
	var jupiter, io, europa, ganymede, callisto Star
	jupiter.red, jupiter.green, jupiter.blue = 223, 227, 202
	io.red, io.green, io.blue = 249, 249, 165
	europa.red, europa.green, europa.blue = 132, 83, 52
	ganymede.red, ganymede.green, ganymede.blue = 76, 0, 153
	callisto.red, callisto.green, callisto.blue = 0, 153, 76

	jupiter.mass = 1.898 * math.Pow(10, 27)
	io.mass = 8.9319 * math.Pow(10, 22)
	europa.mass = 4.7998 * math.Pow(10, 22)
	ganymede.mass = 1.4819 * math.Pow(10, 23)
	callisto.mass = 1.0759 * math.Pow(10, 23)

	jupiter.radius = 71000000
	io.radius = 1821000
	europa.radius = 1569000
	ganymede.radius = 2631000
	callisto.radius = 2410000
	
	// the positions were changed so that the "galaxy" would appear in the center
	center := 1.0e21/2
	jupiter.position.x, jupiter.position.y = center, center
	io.position.x, io.position.y = center-center/4.7438330170778, center
	europa.position.x, europa.position.y = center, center+center/2.98107020420331
	ganymede.position.x, ganymede.position.y = center+center/1.86846038863976, center
	callisto.position.x, callisto.position.y = center, center-center/1.06230413767462

	jupiter.velocity.x, jupiter.velocity.y = 0, 0
	io.velocity.x, io.velocity.y = 0, -17320
	europa.velocity.x, europa.velocity.y = -13740, 0
	ganymede.velocity.x, ganymede.velocity.y = 0, 10870
	callisto.velocity.x, callisto.velocity.y = 8200, 0
	
	g[0]= &jupiter
	g[1]= &io
	g[2]= &europa
	g[3]= &ganymede
	g[4]= &callisto

	var u Universe
	u.width = 1.0e21
	u.stars = make([]*Star, 0, len(g))
	u.stars = append(u.stars, g...)
	
	return &u
}
