package semantic

import (
	"fmt"
	"strings"

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

	// ===============================
	// 1️⃣ UNKNOWN WORDS (STRICT)
	// ===============================
	if len(ast.Unknown) > 0 {
		return intent, fmt.Errorf(
			"kata tidak dikenali: %q",
			strings.Join(ast.Unknown, ", "),
		)
	}

	// ===============================
	// 2️⃣ AKSI (WAJIB)
	// ===============================
	if ast.Aksi == "" {
		return intent, fmt.Errorf("aksi tidak dikenali")
	}
	intent.Aksi = ast.Aksi

	// ===============================
	// 3️⃣ OBJEK (WAJIB)
	// ===============================
	if ast.Objek == "" {
		return intent, fmt.Errorf("objek tidak dikenali")
	}
	intent.Objek = ast.Objek

	// ===============================
	// 4️⃣ LOKASI
	// ===============================
	if ast.Lokasi != "" {
		intent.Lokasi = ast.Lokasi
	} else {
		intent.Lokasi = "semua geladak"
	}

	// ===============================
	// 5️⃣ TARGET (UNTUK CEK)
	// ===============================
	intent.Target = ast.Target

	// ===============================
	// 6️⃣ AKSI-SPECIFIC LOGIC
	// ===============================
	switch intent.Aksi {

	// ===============================
	// LIAT → LIST
	// ===============================
	case "liat":
		if ast.Kondisi != "" {
			intent.Kondisi = ast.Kondisi

			filter, ok := resolveCondition(ast.Kondisi)
			if !ok {
				return intent, fmt.Errorf("kondisi tidak dikenali: %s", ast.Kondisi)
			}

			intent.Filter = filter
			intent.IsDefaultFilter = false
		} else {
			// default: kru sehat
			intent.Filter = "status=Running"
			intent.IsDefaultFilter = true
		}

	// ===============================
	// CEK → DESCRIBE (STRICT)
	// ===============================
	case "cek":
		if intent.Target == "" {
			return intent, fmt.Errorf("cek kru butuh nama kru")
		}

		// cek itu inspect 1 resource → tidak pakai filter
		intent.Filter = ""
		intent.IsDefaultFilter = false

	// ===============================
	// AKSI TIDAK DIKENAL
	// ===============================
	default:
		return intent, fmt.Errorf("aksi tidak dikenali")
	}

	return intent, nil
}

// resolveCondition memetakan kondisi bahasa manusia ke filter Kubernetes.
func resolveCondition(kondisi string) (string, bool) {
	switch kondisi {
	case "rusak":
		return "status!=Running", true
	case "bocor":
		return "reason=OOMKilled", true
	case "sehat":
		return "status=Running", true
	default:
		return "", false
	}
}
