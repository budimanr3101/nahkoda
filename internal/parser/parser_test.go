package parser

import (
	"reflect"
	"testing"

	"nahkoda/internal/errors"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    AST
		wantErr bool
		errType errors.ErrorType
	}{
		{
			name:  "Basic list pod",
			input: "liat kru",
			want: AST{
				Aksi:  "liat",
				Objek: "kru",
			},
			wantErr: false,
		},
		{
			name:  "List pod with condition",
			input: "liat kru rusak",
			want: AST{
				Aksi:    "liat",
				Objek:   "kru",
				Kondisi: "rusak",
			},
			wantErr: false,
		},
		{
			name:  "List pod with location",
			input: "liat kru di geladak auth",
			want: AST{
				Aksi:   "liat",
				Objek:  "kru",
				Lokasi: "geladak auth",
			},
			wantErr: false,
		},
		{
			name:  "Full command",
			input: "liat kru bocor di geladak payment",
			want: AST{
				Aksi:    "liat",
				Objek:   "kru",
				Kondisi: "bocor",
				Lokasi:  "geladak payment",
			},
			wantErr: false,
		},
		{
			name:  "Case insensitive",
			input: "LIAT KRU RUSAK",
			want: AST{
				Aksi:    "liat",
				Objek:   "kru",
				Kondisi: "rusak",
			},
			wantErr: false,
		},
		{
			name:  "Extra spaces",
			input: "  liat   kru    rusak  ",
			want: AST{
				Aksi:    "liat",
				Objek:   "kru",
				Kondisi: "rusak",
			},
			wantErr: false,
		},
		{
			name:  "Cek specific target",
			input: "cek kru pod-123",
			want: AST{
				Aksi:   "cek",
				Objek:  "kru",
				Target: "pod-123",
			},
			wantErr: false,
		},
		{
			name:  "Unknown words",
			input: "liat kru xyz",
			want: AST{
				Aksi:    "liat",
				Objek:   "kru",
				Unknown: []string{"xyz"},
			},
			wantErr: false,
		},
		{
			name:    "Empty input",
			input:   "",
			want:    AST{},
			wantErr: true,
			errType: errors.ErrUnknownAction, // Empty input results in empty action string
		},
		{
			name:    "Missing action",
			input:   "kru rusak", // "kru" treated as object, no action found
			want:    AST{Objek: "kru", Kondisi: "rusak"},
			wantErr: true,
			errType: errors.ErrUnknownAction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				// Check error type if specified
				if ne, ok := err.(*errors.NahkodaError); ok {
					if !ne.IsType(tt.errType) {
						t.Errorf("Error type = %v, want %v", ne.Type, tt.errType)
					}
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
