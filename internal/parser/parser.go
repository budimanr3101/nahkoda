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
	// 1. Tokenize input (handles quotes)
	tokens, err := tokenize(input)
	if err != nil {
		return AST{}, err
	}
	ast := AST{}

	if len(tokens) == 0 {
		return ast, nil
	}

	// 2. Iterate tokens
	for i := 0; i < len(tokens); i++ {
		// Normalize token for keyword matching (lowercase), but keep original for values
		rawTok := tokens[i]
		lowerTok := strings.ToLower(rawTok)

		switch lowerTok {
		// --- ACTIONS ---
		case "liat", "hapus", "cek", "pindah", "baca", "masuk", "bikin", "pantau", "atur", "tukar":
			ast.Aksi = lowerTok // Actions are always lowercase keywords

		// --- OBJECTS ---
		case "kru", "mesin", "kapal", "jurnal", "berita", "geladak", "armada", "penjaga", "pelabuhan", "mercusuar", "peta", "sandi", "kesehatan", "perbekalan":
			if ast.Objek == "perbekalan" {
				// Special case: perbekalan [resource_type]
				// Treat the next token as the resource type (stored in Nilai for consistency)
				ast.Nilai = rawTok // usage: "perbekalan kru" -> Objek=perbekalan, Nilai=kru
			} else {
				ast.Objek = lowerTok
			}

		// --- CONDITIONS ---
		case "rusak", "bocor", "sehat", "terdampar", "siap", "mogok":
			ast.Kondisi = lowerTok

		// --- PREPOSITIONS/KEYWORDS ---
		case "di":
			// di [geladak] [NAME]
			if i+1 < len(tokens) && strings.ToLower(tokens[i+1]) == "geladak" {
				if i+2 < len(tokens) {
					ast.Lokasi = "geladak " + tokens[i+2]
					i += 2
				} else {
					// "di geladak" without name -> maybe implicit? or error?
					// For now let's just consume "geladak"
					i += 1
				}
			} else {
				// "di" without "geladak"? Treat as unknown or part of name if not handled?
				// For strictness, let's treat "di" as unknown if not followed by geladak
				// OR, we can try to be smart. But per rules, stick to structure.
				ast.Unknown = append(ast.Unknown, rawTok)
			}

		case "ke":
			if i+1 < len(tokens) {
				ast.Nilai = tokens[i+1]
				i += 1
			} else {
				ast.Unknown = append(ast.Unknown, rawTok)
			}

		case "terus":
			ast.Follow = true

		case "kabin":
			if i+1 < len(tokens) {
				ast.SubTarget = tokens[i+1]
				i += 1
			} else {
				ast.Unknown = append(ast.Unknown, rawTok)
			}

		default:
			// --- DEFAULT / UNKNOWN ---
			// Check if this should be a Target

			// List of actions that typically accept a target name
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

			// Capture as Target if:
			// 1. We have a valid action that needs target
			// 2. Target is currently empty
			// 3. It's not a keyword (already checked by switch cases)
			if capturingActions[ast.Aksi] && ast.Target == "" {
				ast.Target = rawTok // Preserve original case (e.g. "MyPod")
			} else {
				ast.Unknown = append(ast.Unknown, rawTok)
			}
		}
	}

	if ast.Aksi == "" {
		return ast, errors.NewUnknownAction()
	}

	return ast, nil
}
