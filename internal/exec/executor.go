package exec

import (
	"fmt"
	"nahkoda/internal/semantic"
)

func Execute(intent semantic.Intent) {
	fmt.Println("⚓ Nahkoda menerima perintah:")
	fmt.Printf("Aksi   : %s\n", intent.Aksi)
	fmt.Printf("Objek  : %s\n", intent.Objek)
	fmt.Printf("Lokasi : %s\n", intent.Lokasi)

	if intent.Kondisi != "" {
		fmt.Println("Kondisi:", intent.Kondisi)
	}

	if intent.Filter != "" {
		if intent.IsDefaultFilter {
			fmt.Println("Filter :", intent.Filter, "(aturan default: kru sehat)")
		} else {
			fmt.Println("Filter :", intent.Filter)
		}
	}

	fmt.Println("\n(simulasi, belum menyentuh Kubernetes)")
}
