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
	Nilai   string
	Unknown []string
}

func Parse(input string) (AST, error) {
	tokens := strings.Fields(strings.ToLower(input))
	ast := AST{}

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]

		switch tok {
		case "liat", "hapus", "cek", "pindah", "baca", "masuk", "bikin", "pantau", "atur":
			ast.Aksi = tok

		case "kru", "mesin", "kapal", "jurnal", "berita", "geladak", "armada", "penjaga", "pelabuhan", "mercusuar", "peta", "sandi", "kesehatan":
			ast.Objek = tok

		case "rusak", "bocor", "sehat", "terdampar", "siap", "mogok":
			ast.Kondisi = tok

		case "di":
			if i+2 < len(tokens) && tokens[i+1] == "geladak" {
				ast.Lokasi = "geladak " + tokens[i+2]
				i += 2
			} else {
				ast.Unknown = append(ast.Unknown, tok)
			}

		case "ke":
			if i+1 < len(tokens) {
				ast.Nilai = tokens[i+1]
				i += 1
			} else {
				ast.Unknown = append(ast.Unknown, tok)
			}

		default:
			// khusus untuk aksi "cek", "pindah", "baca", "masuk", token terakhir dianggap target
			capturingActions := map[string]bool{
				"cek":    true,
				"pindah": true,
				"baca":   true,
				"masuk":  true,
				"bikin":  true,
				"hapus":  true,
				"pantau": true,
				"atur":   true,
			}
			if capturingActions[ast.Aksi] && ast.Target == "" {
				ast.Target = tok
			} else {
				ast.Unknown = append(ast.Unknown, tok)
			}
		}
	}

	if ast.Aksi == "" {
		return ast, errors.NewUnknownAction()
	}

	return ast, nil
}
