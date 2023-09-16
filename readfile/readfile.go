package readfile

import (
	"bufio"
	"fmt"
	"lem-in/models"
	"os"
	"strconv"
	"strings"
)

// ParseData analyse les données d'entrée et renvoie les structures LemInData et Ants.
func ParseData(data []string) (models.LemInData, models.Ants) {
	var lemIn models.LemInData
	var ant models.Ants

	lemIn.Links = make(map[string][]string)

	nbrAnts, err := strconv.Atoi(data[0])
	if err != nil {
		fmt.Printf("Erreur de conversion du nombre de fourmis en entier : %v", err)
		os.Exit(1)
	}
	ant.NbrAnts = nbrAnts
	if ant.NbrAnts < 1 {
		fmt.Printf("ERREUR : format de données invalide, nombre de fourmis invalide\n")
		os.Exit(1)
	}

	for i := 1; i < len(data); i++ {
		if data[i] == "##start" {
			fields := strings.Fields(data[i+1])
			start := fields[0]
			lemIn.Start = start
		} else if data[i] == "##end" {
			fields := strings.Fields(data[i+1])
			end := fields[0]
			lemIn.End = end
		} else if strings.Contains(data[i], "-") {
			link := strings.Split(data[i], "-")
			from := link[0]
			to := link[1]
			if from == "" || to == "" {
				fmt.Printf("Format de lien invalide")
				os.Exit(1)
			}
			lemIn.Links[from] = append(lemIn.Links[from], to)
			lemIn.Links[to] = append(lemIn.Links[to], from)
		} else if strings.Contains(data[i], " ") {
			room := strings.Split(data[i], " ")
			if len(room) == 3 {
				v := room[0]
				if strings.HasPrefix(v, "#") {
					fmt.Printf("Nom de salle invalide : %s", v)
					os.Exit(1)
				}
				x, err := strconv.Atoi(room[1])
				y, err1 := strconv.Atoi(room[2])
				if err1 != nil || err != nil {
					fmt.Println("les coordonnée sont invalid")
					os.Exit(1)
				}
				lemIn.X = x
				lemIn.Y = y
				// Vérifier si le nom de la salle est unique
				if RoomAlready(lemIn.Rooms, v) {
					fmt.Printf("Nom de salle déjà utilisé : %s\n", v)
					os.Exit(1)
				}
				lemIn.Rooms = append(lemIn.Rooms, v)
			}
		}
	}
	return lemIn, ant
}

// ReadFile lit le contenu d'un fichier et le renvoie sous forme de slice de chaînes de caractères.
func ReadFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Impossible d'ouvrir le fichier en raison de l'erreur suivante : %s \n", err)
		os.Exit(1)
	}
	defer file.Close()

	var fileLines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	if len(fileLines) == 0 {
		fmt.Println("Le fichier est vide.")
		os.Exit(1)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Erreur de lecture du fichier : %v", err)
		os.Exit(1)
	}

	return fileLines
}
func RoomAlready(rooms []string, name string) bool {
	for _, room := range rooms {
		if room == name {
			return true
		}
	}
	return false
}
