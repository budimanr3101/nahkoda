package semantic

import (
	"strings"

	"nahkoda/internal/errors"
	"nahkoda/internal/parser"
)

type Intent struct {
	Aksi            string
	Objek           string
	Lokasi          string
	Kondisi         string
	Filter          string
	Target          string
	Nilai           string
	SubTarget       string
	Follow          bool
	IsDefaultFilter bool
}

// IntentResolver interface is defined in resolver_interface.go

// Resolve menerjemahkan AST menjadi Intent secara STRICT.
// Jika ada kata / struktur tidak dikenali → ERROR.
func Resolve(ast parser.AST) (Intent, error) { // Sticking to original signature
	return resolveInternal(ast)
}

func resolveInternal(ast parser.AST) (Intent, error) {
	intent := Intent{}

	// 1. UNKNOWN WORDS (STRICT)
	if len(ast.Unknown) > 0 {
		unknownStr := strings.Join(ast.Unknown, ", ")
		err := errors.NewUnknownWord(unknownStr)

		// Coba cari saran untuk kata pertama yang tidak dikenal
		if len(ast.Unknown) > 0 {
			if suggestion := FindSuggestion(ast.Unknown[0]); suggestion != "" {
				err = err.WithSuggestion(suggestion)
			}
		}

		return intent, err
	}

	// 2. AKSI (WAJIB)
	if ast.Aksi == "" {
		return intent, errors.NewUnknownAction()
	}
	intent.Aksi = ast.Aksi

	// 3. MAPPING OBJEK
	// Note: Masuk action can have empty object initially (defaults to kru in resolver),
	// but generic object mapping should run first if object is present.
	if ast.Objek != "" {
		objekMap := map[string]string{
			"kru":        "pod",
			"mesin":      "node",
			"kapal":      "kapal",
			"jurnal":     "jurnal",
			"berita":     "berita",
			"geladak":    "namespace",
			"armada":     "deployment",
			"penjaga":    "daemonset",
			"pelabuhan":  "service",
			"mercusuar":  "ingress",
			"peta":       "configmap",
			"sandi":      "secret",
			"kesehatan":  "kesehatan",
			"perbekalan": "perbekalan",
		}
		if mappedObj, ok := objekMap[ast.Objek]; ok {
			intent.Objek = mappedObj
		} else {
			intent.Objek = ast.Objek
		}
	} else {
		// Special case for 'Masuk' checked in specific resolver
		if intent.Aksi == "masuk" {
			// Let resolver handle default
		} else {
			return intent, errors.NewUnknownObject()
		}
	}

	// 4. LOKASI
	if ast.Lokasi != "" {
		intent.Lokasi = ast.Lokasi
	} else {
		// Default location logic
		if intent.Aksi == "liat" && ast.Target == "" {
			intent.Lokasi = "semua geladak"
		} else {
			intent.Lokasi = "geladak default"
		}
	}

	// 5. TARGET, NILAI, SUBTARGET, FOLLOW
	// Validation Security (P0)
	if err := ValidateResourceName(ast.Target); err != nil {
		return intent, err
	}
	if err := ValidateResourceName(ast.SubTarget); err != nil {
		return intent, err
	}
	// Nilai is sometimes a resource name, sometimes numeric (scale 5). validator allows numbers.
	if err := ValidateResourceName(ast.Nilai); err != nil {
		return intent, err
	}

	intent.Target = ast.Target
	intent.Nilai = ast.Nilai
	intent.SubTarget = ast.SubTarget
	intent.Follow = ast.Follow

	// 6. ACTION-SPECIFIC LOGIC via Strategy Pattern
	resolvers := map[string]IntentResolver{
		"liat":   &LiatResolver{},
		"cek":    &CekResolver{},
		"hapus":  &HapusResolver{},
		"pindah": &PindahResolver{},
		"baca":   &BacaResolver{},
		"masuk":  &MasukResolver{},
		"bikin":  &BikinResolver{},
		"pantau": &PantauResolver{},
		"atur":   &AturResolver{},
		"tukar":  &TukarResolver{},
	}

	resolver, ok := resolvers[intent.Aksi]
	if !ok {
		return intent, errors.NewUnknownAction()
	}

	if err := resolver.Resolve(ast, &intent); err != nil {
		return intent, err
	}

	return intent, nil
}
