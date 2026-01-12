package parser

import (
	"strings"

	"nahkoda/internal/errors"
)

type AST struct {
	Aksi      string
	Objek     string
	Lokasi    string
	Kondisi   string
	Target    string
	Nilai     string
	SubTarget string
	Follow    bool
	Unknown   []string
}

func Parse(input string) (AST, error) {
	tokens := strings.Fields(strings.ToLower(input))
	ast := AST{}

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]

		switch tok {
		case "liat", "hapus", "cek", "pindah", "baca", "masuk", "bikin", "pantau", "atur", "tukar":
			ast.Aksi = tok

		case "kru", "mesin", "kapal", "jurnal", "berita", "geladak", "armada", "penjaga", "pelabuhan", "mercusuar", "peta", "sandi", "kesehatan", "perbekalan":
			if ast.Objek == "perbekalan" {
				// perbekalan [objek] -> capture the [objek] as a special sub-object
				// here we can reuse SubTarget or Nilai, but let's add a new field if needed
				// For now, let's use ast.Target as the resource type if it's empty
				// actually, let's just use Nilai to store the resource type for perbekalan
				ast.Nilai = tok
			} else {
				ast.Objek = tok
			}

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

		case "terus":
			ast.Follow = true

		case "kabin":
			if i+1 < len(tokens) {
				ast.SubTarget = tokens[i+1]
				i += 1
			} else {
				ast.Unknown = append(ast.Unknown, tok)
			}

		default:
			// khusus untuk aksi "cek", "pindah", "baca", "masuk", token terakhir dianggap target
			capturingActions := map[string]bool{
				"liat":   true,
				"cek":    true,
				"pindah": true,
				"baca":   true,
				"masuk":  true,
				"bikin":  true,
				"hapus":  true,
				"pantau": true,
				"atur":   true,
				"tukar":  true,
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
