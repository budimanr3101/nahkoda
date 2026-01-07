package parser

import (
	"fmt"
	"strings"
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
		case "liat", "hapus", "cek":
			ast.Aksi = tok

		// ===== OBJEK =====
		case "kru":
			ast.Objek = tok

		// ===== KONDISI =====
		case "rusak", "bocor", "sehat":
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
			// khusus untuk aksi "cek", token terakhir dianggap target
			if ast.Aksi == "cek" && i == len(tokens)-1 {
				ast.Target = tok
			} else {
				ast.Unknown = append(ast.Unknown, tok)
			}
		}
	}

	// ===== VALIDASI MINIMAL =====
	if ast.Aksi == "" {
		return ast, fmt.Errorf("aksi tidak dikenali")
	}

	return ast, nil
}
