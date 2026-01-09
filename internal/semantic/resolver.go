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

	// ===============================
	// 1️⃣ UNKNOWN WORDS (STRICT)
	// ===============================
	if len(ast.Unknown) > 0 {
		return intent, errors.NewUnknownWord(strings.Join(ast.Unknown, ", "))
	}

	// ===============================
	// 2️⃣ AKSI (WAJIB)
	// ===============================
	if ast.Aksi == "" {
		return intent, errors.NewUnknownAction()
	}
	intent.Aksi = ast.Aksi

	// ===============================
	// 3️⃣ OBJEK (WAJIB)
	// ===============================
	if ast.Objek == "" {
		return intent, errors.NewUnknownObject()
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
			return intent, errors.NewMissingTarget(intent.Objek)
		}

		// cek itu inspect 1 resource → tidak pakai filter
		intent.Filter = ""
		intent.IsDefaultFilter = false

	// ===============================
	// PINDAH → USE CONTEXT
	// ===============================
	case "pindah":
		if intent.Objek != "kapal" {
			// Pindah hanya support kapal untuk saat ini
			return intent, errors.NewUnknownObject()
		}
		if intent.Target == "" {
			return intent, errors.NewMissingTarget(intent.Objek)
		}

	// ===============================
	// BACA → LOGS
	// ===============================
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

	// ===============================
	// MASUK → EXEC
	// ===============================
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

	// ===============================
	// BIKIN → CREATE / RUN
	// ===============================
	case "bikin":
		// bikin geladak (create ns) or bikin kru (run pod)
		if intent.Objek != "geladak" && intent.Objek != "kru" {
			return intent, errors.NewUnknownObject()
		}
		if intent.Target == "" {
			return intent, errors.NewMissingTarget(intent.Objek)
		}

	// ===============================
	// PANTAU → TOP
	// ===============================
	case "pantau":
		// pantau mesin (top node) or pantau kru (top pod)
		if intent.Objek != "mesin" && intent.Objek != "kru" {
			return intent, errors.NewUnknownObject()
		}

	// ===============================
	// AKSI TIDAK DIKENAL
	// ===============================
	default:
		return intent, errors.NewUnknownAction()
	}

	return intent, nil
}
