// main.go

package main

import (
	"fmt"
	"lem-in/controllers"
	"lem-in/readfile"
	"os"
)

func main() {
	// Vérifie si un argument (nom de fichier) est fourni
	if len(os.Args) != 2 {
		fmt.Println("Error: missing filename")
		fmt.Println("USAGE: go run main.go example00.txt")
		return
	}

	// Récupère le chemin du fichier à partir des arguments de la ligne de commande
	filePath := os.Args[1]

	// Lit les données à partir du fichier
	data := readfile.ReadFile(filePath)

	// Analyse les données pour obtenir le graphe et les informations sur les fourmis
	graph, ants := readfile.ParseData(data)

	// Trouve les chemins possibles dans le graphe
	paths := controllers.FindPaths(graph)

	// Filtre les chemins valides et détermine la meilleure combinaison de chemins
	validPaths := controllers.ValidPaths(paths)
	sortedCombPaths := controllers.SortCombPaths(validPaths)
	bestCombPaths := controllers.BestCombPaths(ants, sortedCombPaths)

	// Envoie les fourmis sur les chemins et calcule les mouvements
	sendAnts := controllers.SendAnts(ants, data, bestCombPaths)

	// Affiche les mouvements des fourmis
	turns := 0
	for _, v := range sendAnts {
		fmt.Println(v)
		turns++
	}
	fmt.Println()
	fmt.Printf("Turns number: %v\n", turns)
}
