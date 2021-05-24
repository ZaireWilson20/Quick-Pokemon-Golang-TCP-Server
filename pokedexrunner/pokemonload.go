package pokeload

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//Pokemon Struct
type Pokemon struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func main() {

	fmt.Println("Wassup")
	Load(os.Args[1])
	//fmt.Println("HI")
}

//Load CSV files
func Load(file string) []Pokemon {
	fmt.Println(file)
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(csvFile)

	var pokemon []Pokemon
	lineNum := 0
	//fmt.Println("hello")

	for {
		if lineNum != 0 {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				fmt.Println("come on...")
				log.Fatal(error)
			}
			pokemon = append(pokemon, Pokemon{
				Name: line[2],
				Type: line[10],
			})
		} else {
			lineNum++
		}
	}
	//pokemonJSON, _ := json.Marshal(pokemon)
	//fmt.Println(string(pokemonJSON))
	return pokemon
}

//Find pokemon
func Find(table []Pokemon, filter string) []Pokemon {
	if filter == "" || filter == "*" {
		return table
	}
	result := make([]Pokemon, 0)
	filter = strings.ToUpper(filter)
	for _, pok := range table {
		if strings.ToUpper(pok.Type) == filter ||
			strings.ToUpper(pok.Name) == filter {
			result = append(result, pok)
		}
	}
	return result
}
