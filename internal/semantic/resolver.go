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

	// ===== DEFAULT OBJEK =====
	if intent.Objek == "" {
		intent.Objek = "kru"
	}

	// ===== LOKASI =====
	if ast.Lokasi != "" {
		// parser sekarang hanya memberi nama namespace: "auth"
		// semantic bertugas membentuk makna bahasa
		intent.Lokasi = "geladak " + ast.Lokasi
	} else {
		intent.Lokasi = "semua geladak"
	}

	// ===== KONDISI =====
	if ast.Kondisi != "" {
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
