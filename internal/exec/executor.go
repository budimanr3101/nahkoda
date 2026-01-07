package exec

import (
	"fmt"
	"nahkoda/internal/planner"
)

func Execute(plan planner.Plan) error {
	fmt.Println("⚓ Rencana eksekusi:")
	fmt.Println("Operation :", plan.Operation)
	fmt.Println("Resource  :", plan.Resource)
	fmt.Println("Namespace :", plan.Namespace)

	fmt.Println("Filters:")
	if len(plan.Filters) == 0 {
		fmt.Println("  - none")
	} else {
		for k, v := range plan.Filters {
			// FIX UTAMA ADA DI SINI
			if len(v) > 0 && (v[0] == '!' || v[0] == '=') {
				// contoh: status!=Running
				fmt.Printf("  - %s%s\n", k, v)
			} else {
				// contoh: status=Running
				fmt.Printf("  - %s=%s\n", k, v)
			}
		}
	}

	fmt.Println()
	fmt.Println("(simulasi eksekusi, belum menyentuh Kubernetes)")
	return nil
}
