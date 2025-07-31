package validators

import "testing"

func TestValidateCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{
			name:     "CEP válido com 8 dígitos",
			cep:      "29902555",
			expected: true,
		},
		{
			name:     "CEP válido com zeros",
			cep:      "00000000",
			expected: true,
		},
		{
			name:     "CEP inválido - menos de 8 dígitos",
			cep:      "1234567",
			expected: false,
		},
		{
			name:     "CEP inválido - mais de 8 dígitos",
			cep:      "123456789",
			expected: false,
		},
		{
			name:     "CEP inválido - contém letras",
			cep:      "1234567a",
			expected: false,
		},
		{
			name:     "CEP inválido - contém caracteres especiais",
			cep:      "12345-67",
			expected: false,
		},
		{
			name:     "CEP inválido - string vazia",
			cep:      "",
			expected: false,
		},
		{
			name:     "CEP inválido - apenas espaços",
			cep:      "        ",
			expected: false,
		},
		{
			name:     "CEP inválido - com espaços no meio",
			cep:      "123 567 8",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCEP(tt.cep)
			if result != tt.expected {
				t.Errorf("ValidateCEP(%q) = %v, expected %v", tt.cep, result, tt.expected)
			}
		})
	}
}
