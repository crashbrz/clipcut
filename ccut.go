package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	delimiter := flag.String("d", "\t", "Delimitador (padrão: TAB)")
	fieldsStr := flag.String("f", "", "Campos a selecionar (ex: 1,3-5,7-)")
	charsStr := flag.String("c", "", "Selecionar por caracteres (ex: 1-10,15)")
	suppress := flag.Bool("s", false, "Suprimir linhas sem o delimitador (-s)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: clipcut [opções]\n")
		fmt.Fprintf(os.Stderr, "   Faz 'cut' direto no conteúdo da área de transferência.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExemplos:\n")
		fmt.Fprintf(os.Stderr, "  clipcut -f 1,3\n")
		fmt.Fprintf(os.Stderr, "  clipcut -d, -f 2-4\n")
		fmt.Fprintf(os.Stderr, "  clipcut -d' ' -f 1-\n")
		fmt.Fprintf(os.Stderr, "  clipcut -c 1-20\n")
	}
	flag.Parse()

	// Lê o clipboard
	text, err := clipboard.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler área de transferência: %v\n", err)
		os.Exit(1)
	}

	if text == "" {
		fmt.Println("Área de transferência vazia.")
		return
	}

	// Normaliza quebras de linha
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")

	for _, line := range lines {
		if line == "" {
			fmt.Println()
			continue
		}

		var result string

		if *charsStr != "" {
			// Modo caracteres (-c)
			result = cutByChars(line, *charsStr)
		} else if *fieldsStr != "" {
			// Modo campos (-f)
			result = cutByFields(line, *delimiter, *fieldsStr, *suppress)
		} else {
			// Sem -f nem -c → comportamento padrão do cut (mostra linha inteira)
			result = line
		}

		if result != "" || !*suppress {
			fmt.Println(result)
		}
	}
}

// cutByFields implementa a lógica de -f (campos)
func cutByFields(line, delim, fieldsStr string, suppress bool) string {
	if !strings.Contains(line, delim) {
		if suppress {
			return ""
		}
		return line // cut original repete a linha quando não tem delimitador (sem -s)
	}

	fields := strings.Split(line, delim)
	selected := parseFieldList(fieldsStr, len(fields))

	var parts []string
	for _, idx := range selected {
		if idx > 0 && idx <= len(fields) {
			parts = append(parts, fields[idx-1])
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, delim)
}

// cutByChars implementa a lógica de -c (caracteres)
func cutByChars(line, charsStr string) string {
	selected := parseRangeList(charsStr, len(line))
	if len(selected) == 0 {
		return ""
	}

	var result []rune
	for _, pos := range selected {
		if pos > 0 && pos <= len(line) {
			result = append(result, rune(line[pos-1]))
		}
	}
	return string(result)
}

// parseFieldList converte "1,3-5,7-" em slice de índices
func parseFieldList(s string, maxFields int) []int {
	var result []int
	parts := strings.Split(s, ",")

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if strings.Contains(p, "-") {
			rng := strings.Split(p, "-")
			start := 1
			end := maxFields

			if rng[0] != "" {
				start, _ = strconv.Atoi(rng[0])
			}
			if len(rng) > 1 && rng[1] != "" {
				end, _ = strconv.Atoi(rng[1])
			}

			for i := start; i <= end && i <= maxFields; i++ {
				result = append(result, i)
			}
		} else {
			num, _ := strconv.Atoi(p)
			if num > 0 {
				result = append(result, num)
			}
		}
	}
	return result
}

// parseRangeList para -c (suporte simples a ranges)
func parseRangeList(s string, maxLen int) []int {
	var result []int
	parts := strings.Split(s, ",")

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.Contains(p, "-") {
			rng := strings.Split(p, "-")
			start := 1
			end := maxLen

			if rng[0] != "" {
				start, _ = strconv.Atoi(rng[0])
			}
			if len(rng) > 1 && rng[1] != "" {
				end, _ = strconv.Atoi(rng[1])
			}

			for i := start; i <= end && i <= maxLen; i++ {
				result = append(result, i)
			}
		} else {
			num, _ := strconv.Atoi(p)
			if num > 0 {
				result = append(result, num)
			}
		}
	}
	return result
}
