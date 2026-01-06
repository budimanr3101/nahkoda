package semantic

import "nahkoda/internal/parser"

type Intent struct {
	Aksi            string
	Objek           string
	Lokasi          string
	Kondisi         string
	Filter          string
	IsDefaultFilter bool
}

func Resolve(ast parser.AST) Intent {
	intent := Intent{
		Aksi:  ast.Aksi,
		Objek: ast.Objek,
	}

	// default objek
	if intent.Objek == "" {
		intent.Objek = "kru"
	}

	// lokasi
	if ast.Lokasi != "" {
		intent.Lokasi = ast.Lokasi
	} else {
		intent.Lokasi = "semua geladak"
	}

	// ===== KONDISI =====
	if ast.Kondisi != "" {
		// kondisi eksplisit
		intent.Kondisi = ast.Kondisi

		if filter, ok := ResolveCondition(ast.Kondisi); ok {
			intent.Filter = filter
			intent.IsDefaultFilter = false
		}
	} else {
		// default: kru sehat
		intent.Filter = "status=Running"
		intent.IsDefaultFilter = true
	}

	return intent
}
