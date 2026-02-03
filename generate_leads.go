package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const totalLeads = 10000

var firstNames = []string{
	"João", "Maria", "Pedro", "Ana", "Carlos", "Lucas",
	"Fernanda", "Rafael", "Juliana", "Bruno", "Camila",
}

var lastNames = []string{
	"Silva", "Santos", "Oliveira", "Pereira", "Costa",
	"Rodrigues", "Alves", "Lima", "Gomes", "Ribeiro",
}

var tagsPool = []string{
	"novo", "vip", "retorno", "interessado",
	"premium", "lead-frio", "lead-quente",
}

func randomPhone() string {
	ddd := rand.Intn(90) + 10
	number := rand.Intn(90000000) + 10000000
	return fmt.Sprintf("55%d9%d", ddd, number)
}

func randomCPF() string {
	return fmt.Sprintf("%011d", rand.Int63n(99999999999))
}

func randomName() string {
	return fmt.Sprintf(
		"%s %s",
		firstNames[rand.Intn(len(firstNames))],
		lastNames[rand.Intn(len(lastNames))],
	)
}

func randomEmail(name string) string {
	base := strings.ToLower(strings.ReplaceAll(name, " ", "."))
	return fmt.Sprintf("%s%d@email.com", base, rand.Intn(1000))
}

func randomTags() string {
	count := rand.Intn(6) // 0 a 5
	if count == 0 {
		return ""
	}

	selected := map[string]bool{}
	for len(selected) < count {
		tag := tagsPool[rand.Intn(len(tagsPool))]
		selected[tag] = true
	}

	var tags []string
	for tag := range selected {
		tags = append(tags, tag)
	}

	return strings.Join(tags, ", ")
}

func main() {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())

	file, err := os.Create("leads_10000.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	writer.Write([]string{"name", "phone", "cpf", "email", "tags"})

	for i := 0; i < totalLeads; i++ {
		name := randomName()

		record := []string{
			name,
			randomPhone(),
			randomCPF(),
			randomEmail(name),
			randomTags(),
		}

		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("✅ Arquivo leads_10000.csv gerado com sucesso! Tempo de execução: %.2f segundos\n", elapsed.Seconds())
}
