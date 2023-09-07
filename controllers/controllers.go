package controllers

import (
	"fmt"
	"lem-in/models"
	"math"
	"os"
	"strconv"
)

// BestCombPaths trouve la meilleure combinaison de chemins pour les fourmis.
func BestCombPaths(ants models.Ants, pathGroups models.Paths) models.Paths {
	minLevel := math.MaxInt32
	var bestCombination [][]string

	for _, combination := range pathGroups.SortComb {
		spaceInPath := 0
		totalPathLength := 0
		levelOfAnts := 0
		longestPath := len(combination[len(combination)-1])
		for i := 1; i < len(combination); i++ {
			spaceInPath = spaceInPath + longestPath - len(combination[i])
			totalPathLength = totalPathLength + len(combination[i])
		}
		levelOfAnts = (ants.NbrAnts-spaceInPath)/len(combination) + longestPath

		if levelOfAnts < minLevel {
			minLevel = levelOfAnts
			bestCombination = combination
		}
	}

	return models.Paths{BestComb: bestCombination}
}

// SortCombPaths trie les combinaisons de chemins.
func SortCombPaths(combinations models.Paths) models.Paths {
	pathGroups := make(map[int][][]string)
	for _, combination := range combinations.ValidPaths {
		category := len(combination)
		currentCombLength := combLength(combination)
		if _, ok := pathGroups[category]; ok {
			valueInMap := pathGroups[category]
			if currentCombLength < combLength(valueInMap) {
				pathGroups[category] = combination
			}
		} else {
			pathGroups[category] = combination
		}
	}
	return models.Paths{SortComb: pathGroups}
}

// combLength calcule la longueur d'une combinaison de chemins.
func combLength(combination [][]string) int {
	length := 0
	for _, path := range combination {
		length = length + len(path)
	}
	return length
}

// SendAnts attribue les fourmis aux chemins et calcule leurs mouvements.
func SendAnts(ants models.Ants, data []string, bestCombination models.Paths) []string {
	inputData(data)
	queue := assignAnts(ants, bestCombination)
	determineOrder(queue)
	return calculateSteps(queue, bestCombination)
}

// determineOrder détermine l'ordre dans lequel les fourmis se déplaceront.
func determineOrder(queue [][]string) []int {
	order := []int{}
	longest := len(queue[0])
	for i := 0; i < len(queue); i++ {
		if len(queue[i]) > longest {
			longest = len(queue[i])
		}
	}
	for j := 0; j < longest; j++ {
		for i := 0; i < len(queue); i++ {
			if j < len(queue[i]) {
				x, _ := strconv.Atoi(queue[i][j])
				order = append(order, x)
			}
		}
	}
	return order
}

// assignAnts attribue les fourmis aux chemins.
func assignAnts(ants models.Ants, bestCombination models.Paths) [][]string {
	queue := make([][]string, len(bestCombination.BestComb))

	for i := 1; i <= ants.NbrAnts; i++ {
		ant := strconv.Itoa(i)
		minSteps := len(bestCombination.BestComb[0]) + len(queue[0])
		minIndex := 0

		for j, path := range bestCombination.BestComb {
			steps := len(path) + len(queue[j])
			if steps < minSteps {
				minSteps = steps
				minIndex = j
			}
		}
		queue[minIndex] = append(queue[minIndex], ant)
	}
	return queue
}

// inputData affiche les données d'entrée (pour le débogage).
func inputData(data []string) {
	for _, v := range data {
		fmt.Println(v)
	}
	fmt.Println()
}

// calculateSteps calcule les mouvements des fourmis.
func calculateSteps(queue [][]string, bestCombination models.Paths) []string {
	container := make([][][]string, len(queue))
	for i, path := range queue {
		for _, ant := range path {
			adder := []string{}
			for _, vertex := range bestCombination.BestComb[i] {
				str := "L" + ant + "-" + vertex
				adder = append(adder, str)
			}
			container[i] = append(container[i], adder)
		}
	}

	finalMoves := []string{}
	for _, paths := range container {
		for j, moves := range paths {
			for k, vertex := range moves {
				if j+k > len(finalMoves)-1 {
					finalMoves = append(finalMoves, vertex+" ")
				} else {
					finalMoves[j+k] = finalMoves[j+k] + vertex + " "
				}
			}
		}
	}
	return finalMoves
}

// FindPaths trouve les chemins possibles dans le graphe.
func FindPaths(lemin models.LemInData) models.Paths {
	_, ok := lemin.Links[lemin.Start]
	if !ok || lemin.Start == "" {
		fmt.Println("ERROR: invalid data format, no start room found")
		os.Exit(1)
	}
	_, ok = lemin.Links[lemin.End]
	if !ok || lemin.End == "" {
		fmt.Println("ERROR: invalid data format, no end room found")
		os.Exit(1)
	}

	visited := make(map[string]bool)
	var path []string
	var paths [][]string

	find(lemin.Start, lemin.End, lemin.Links, visited, path, &paths)
	return models.Paths{AllPaths: paths}
}

// find trouve les chemins possibles entre le début et la fin.
func find(start, end string, edges map[string][]string, visited map[string]bool, path []string, paths *[][]string) {
	visited[start] = true
	path = append(path, start)

	if start == end {
		temp := make([]string, len(path[1:]))
		copy(temp, path[1:])
		*paths = append(*paths, temp)
	} else {
		for _, vertex := range edges[start] {
			if !visited[vertex] {
				find(vertex, end, edges, visited, path, paths)
			}
		}
	}
	visited[start] = false
}

// ValidPaths renvoie les chemins valides après vérification des interceptions.
func ValidPaths(paths models.Paths) models.Paths {
	sort(paths)
	var result [][][]string
	for i, path := range paths.AllPaths {
		var ValidPaths [][]string
		ValidPaths = append(ValidPaths, path)
		result = append(result, ValidPaths)
		for j := i + 1; j < len(paths.AllPaths); j++ {
			if !checkIntercept(ValidPaths, paths.AllPaths[j]) {
				ValidPaths = append(ValidPaths, paths.AllPaths[j])
				result = append(result, ValidPaths)
			}
		}
	}
	return models.Paths{ValidPaths: result}
}

// sort trie les chemins par longueur.
func sort(paths models.Paths) {
	for i := 0; i < len(paths.AllPaths); i++ {
		for j := i + 1; j < len(paths.AllPaths)-1; j++ {
			if len(paths.AllPaths[i]) > len(paths.AllPaths[j+1]) {
				paths.AllPaths[j+1], paths.AllPaths[i] = paths.AllPaths[i], paths.AllPaths[j+1]
			}
		}
	}
}

// checkIntercept vérifie s'il y a une interception entre deux chemins.
func checkIntercept(ValidPaths [][]string, path []string) bool {
	for _, paths := range ValidPaths {
		for _, room1 := range paths[:len(paths)-1] {
			for _, room2 := range path[:len(path)-1] {
				if room1 == room2 {
					return true // Il y a une interception
				}
			}
		}
	}
	return false // Pas d'interception
}
