// package main

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"nahkoda/internal/parser"
// 	"nahkoda/internal/semantic"
// )

// func main() {
// 	if len(os.Args) < 2 {
// 		printHelp()
// 		return
// 	}

// 	input := strings.Join(os.Args[1:], " ")

// 	ast, err := parser.Parse(input)
// 	if err != nil {
// 		fmt.Println("❌", err.Error())
// 		return
// 	}
// 	intent := semantic.Resolve(ast)

// 	printIntent(intent)
// }

// func printHelp() {
// 	fmt.Println("⚓ Nahkoda — Natural CLI for Kubernetes\n")
// 	fmt.Println("Contoh penggunaan:")
// 	fmt.Println("  nahkoda liat kru")
// 	fmt.Println("  nahkoda liat kru rusak")
// 	fmt.Println("  nahkoda liat kru di geladak auth")
// 	fmt.Println("  nahkoda liat kru bocor di geladak auth")
// }

// func printIntent(intent semantic.Intent) {
// 	fmt.Println("⚓ Nahkoda menerima perintah:")

// 	fmt.Printf("Aksi   : %s\n", intent.Aksi)
// 	fmt.Printf("Objek  : %s\n", intent.Objek)
// 	fmt.Printf("Lokasi : %s\n", intent.Lokasi)

// 	if intent.Kondisi != "" {
// 		fmt.Printf("Kondisi: %s\n", intent.Kondisi)
// 	}

// 	if intent.Filter != "" {
// 		if intent.IsDefaultFilter {
// 			fmt.Printf("Filter : %s (aturan default: kru sehat)\n", intent.Filter)
// 		} else {
// 			fmt.Printf("Filter : %s\n", intent.Filter)
// 		}
// 	}

// 	fmt.Println("\n(simulasi, belum menyentuh Kubernetes)")
// }

package main

import "nahkoda/cmd"

func main() {
	cmd.Execute()
}
