package parser

import (
	"errors"
	"fmt"
	"strings"
)

type AST struct {
	Aksi    string
	Objek   string
	Lokasi  string
	Kondisi string
	Filter  string
}

func Parse(input string) (AST, error) {
	// v0 parser: super naif, tapi sesuai SPEC
	tokens := split(input)
	ast := AST{}

	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {

		case "liat":
			ast.Aksi = "liat"

		case "kru":
			ast.Objek = "kru"

		case "rusak", "terdampar", "bocor", "sehat":
			ast.Kondisi = tokens[i]

		case "di":
			if i+2 < len(tokens) && tokens[i+1] == "geladak" {
				ast.Lokasi = "geladak " + tokens[i+2]
			}
			fmt.Printf("DEBUG TOKENS: %#v\n", tokens)
			fmt.Printf("DEBUG AST: %#v\n", ast)
		}
	}
	if ast.Aksi == "" {
		return ast, errors.New("aksi tidak dikenali")
	}
	return ast, nil
}

func split(s string) []string {
	return strings.Fields(strings.ToLower(s))
}
