package semantic

import (
	"testing"

	"nahkoda/internal/errors"
	"nahkoda/internal/parser"
)

func TestResolve_ValidIntents(t *testing.T) {
	tests := []struct {
		name          string
		ast           parser.AST
		wantAksi      string
		wantObjek     string
		wantLokasi    string
		wantFilter    string
		wantIsDefault bool
		wantErr       bool
	}{
		{
			name: "liat kru - default sehat",
			ast: parser.AST{
				Aksi:  "liat",
				Objek: "kru",
			},
			wantAksi:      "liat",
			wantObjek:     "pod",
			wantLokasi:    "semua geladak",
			wantFilter:    "status=Running",
			wantIsDefault: true,
			wantErr:       false,
		},
		{
			name: "liat mesin - no default filter",
			ast: parser.AST{
				Aksi:  "liat",
				Objek: "mesin",
			},
			wantAksi:      "liat",
			wantObjek:     "node",
			wantLokasi:    "semua geladak",
			wantFilter:    "",
			wantIsDefault: false,
			wantErr:       false,
		},
		{
			name: "liat kru rusak",
			ast: parser.AST{
				Aksi:    "liat",
				Objek:   "kru",
				Kondisi: "rusak",
			},
			wantAksi:      "liat",
			wantObjek:     "pod",
			wantFilter:    "status!=Running",
			wantIsDefault: false,
			wantErr:       false,
		},
		{
			name: "liat kru with location",
			ast: parser.AST{
				Aksi:   "liat",
				Objek:  "kru",
				Lokasi: "geladak auth",
			},
			wantAksi:      "liat",
			wantObjek:     "pod",
			wantLokasi:    "geladak auth",
			wantFilter:    "status=Running",
			wantIsDefault: true, // kru still gets default filter
			wantErr:       false,
		},
		{
			name: "cek kru with target",
			ast: parser.AST{
				Aksi:   "cek",
				Objek:  "kru",
				Target: "my-pod",
			},
			wantAksi:   "cek",
			wantObjek:  "pod",
			wantLokasi: "geladak default",
			wantFilter: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent, err := Resolve(tt.ast)

			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if intent.Aksi != tt.wantAksi {
				t.Errorf("Aksi = %v, want %v", intent.Aksi, tt.wantAksi)
			}

			if intent.Objek != tt.wantObjek {
				t.Errorf("Objek = %v, want %v", intent.Objek, tt.wantObjek)
			}

			if tt.wantLokasi != "" && intent.Lokasi != tt.wantLokasi {
				t.Errorf("Lokasi = %v, want %v", intent.Lokasi, tt.wantLokasi)
			}

			if tt.wantFilter != "" && intent.Filter != tt.wantFilter {
				t.Errorf("Filter = %v, want %v", intent.Filter, tt.wantFilter)
			}

			if intent.IsDefaultFilter != tt.wantIsDefault {
				t.Errorf("IsDefaultFilter = %v, want %v", intent.IsDefaultFilter, tt.wantIsDefault)
			}
		})
	}
}

func TestResolve_Errors(t *testing.T) {
	tests := []struct {
		name        string
		ast         parser.AST
		wantErrType errors.ErrorType
	}{
		{
			name: "unknown word",
			ast: parser.AST{
				Aksi:    "liat",
				Objek:   "kru",
				Unknown: []string{"xyz"},
			},
			wantErrType: errors.ErrUnknownWord,
		},
		{
			name: "missing action",
			ast: parser.AST{
				Objek: "kru",
			},
			wantErrType: errors.ErrUnknownAction,
		},
		{
			name: "missing object",
			ast: parser.AST{
				Aksi: "liat",
			},
			wantErrType: errors.ErrUnknownObject,
		},
		{
			name: "unknown condition",
			ast: parser.AST{
				Aksi:    "liat",
				Objek:   "kru",
				Kondisi: "invalid",
			},
			wantErrType: errors.ErrUnknownCondition,
		},
		{
			name: "cek without target",
			ast: parser.AST{
				Aksi:  "cek",
				Objek: "kru",
			},
			wantErrType: errors.ErrMissingTarget,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Resolve(tt.ast)

			if err == nil {
				t.Error("Expected error but got nil")
				return
			}

			if ne, ok := err.(*errors.NahkodaError); ok {
				if !ne.IsType(tt.wantErrType) {
					t.Errorf("Error type = %v, want %v", ne.Type, tt.wantErrType)
				}
			} else {
				t.Error("Expected NahkodaError type")
			}
		})
	}
}

func TestResolveCondition(t *testing.T) {
	tests := []struct {
		kondisi string
		want    string
		wantOk  bool
	}{
		{"rusak", "status!=Running", true},
		{"sehat", "status=Running", true},
		{"terdampar", "status=Pending", true},
		{"siap", "status=Ready", true},
		{"mogok", "status=NotReady", true},
		{"bocor", "status.reason=OOMKilled", true},
		{"invalid", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.kondisi, func(t *testing.T) {
			got, ok := ResolveCondition(tt.kondisi)

			if ok != tt.wantOk {
				t.Errorf("ResolveCondition() ok = %v, want %v", ok, tt.wantOk)
			}

			if got != tt.want {
				t.Errorf("ResolveCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolve_DefaultLocationLogic(t *testing.T) {
	tests := []struct {
		name       string
		aksi       string
		wantLokasi string
	}{
		{
			name:       "liat defaults to all namespaces",
			aksi:       "liat",
			wantLokasi: "semua geladak",
		},
		{
			name:       "cek defaults to default namespace",
			aksi:       "cek",
			wantLokasi: "geladak default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast := parser.AST{
				Aksi:   tt.aksi,
				Objek:  "kru",
				Target: "test-pod",
			}

			intent, err := Resolve(ast)
			if err != nil {
				t.Errorf("Resolve() error = %v", err)
				return
			}

			if intent.Lokasi != tt.wantLokasi {
				t.Errorf("Lokasi = %v, want %v", intent.Lokasi, tt.wantLokasi)
			}
		})
	}
}
