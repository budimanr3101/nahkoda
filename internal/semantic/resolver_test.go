package semantic_test

import (
	"nahkoda/internal/parser"
	"nahkoda/internal/semantic"
	"testing"
)

func TestResolveIntent(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect semantic.Intent
	}{
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
		},
		{
			name:  "kondisi sehat (dulu tidak dikenal, sekarang default explicit)",
			input: "liat kru sehat",
			expect: semantic.Intent{
				Aksi:            "liat",
				Objek:           "kru",
				Lokasi:          "semua geladak",
				Kondisi:         "sehat",
				Filter:          "status=Running",
				IsDefaultFilter: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast := parser.Parse(tt.input)
			intent := semantic.Resolve(ast)

			if intent != tt.expect {
				t.Errorf(
					"\ninput : %q\nexpect: %#v\ngot   : %#v",
					tt.input,
					tt.expect,
					intent,
				)
			}
		})
	}
}
