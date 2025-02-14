package utils

import (
	"strings"
	"unicode"
)

var streetTypeMap = map[string]string{
	"R":     "RUA",
	"R.":    "RUA",
	"RUA":   "RUA",
	"AV":    "AVENIDA",
	"AV.":   "AVENIDA",
	"AVEN":  "AVENIDA",
	"AVE":   "AVENIDA",
	"AL":    "ALAMEDA",
	"AL.":   "ALAMEDA",
	"EST":   "ESTRADA",
	"EST.":  "ESTRADA",
	"PÇ":    "PRACA",
	"PC":    "PRACA",
	"PCA":   "PRACA",
	"PÇA":   "PRACA",
	"PRAÇA": "PRACA",
}

func NormalizeStreetType(streetType string) string {
	normalized := strings.ToUpper(strings.TrimSpace(streetType))
	if standardType, exists := streetTypeMap[normalized]; exists {
		return standardType
	}
	return normalized
}

func NormalizeStreetName(name string) string {
	// Remove extra spaces and convert to uppercase
	name = strings.TrimSpace(strings.ToUpper(name))

	// Remove accents
	name = removeAccents(name)

	// Replace multiple spaces with single space
	name = strings.Join(strings.Fields(name), " ")

	return name
}

func NormalizeNumber(number string) string {
	// Remove everything except digits
	var normalized strings.Builder
	for _, r := range number {
		if unicode.IsDigit(r) {
			normalized.WriteRune(r)
		}
	}
	return normalized.String()
}

func NormalizeCEP(cep string) string {
	// Remove everything except digits
	return NormalizeNumber(cep)
}

func removeAccents(s string) string {
	replacements := map[rune]rune{
		'Á': 'A', 'À': 'A', 'Ã': 'A', 'Â': 'A',
		'É': 'E', 'È': 'E', 'Ê': 'E',
		'Í': 'I', 'Ì': 'I', 'Î': 'I',
		'Ó': 'O', 'Ò': 'O', 'Õ': 'O', 'Ô': 'O',
		'Ú': 'U', 'Ù': 'U', 'Û': 'U',
		'Ç': 'C',
	}

	var result strings.Builder
	for _, r := range s {
		if replacement, ok := replacements[r]; ok {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func ExtractCEPPrefix(cep string) string {
	normalized := NormalizeCEP(cep)
	if len(normalized) >= 5 {
		return normalized[:5]
	}
	return normalized
}

// ValidateEvenOdd checks if a number matches the even/odd rule
func ValidateEvenOdd(number int, rule string) bool {
	switch rule {
	case "even":
		return number%2 == 0
	case "odd":
		return number%2 != 0
	case "all":
		return true
	default:
		return false
	}
}
