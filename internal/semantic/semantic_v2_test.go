package semantic

import (
	"testing"

	"nahkoda/internal/parser"
)

func TestResolve_Features_V0_10_0(t *testing.T) {
	tests := []struct {
		name    string
		input   parser.AST
		want    Intent
		wantErr bool
	}{
		// === BIKIN ===
		{
			name: "bikin geladak valid",
			input: parser.AST{
				Aksi:   "bikin",
				Objek:  "geladak",
				Target: "baru",
			},
			want: Intent{
				Aksi:   "bikin",
				Objek:  "geladak",
				Target: "baru",
				Lokasi: "geladak default",
			},
			wantErr: false,
		},
		{
			name: "bikin kru valid",
			input: parser.AST{
				Aksi:   "bikin",
				Objek:  "kru",
				Target: "nginx-pod",
			},
			want: Intent{
				Aksi:   "bikin",
				Objek:  "kru",
				Target: "nginx-pod",
				Lokasi: "geladak default",
			},
			wantErr: false,
		},
		{
			name: "bikin tanpa target error",
			input: parser.AST{
				Aksi:  "bikin",
				Objek: "geladak",
			},
			wantErr: true,
		},

		// === PANTAU ===
		{
			name: "pantau mesin valid",
			input: parser.AST{
				Aksi:  "pantau",
				Objek: "mesin",
			},
			want: Intent{
				Aksi:   "pantau",
				Objek:  "mesin",
				Lokasi: "geladak default",
			},
			wantErr: false,
		},
		{
			name: "pantau kru valid",
			input: parser.AST{
				Aksi:  "pantau",
				Objek: "kru",
			},
			want: Intent{
				Aksi:   "pantau",
				Objek:  "kru",
				Lokasi: "geladak default",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Resolve(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Aksi != tt.want.Aksi || got.Objek != tt.want.Objek {
					t.Errorf("Resolve() = %v, want %v", got, tt.want)
				}
				if got.Target != tt.want.Target {
					t.Errorf("Target mismatch: got %v, want %v", got.Target, tt.want.Target)
				}
			}
		})
	}
}
