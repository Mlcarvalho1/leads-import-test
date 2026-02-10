//go:build ignore

package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

var (
	firstNames = []string{"Ana", "Bruno", "Carlos", "Diana", "Eduardo", "Fernanda", "Gabriel", "Helena", "Igor", "Julia", "Lucas", "Maria", "Nicolas", "Olivia", "Pedro", "Rafaela", "Samuel", "Tatiana", "Vitor", "Yasmin"}
	lastNames  = []string{"Silva", "Santos", "Oliveira", "Souza", "Pereira", "Costa", "Ferreira", "Almeida", "Nascimento", "Lima", "Araújo", "Ribeiro", "Carvalho", "Gomes", "Martins", "Rocha", "Rodrigues", "Moreira", "Barbosa", "Lopes"}
	tagPool    = []string{"ortho", "implant", "cleaning", "whitening", "braces", "checkup", "urgent", "vip", "returning", "new-patient"}
)

func randomName() string {
	return firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
}

// Valid Brazilian DDDs (area codes)
var validDDDs = []int{
	11, 12, 13, 14, 15, 16, 17, 18, 19, // SP
	21, 22, 24, // RJ
	27, 28, // ES
	31, 32, 33, 34, 35, 37, 38, // MG
	41, 42, 43, 44, 45, 46, // PR
	47, 48, 49, // SC
	51, 53, 54, 55, // RS
	61, // DF
	62, 64, // GO
	63, // TO
	65, 66, // MT
	67, // MS
	68, // AC
	69, // RO
	71, 73, 74, 75, 77, // BA
	79, // SE
	81, 82, // PE/AL
	83, // PB
	84, // RN
	85, 88, // CE
	86, 89, // PI
	87, // PE
	91, 93, 94, // PA
	92, 97, // AM
	95, // RR
	96, // AP
	98, 99, // MA
}

func randomPhone(i int) string {
	ddd := validDDDs[rand.Intn(len(validDDDs))]
	// Generate 8-digit subscriber number: first digit 1-9, rest 0-9
	subscriber := 10000000 + rand.Intn(90000000)
	return fmt.Sprintf("+55%d9%d", ddd, subscriber)
}

func randomCPF() string {
	digits := make([]int, 9)
	for i := range digits {
		digits[i] = rand.Intn(10)
	}

	// First check digit
	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (10 - i)
	}
	r := (sum * 10) % 11
	if r == 10 {
		r = 0
	}
	digits = append(digits, r)

	// Second check digit
	sum = 0
	for i := 0; i < 10; i++ {
		sum += digits[i] * (11 - i)
	}
	r = (sum * 10) % 11
	if r == 10 {
		r = 0
	}
	digits = append(digits, r)

	s := ""
	for _, d := range digits {
		s += fmt.Sprintf("%d", d)
	}
	return s
}

func randomTags() string {
	n := rand.Intn(4) // 0-3 tags
	if n == 0 {
		return ""
	}
	picked := make([]string, 0, n)
	used := map[int]bool{}
	for len(picked) < n {
		idx := rand.Intn(len(tagPool))
		if !used[idx] {
			used[idx] = true
			picked = append(picked, tagPool[idx])
		}
	}
	return strings.Join(picked, ", ")
}

func main() {
	f, err := os.Create("leads_5000.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"name", "phone", "cpf", "email", "tags"})

	for i := 1; i <= 5000; i++ {
		name := randomName()
		phone := randomPhone(i)

		cpf := ""
		if rand.Float64() < 0.7 {
			cpf = randomCPF()
		}

		email := ""
		if rand.Float64() < 0.8 {
			email = fmt.Sprintf("lead%d@example.com", i)
		}

		tags := randomTags()

		w.Write([]string{name, phone, cpf, email, tags})
	}

	fmt.Println("Generated leads_5000.csv with 5000 rows")
}
