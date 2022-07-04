package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

var OPTIMAL = "Let's Go!"
var DNA_SIZE = len(OPTIMAL)
var DNA_SIZE_FLOAT = float32(DNA_SIZE)
var POP_SIZE = 500
var GENERATIONS = 5000
var MUTATION_CHANCE float32 = 0.05
var OPTIMAL_ARRAY = get_optimal_array(OPTIMAL)

func get_optimal_array(optimal string) []int {
	opt := make([]int, len(optimal))
	for i := 0; i < len(optimal); i++ {
		opt[i] = int(optimal[i])
	}
	return opt
}

func binarySearchApprox(numbers []float32, leftBound int, rightBound int, numberToFind float32) int {

	// stop when num >= numToFind

	if rightBound >= leftBound {
		midPoint := leftBound + (rightBound-leftBound)/2

		if numbers[midPoint] == numberToFind {
			return midPoint
		}

		if midPoint-1 > 0 && numbers[midPoint] > numberToFind && numbers[midPoint-1] < numberToFind {
			return midPoint
		}

		if numbers[midPoint] > numberToFind {
			return binarySearchApprox(numbers, leftBound, midPoint-1, numberToFind)
		}

		return binarySearchApprox(numbers, midPoint+1, rightBound, numberToFind)
	}

	return -1
}

func weighted_choice(dnas []string, weight_total float32, weights_cumulated []float32) string {

	n := weight_total * rand.Float32()
	i := binarySearchApprox(weights_cumulated, 0, len(weights_cumulated)-1, n)
	if i > -1 {
		return dnas[i]
	}
	return dnas[POP_SIZE-1]
}

func random_char() rune {
	return rune(rand.Intn(126-32) + 32)
}

func random_population() []string {
	pop := make([]string, POP_SIZE)

	for i := 0; i < POP_SIZE; i++ {
		var buffer bytes.Buffer
		for j := 0; j < DNA_SIZE; j++ {
			buffer.WriteRune(random_char())
		}
		pop[i] = buffer.String()
	}
	return pop
}

func fitness(dna string) int {
	var fitness_score int
	for i := 0; i < DNA_SIZE; i++ {
		diff := int(dna[i]) - OPTIMAL_ARRAY[i]
		if diff < 0 {
			diff = diff * -1
		}
		fitness_score += diff
	}
	return fitness_score
}

func mutate(dna string) string {

	if rand.Float32() < MUTATION_CHANCE {
		out := []rune(dna)
		pos := int(rand.Float32() * DNA_SIZE_FLOAT)
		out[pos] = random_char()
		return string(out)
	} else {
		return dna
	}

}

func crossover(dna1 string, dna2 string) (string, string) {
	pos := int(rand.Float32() * DNA_SIZE_FLOAT)
	new_dna1 := dna1[0:pos] + dna2[pos:]
	new_dna2 := dna1[pos:] + dna2[0:pos]
	return new_dna1, new_dna2
}

func main() {
	start_time := time.Now()
	rand.Seed(time.Now().UnixNano())

	population := random_population()

	for i := 0; i < GENERATIONS; i++ {
		//fmt.Printf("Generation %d... Random sample: '%s'... Fitness: %d\n", i, population[0], fitness(population[0]))

		weights := make([]float32, POP_SIZE)
		weighted_population := make([]string, POP_SIZE)

		for j, individual := range population {
			fitness_val := fitness(population[j])

			if fitness_val == 0 {
				weights[j] = 1.0
			} else {
				weights[j] = 1 / float32(fitness_val)
			}
			weighted_population[j] = individual
		}

		population = make([]string, POP_SIZE)

		var weight_total float32
		weights_cumulated := make([]float32, POP_SIZE)

		for i, w := range weights {
			weight_total += w
			if i != 0 {
				weights_cumulated[i] = weights_cumulated[i-1] + w
			} else {
				weights_cumulated[i] = w
			}
		}

		for j := 0; j < POP_SIZE/2; j++ {

			ind1 := weighted_choice(weighted_population, weight_total, weights_cumulated)
			ind2 := weighted_choice(weighted_population, weight_total, weights_cumulated)

			ind1, ind2 = crossover(ind1, ind2)
			population[j] = mutate(ind1)
			population[j+POP_SIZE/2] = mutate(ind2)

		}
	}

	fittest_string := population[0]
	minimum_fitness := fitness(fittest_string)

	for _, individual := range population {
		ind_fitness := fitness(individual)
		if ind_fitness < minimum_fitness {
			fittest_string = individual
			minimum_fitness = ind_fitness
		}
	}

	fmt.Println(fittest_string, minimum_fitness, time.Since(start_time))

}
