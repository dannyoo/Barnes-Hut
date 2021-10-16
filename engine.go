package main

import (
	"fmt"
	"math"
)

//BarnesHut is our highest level function.
//Input: initial Universe object, a number of generations, and a time interval.
//Output: collection of Universe objects corresponding to updating the system
//over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens+1)
	timePoints[0] = initialUniverse

	for gen := 0; gen < numGens; gen++ {
		timePoints[gen+1] = UpdateUniverse(timePoints[gen], time, theta)
	}

	return timePoints
}

// func NewNode(star *Star, quad Quadrant) *Node {
// 	return &Node{
// 		children: make([]*Node, 0),
// 		star:   star,
// 		// sector: nil,
// 	}
// }

func CreateQuateTree(univ *Universe) QuadTree {
	tree := QuadTree{
		root: nil,
	}

	tree.root = &Node{
		children: make([]*Node, 4),
		// star: &Star{
		// 	position:     OrderedPair{},
		// 	velocity:     OrderedPair{},
		// 	acceleration: OrderedPair{},
		// 	mass:         0,
		// 	radius:       0,
		// 	red:          0,
		// 	blue:         0,
		// 	green:        0,
		// },
		sector: Quadrant{
			x:     univ.width,
			y:     univ.width,
			width: univ.width,
		},
	}

	for _, s := range univ.stars {
		insert(tree.root, s)
	}
	return tree
}

func insert(node *Node, star *Star) {
	center := OrderedPair{
		x: node.sector.x + (node.sector.width / 2.0),
		y: node.sector.y / 2.0,
	}
	quadIndex := findQuadrantIndex(star, center)
	if quadIndex == -1 {
		panic("problem generating quadIndex")
	}

	// checking if it's there, may be a problem
	// if node.children[quadIndex] == nil || node.children[quadIndex].sector.width == 0
	if node.children[quadIndex] == nil {
		new := Node{
			children: make([]*Node, 0),
			star:     star,
			sector:   getQuadrant(quadIndex, node),
		}
		node.children[quadIndex] = &new

	} else if len(node.children[quadIndex].children) == 4 {
		// node x is an internal node
		// update the center-of-mass and total mass of x
		ratio := (node.star.mass / star.mass)
		node.star.position = OrderedPair{
			x: (ratio*node.star.position.x + star.position.x) / (1 + ratio),
			y: (ratio*node.star.position.y + star.position.y) / (1 + ratio),
		}
		node.star.mass += star.mass
		insert(node.children[quadIndex], star)
	} else if len(node.children[quadIndex].children) == 0 {
		//external node
		// copy  := make([]*Node, 1)
		// copy[0] = node.children[quadIndex]
		// old := copy[0].star
		// way2
		// old:= node.children[quadIndex].star // may not work
		// psuedostar := Node{
		// 	children: make([]*Node, 4),
		// 	sector: node.children[quadIndex].sector,
		// }
		// node.children[quadIndex] = &psuedostar
		// insert(&psuedostar, old)
		// insert(&psuedostar, star)

		// way 3
		// var test *Star
		node.children[quadIndex].children = make([]*Node, 4) //create 4 children
		// test = node.children[quadIndex].star
		// node.children[quadIndex].star = nil
		// fmt.Println(test, "<- if this is nil value, the algo will be broken, but if not... we gucci")
		insert(node.children[quadIndex], node.children[quadIndex].star)
		insert(node.children[quadIndex], star)
		node.children[quadIndex].star.mass += star.mass
		ratio := (node.children[quadIndex].star.mass / star.mass)
		node.children[quadIndex].star.position = OrderedPair{
			x: (ratio*node.children[quadIndex].star.position.x + star.position.x) / (1 + ratio),
			y: (ratio*node.children[quadIndex].star.position.y + star.position.y) / (1 + ratio),
		}
		// update the center-of-mass and total mass of x

	} else {
		panic("I didn't plan for this....")
	}
}

func getQuadrant(quadIndex int, node *Node) Quadrant {
	var quad Quadrant
	switch quadIndex {
	case 0:
		// NW
		quad = Quadrant{
			x:     node.sector.x,
			y:     node.sector.y / 2.0,
			width: node.sector.width / 2.0,
		}
	case 1:
		// NE
		quad = Quadrant{
			x:     node.sector.x + (node.sector.width / 2.0),
			y:     node.sector.y / 2.0,
			width: node.sector.width / 2.0,
		}
	case 2:
		// SW
		quad = Quadrant{
			x:     node.sector.x,
			y:     node.sector.y,
			width: node.sector.width / 2.0,
		}
	case 3:
		//SE
		quad = Quadrant{
			x:     node.sector.x + (node.sector.width / 2.0),
			y:     node.sector.y,
			width: node.sector.width / 2.0,
		}
	default:
		fmt.Println("go, sucks")
		panic("not good generating quadrant from quadIndex")
	}

	return quad
}

func findQuadrantIndex(star *Star, center OrderedPair) int {
	// var index int
	// find center, make center 0,0 by dividing
	if star.position.x <= center.x && star.position.y <= center.y {
		// bottom left SW
		// index = 2
		return 2
	} else if star.position.x <= center.x && star.position.y >= center.y {
		// top left NW
		// index = 0
		return 0
	} else if star.position.x >= center.x && star.position.y <= center.y {
		// bottom right SE
		// index = 3
		return 3
	} else if star.position.x >= center.x && star.position.y >= center.y {
		// top right NE
		// index = 1
		return 1
	}
	return -1

}

func CreateComparableUniverse(tree *QuadTree, theta float64, X *Star) (uni Universe) {

	uni.stars = thetaStars(tree.root, theta, X)

	return
}

func thetaStars(node *Node, theta float64, X *Star) []*Star {
	var starrys = make([]*Star, 0)
	for _, single := range node.children{
		if single.star == nil{
			continue
		} else if len(single.children) == 4{ // its a internal node
			s := single.sector.width
			d:= Dist(*single.star, *X)
			heuristic := s/d
			if heuristic > theta {
				starrys = append(starrys, thetaStars(single, theta, X)...)
			} else if heuristic <= theta{
				starrys = append(starrys, single.star)
			}
		} else if len(single.children) == 0 { // external node
			starrys = append(starrys, single.star)
		}
	}
	return starrys
}

// UpdateUniverse returns a new universe after time t.
func UpdateUniverse(univ *Universe, t, theta float64) *Universe {

	newUniverse := CopyUniverse(univ)
	tree := CreateQuateTree(univ)
	for b := range univ.stars {
		comparableUniverse := CreateComparableUniverse(&tree, theta, newUniverse.stars[b])
		// update pos, vel and accel
		newUniverse.stars[b].Update(&comparableUniverse, t)
	}

	return newUniverse
}

func (b *Star) Update(univ *Universe, t float64) {
	acc := b.NewAccel(univ)
	vel := b.NewVelocity(t)
	pos := b.NewPosition(t)
	b.acceleration, b.velocity, b.position = acc, vel, pos
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
