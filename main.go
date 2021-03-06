package main

import (
	"fmt"
	"gifhelper"
	"os"
)

func main() {
	command := os.Args[1]
	if command != "jupiter" && command != "galaxy" && command != "collision" {
		panic("Usage './barnes-hut command' where command is 'jupiter' or 'galaxy' or 'collision'")
	}

	var initialUniverse *Universe

	if command == "collision" {
		// the following sample parameters may be helpful for the "collide" command
		// all units are in SI (meters, kg, etc.)
		// but feel free to change the positions of the galaxies.

		g0 := InitializeGalaxy(500, 4e21, 3.5e22, 3.8e22)
		g1 := InitializeGalaxy(500, 4e21, 5e22, 3.1e22)

		// This code is optional if pushing every star is desired
		// for _, x := range g0{
		// 	x.velocity.x = 1e3
		// }
		// for _, x := range g1{
		// 	x.velocity.x = -1e3
		// }
		
		//just pushing the blackholes looks more interesting
		g0[len(g0)-1].velocity.x = 1e3
		g1[len(g0)-1].velocity.x = -1e3

		// you probably want to apply a "push" function at this point to these galaxies to move
		// them toward each other to collide.
		// be careful: if you push them too fast, they'll just fly through each other.
		// too slow and the black holes at the center collide and hilarity ensues.

		width := 1.0e23
		galaxies := []Galaxy{g0, g1}
		initialUniverse = InitializeUniverse(galaxies, width)

	} else if command == "galaxy" {
		g0 := InitializeGalaxy(500, 10e21, 1.0e23/2, 1.0e23/2)
		width := 1.0e23
		galaxies := []Galaxy{g0}
		initialUniverse = InitializeUniverse(galaxies, width)
	} else if command == "jupiter" {
		initialUniverse = Jupiter()
	}

	// now evolve the universe: feel free to adjust the following parameters.
	numGens := 10000
	time := 2e14
	if command == "jupiter" {
		time = 1.5e11 // this makes it easier to see
	}
	theta := 0.5

	timePoints := BarnesHut(initialUniverse, numGens, time, theta)

	fmt.Println("Simulation run. Now drawing images.")
	canvasWidth := 1000
	frequency := 1000
	scalingFactor := 1e11 
	if command == "jupiter" {
		scalingFactor = 10e11 // this makes the moons easier to see
	}
	// a scaling factor is needed to inflate size of stars when drawn because galaxies are very sparse
	imageList := AnimateSystem(timePoints, canvasWidth, frequency, scalingFactor)

	fmt.Println("Images drawn. Now generating GIF.")
	gifhelper.ImagesToGIF(imageList, "galaxy")
	fmt.Println("GIF drawn.")
}
