package parser

import (
	"testing"

	"nahkoda/internal/errors"
)

func TestParse_ValidCommands(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantAksi    string
		wantObjek   string
		wantLokasi  string
		wantKondisi string
		wantTarget  string
		wantErr     bool
	}{
		{
			name:      "simple liat kru",
			input:     "liat kru",
			wantAksi:  "liat",
			wantObjek: "kru",
			wantErr:   false,
		},
		{
			name:        "liat kru with condition",
			input:       "liat kru rusak",
			wantAksi:    "liat",
			wantObjek:   "kru",
			wantKondisi: "rusak",
			wantErr:     false,
		},
		{
			name:       "liat kru with location",
			input:      "liat kru di geladak auth",
			wantAksi:   "liat",
			wantObjek:  "kru",
			wantLokasi: "geladak auth",
			wantErr:    false,
		},
		{
			name:        "liat kru with condition and location",
			input:       "liat kru rusak di geladak payment",
			wantAksi:    "liat",
			wantObjek:   "kru",
			wantKondisi: "rusak",
			wantLokasi:  "geladak payment",
			wantErr:     false,
		},
		{
			name:       "cek kru with target",
			input:      "cek kru my-pod-123",
			wantAksi:   "cek",
			wantObjek:  "kru",
			wantTarget: "my-pod-123",
			wantErr:    false,
		},
		{
			name:      "liat mesin",
			input:     "liat mesin",
			wantAksi:  "liat",
			wantObjek: "mesin",
			wantErr:   false,
		},
		{
			name:        "liat mesin siap",
			input:       "liat mesin siap",
			wantAksi:    "liat",
			wantObjek:   "mesin",
			wantKondisi: "siap",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if ast.Aksi != tt.wantAksi {
				t.Errorf("Aksi = %v, want %v", ast.Aksi, tt.wantAksi)
			}

			if ast.Objek != tt.wantObjek {
				t.Errorf("Objek = %v, want %v", ast.Objek, tt.wantObjek)
			}

			if tt.wantLokasi != "" && ast.Lokasi != tt.wantLokasi {
				t.Errorf("Lokasi = %v, want %v", ast.Lokasi, tt.wantLokasi)
			}

			if tt.wantKondisi != "" && ast.Kondisi != tt.wantKondisi {
				t.Errorf("Kondisi = %v, want %v", ast.Kondisi, tt.wantKondisi)
			}

			if tt.wantTarget != "" && ast.Target != tt.wantTarget {
				t.Errorf("Target = %v, want %v", ast.Target, tt.wantTarget)
			}
		})
	}
}

func TestParse_UnknownWords(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantUnknown []string
	}{
		{
			name:        "unknown word xyz",
			input:       "liat kru xyz",
			wantUnknown: []string{"xyz"},
		},
		{
			name:        "multiple unknown words",
			input:       "liat kru foo bar",
			wantUnknown: []string{"foo", "bar"},
		},
		{
			name:        "invalid di without geladak",
			input:       "liat kru di auth",
			wantUnknown: []string{"di", "auth"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, _ := Parse(tt.input)

			if len(ast.Unknown) != len(tt.wantUnknown) {
				t.Errorf("Unknown count = %v, want %v", len(ast.Unknown), len(tt.wantUnknown))
				return
			}

			for i, unknown := range tt.wantUnknown {
				if ast.Unknown[i] != unknown {
					t.Errorf("Unknown[%d] = %v, want %v", i, ast.Unknown[i], unknown)
				}
			}
		})
	}
}

func TestParse_MissingAction(t *testing.T) {
	input := "kru rusak"
	_, err := Parse(input)

	if err == nil {
		t.Error("Expected error for missing action")
		return
	}

	if ne, ok := err.(*errors.NahkodaError); ok {
		if !ne.IsType(errors.ErrUnknownAction) {
			t.Errorf("Expected ErrUnknownAction, got %v", ne.Type)
		}
	} else {
		t.Error("Expected NahkodaError type")
	}
}

func TestParse_CaseInsensitive(t *testing.T) {
	tests := []string{
		"LIAT KRU",
		"Liat Kru",
		"lIaT kRu",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			ast, err := Parse(input)

			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if ast.Aksi != "liat" || ast.Objek != "kru" {
				t.Errorf("Case insensitive failed: aksi=%v, objek=%v", ast.Aksi, ast.Objek)
			}
		})
	}
}

func TestParse_AllConditions(t *testing.T) {
	conditions := []string{"rusak", "bocor", "sehat", "terdampar", "siap", "mogok"}

	for _, cond := range conditions {
		t.Run(cond, func(t *testing.T) {
			input := "liat kru " + cond
			ast, err := Parse(input)

			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if ast.Kondisi != cond {
				t.Errorf("Kondisi = %v, want %v", ast.Kondisi, cond)
			}
		})
	}
}

func TestParse_EmptyInput(t *testing.T) {
	_, err := Parse("")

	if err == nil {
		t.Error("Expected error for empty input")
	}
}
