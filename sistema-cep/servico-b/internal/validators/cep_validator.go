package validators

import "regexp"

// ValidateCEP valida se o CEP tem o formato correto
func ValidateCEP(cep string) bool {
	// Verifica se é string e tem exatamente 8 dígitos
	if len(cep) != 8 {
		return false
	}

	// Verifica se contém apenas dígitos
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}
