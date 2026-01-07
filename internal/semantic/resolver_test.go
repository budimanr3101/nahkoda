package semantic_test

import (
	"testing"

	"nahkoda/internal/parser"
	"nahkoda/internal/semantic"
)

func TestResolveIntent_Strict(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expect    semantic.Intent
		expectErr bool
	}{
		// =============================
		// ✅ VALID CASES
		// =============================
		{
			name:  "liat kru default",
			input: "liat kru",
			expect: semantic.Intent{
				Aksi:            "liat",
				Objek:           "kru",
				Lokasi:          "semua geladak",
				Filter:          "status=Running",
				IsDefaultFilter: true,
			},
			expectErr: false,
		},
		{
			name:  "liat kru di geladak auth",
			input: "liat kru di geladak auth",
			expect: semantic.Intent{
				Aksi:            "liat",
				Objek:           "kru",
				Lokasi:          "geladak auth",
				Filter:          "status=Running",
				IsDefaultFilter: true,
			},
			expectErr: false,
		},
		{
			name:  "liat kru rusak",
			input: "liat kru rusak",
			expect: semantic.Intent{
				Aksi:            "liat",
				Objek:           "kru",
				Lokasi:          "semua geladak",
				Kondisi:         "rusak",
				Filter:          "status!=Running",
				IsDefaultFilter: false,
			},
			expectErr: false,
		},
		{
			name:  "liat kru rusak di geladak auth",
			input: "liat kru rusak di geladak auth",
			expect: semantic.Intent{
				Aksi:            "liat",
				Objek:           "kru",
				Lokasi:          "geladak auth",
				Kondisi:         "rusak",
				Filter:          "status!=Running",
				IsDefaultFilter: false,
			},
			expectErr: false,
		},
		{
			name:  "liat kru bocor",
			input: "liat kru bocor",
			expect: semantic.Intent{
				Aksi:            "liat",
				Objek:           "kru",
				Lokasi:          "semua geladak",
				Kondisi:         "bocor",
				Filter:          "reason=OOMKilled",
				IsDefaultFilter: false,
			},
			expectErr: false,
		},
		{
			name:  "liat kru sehat explicit",
			input: "liat kru sehat",
			expect: semantic.Intent{
				Aksi:            "liat",
				Objek:           "kru",
				Lokasi:          "semua geladak",
				Kondisi:         "sehat",
				Filter:          "status=Running",
				IsDefaultFilter: false,
			},
			expectErr: false,
		},

		// =============================
		// ❌ INVALID / STRICT FAILS
		// =============================
		{
			name:      "kata tidak dikenali",
			input:     "liat kru xyz",
			expectErr: true,
		},
		{
			name:      "aksi tidak dikenali",
			input:     "terbangkan kapal",
			expectErr: true,
		},
		{
			name:      "kondisi tidak dikenali",
			input:     "liat kru aneh",
			expectErr: true,
		},
		{
			name:      "lokasi rusak struktur",
			input:     "liat kru di auth",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := parser.Parse(tt.input)
			if err != nil {
				if tt.expectErr {
					return
				}
				t.Fatalf("unexpected parser error: %v", err)
			}

			intent, err := semantic.Resolve(ast)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got success: %#v", intent)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if intent != tt.expect {
				t.Errorf(
					"\nINPUT : %q\nEXPECT: %#v\nGOT   : %#v",
					tt.input,
					tt.expect,
					intent,
				)
			}
		})
	}
}
