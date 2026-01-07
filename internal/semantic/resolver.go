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
	// ===============================
	// 4️⃣ LOKASI
	// ===============================
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

			filter, ok := ResolveCondition(ast.Kondisi)
			if !ok {
				return intent, fmt.Errorf("kondisi tidak dikenali: %s", ast.Kondisi)
			}

			intent.Filter = filter
			intent.IsDefaultFilter = false
		} else {
			// default: hanya kru yang dianggap "sehat" (running) secara default
			// resource lain (node/mesin) ditampilkan apa adanya (tanpa filter)
			if intent.Objek == "kru" {
				intent.Filter = "status=Running"
				intent.IsDefaultFilter = true
			}
		}

	// ===============================
	// CEK → DESCRIBE (STRICT)
	// ===============================
	case "cek":
		if intent.Target == "" {
			return intent, fmt.Errorf("cek %s butuh nama %s", intent.Objek, intent.Objek)
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
