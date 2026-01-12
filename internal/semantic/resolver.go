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
	IsDefaultFilter bool
}

// Resolve menerjemahkan AST menjadi Intent secara STRICT.
// Jika ada kata / struktur tidak dikenali → ERROR.
func Resolve(ast parser.AST) (Intent, error) {
	intent := Intent{}

	// 1. UNKNOWN WORDS (STRICT)
	if len(ast.Unknown) > 0 {
		unknownStr := strings.Join(ast.Unknown, ", ")
		err := errors.NewUnknownWord(unknownStr)

		// Coba cari saran untuk kata pertama yang tidak dikenal
		if suggestion := FindSuggestion(ast.Unknown[0]); suggestion != "" {
			err = err.WithSuggestion(suggestion)
		}

		return intent, err
	}

	// 2. AKSI (WAJIB)
	if ast.Aksi == "" {
		return intent, errors.NewUnknownAction()
	}
	intent.Aksi = ast.Aksi

	// 3. OBJEK (WAJIB)
	if ast.Objek == "" {
		return intent, errors.NewUnknownObject()
	}
	// Mapping Objek Nahkoda -> Kubernetes Standard
	objekMap := map[string]string{
		"kru":       "pod",
		"mesin":     "node",
		"kapal":     "kapal", // special handling
		"jurnal":    "jurnal",
		"berita":    "berita",
		"geladak":   "namespace",
		"armada":    "deployment",
		"penjaga":   "daemonset",
		"pelabuhan": "service",
		"mercusuar": "ingress",
		"peta":      "configmap",
		"sandi":     "secret",
	}
	if mappedObj, ok := objekMap[ast.Objek]; ok {
		intent.Objek = mappedObj
	} else {
		intent.Objek = ast.Objek
	}

	// 4. LOKASI
	if ast.Lokasi != "" {
		intent.Lokasi = ast.Lokasi
	} else {
		// Default location logic
		if intent.Aksi == "liat" {
			intent.Lokasi = "semua geladak"
		} else {
			intent.Lokasi = "geladak default"
		}
	}

	// 5. TARGET (UNTUK CEK)
	intent.Target = ast.Target

	// 6. AKSI-SPECIFIC LOGIC
	switch intent.Aksi {

	case "liat":
		// Handle "liat berita"
		if intent.Objek == "berita" {
			// No filter needed for events usually, or maybe handled later
			intent.Filter = ""
			intent.IsDefaultFilter = false
			break // Exit case "liat"
		} else if intent.Objek == "kapal" {
			// Already handled in planner
			intent.Filter = ""
			intent.IsDefaultFilter = false
			break
		}

		if ast.Kondisi != "" {
			intent.Kondisi = ast.Kondisi

			filter, ok := ResolveCondition(ast.Kondisi)
			if !ok {
				return intent, errors.NewUnknownCondition(ast.Kondisi)
			}

			intent.Filter = filter
			intent.IsDefaultFilter = false
		} else {
			// default: hanya kru yang dianggap "sehat" (running) secara default
			// resource lain (node/mesin) ditampilkan apa adanya (tanpa filter)
			if intent.Objek == "pod" {
				intent.Filter = "status=Running"
				intent.IsDefaultFilter = true
			}
		}

	case "cek":
		if intent.Target == "" {
			return intent, errors.NewMissingTarget(intent.Objek)
		}

		// cek itu inspect 1 resource → tidak pakai filter
		intent.Filter = ""
		intent.IsDefaultFilter = false

	case "pindah":
		if intent.Objek != "kapal" {
			// Pindah hanya support kapal untuk saat ini
			return intent, errors.NewUnknownObject()
		}
		if intent.Target == "" {
			return intent, errors.NewMissingTarget(intent.Objek)
		}

	case "baca":
		if intent.Objek != "jurnal" {
			return intent, errors.NewUnknownObject()
		}
		if intent.Target == "" {
			return intent, errors.NewMissingTarget(intent.Objek)
		}
		// "baca jurnal" itu logs 1 pod -> no filter needed usually
		intent.Filter = ""
		intent.IsDefaultFilter = false

	case "masuk":
		// User bisa bilang "masuk [target]" (objek kosong, default kru) atau "masuk kru [target]"
		// Tapi Parser menaruh token non-keyword ke 'Target' hanya jika aksi=masuk & token terakhir.
		// Jika Objek kosong, kita set ke "kru".
		if intent.Objek == "" {
			intent.Objek = "kru"
		} else if intent.Objek != "kru" {
			// masuk hanya support kru (pod)
			return intent, errors.NewUnknownObject()
		}

		if intent.Target == "" {
			return intent, errors.NewMissingTarget("kru")
		}

	case "bikin":
		// bikin geladak (create ns) or bikin kru (run pod)
		allowedBikin := map[string]bool{
			"namespace":  true,
			"pod":        true,
			"deployment": true,
			"service":    true,
			"ingress":    true,
			"configmap":  true,
			"secret":     true,
		}
		if !allowedBikin[intent.Objek] {
			return intent, errors.NewUnknownObject()
		}
		if intent.Target == "" {
			return intent, errors.NewMissingTarget(intent.Objek)
		}

	case "pantau":
		// pantau mesin (top node) or pantau kru (top pod)
		if intent.Objek != "node" && intent.Objek != "pod" {
			return intent, errors.NewUnknownObject()
		}

	default:
		return intent, errors.NewUnknownAction()
	}

	return intent, nil
}
