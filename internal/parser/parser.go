package parser

import (
	"strings"

	"nahkoda/internal/errors"
)

type AST struct {
	Aksi    string
	Objek   string
	Lokasi  string
	Kondisi string
	Target  string
	Unknown []string
}

func Parse(input string) (AST, error) {
	tokens := strings.Fields(strings.ToLower(input))
	ast := AST{}

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]

		switch tok {

		// ===== AKSI =====
		case "liat", "hapus", "cek", "pindah", "baca", "masuk", "bikin", "pantau":
			ast.Aksi = tok

		// ===== OBJEK =====
		case "kru", "mesin", "kapal", "jurnal", "berita":
			ast.Objek = tok

		// ===== KONDISI =====
		case "rusak", "bocor", "sehat", "terdampar", "siap", "mogok":
			ast.Kondisi = tok

		// ===== LOKASI =====
		case "di":
			if i+2 < len(tokens) && tokens[i+1] == "geladak" {
				ast.Lokasi = "geladak " + tokens[i+2]
				i += 2
			} else {
				ast.Unknown = append(ast.Unknown, tok)
			}

		// ===== DEFAULT =====
		default:
			// khusus untuk aksi "cek", "pindah", "baca", "masuk", token terakhir dianggap target
			// "liat" tidak butuh target biasanya, kecuali kustomisasi (tp "liat berita" gapake target)
			capturingActions := map[string]bool{
				"cek":    true,
				"pindah": true,
				"baca":   true,
				"masuk":  true,
				"bikin":  true,
			}
			if capturingActions[ast.Aksi] && ast.Target == "" {
				ast.Target = tok
			} else {
				ast.Unknown = append(ast.Unknown, tok)
			}
		}
	}

	// ===== VALIDASI MINIMAL =====
	if ast.Aksi == "" {
		return ast, errors.NewUnknownAction()
	}

	return ast, nil
}
