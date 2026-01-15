package parser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    AST
		wantErr bool
	}{
		{
			name:  "Basic Command",
			input: "liat kru",
			want: AST{
				Aksi:  "liat",
				Objek: "kru",
			},
		},
		{
			name:  "Command with Target",
			input: "cek kru podium-1",
			want: AST{
				Aksi:   "cek",
				Objek:  "kru",
				Target: "podium-1",
			},
		},
		{
			name:  "Command with Quoted Target",
			input: "cek kru 'podium satu'",
			want: AST{
				Aksi:   "cek",
				Objek:  "kru",
				Target: "podium satu",
			},
		},
		{
			name:  "Command with Double Quoted Target",
			input: "cek kru \"podium dua\"",
			want: AST{
				Aksi:   "cek",
				Objek:  "kru",
				Target: "podium dua",
			},
		},
		{
			name:  "Command with Location",
			input: "liat kru di geladak produksi",
			want: AST{
				Aksi:   "liat",
				Objek:  "kru",
				Lokasi: "geladak produksi",
			},
		},
		{
			name:  "Command with Condition",
			input: "liat kru rusak",
			want: AST{
				Aksi:    "liat",
				Objek:   "kru",
				Kondisi: "rusak",
			},
		},
		{
			name:  "Complex Command",
			input: "baca jurnal 'app backend' terus",
			want: AST{
				Aksi:   "baca",
				Objek:  "jurnal",
				Target: "app backend",
				Follow: true,
			},
		},
		{
			name:  "Scale Command",
			input: "atur armada frontend ke 5",
			want: AST{
				Aksi:   "atur",
				Objek:  "armada",
				Target: "frontend",
				Nilai:  "5",
			},
		},
		{
			name:    "Unknown Action",
			input:   "terbang ke bulan",
			want:    AST{Unknown: []string{"terbang"}, Nilai: "bulan"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{name: "Basic", input: `liat kru`, want: []string{"liat", "kru"}},
		{name: "Quotes", input: `liat "kru detail"`, want: []string{"liat", "kru detail"}},
		{name: "Single Quotes", input: `liat 'kru detail'`, want: []string{"liat", "kru detail"}},
		{name: "Flags", input: `baca jurnal -f`, want: []string{"baca", "jurnal", "-f"}},
		{name: "Spaces", input: `   banyak    spasi   `, want: []string{"banyak", "spasi"}},
		{name: "Unclosed Quote", input: `liat "kru`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
