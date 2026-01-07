package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"nahkoda/internal/planner"
)

func Execute(plan planner.Plan) error {
	// Construct args: kubectl [operation] [resource] [target]
	args := []string{plan.Operation, plan.Resource}

	if plan.Target != "" {
		args = append(args, plan.Target)
	}

	// Namespace
	if plan.Namespace == "all" {
		args = append(args, "-A")
	} else {
		args = append(args, "-n", plan.Namespace)
	}

	// Filter (field-selector)
	var selectors []string
	for k, v := range plan.Filters {
		selectors = append(selectors, k+v)
	}
	if len(selectors) > 0 {
		args = append(args, "--field-selector="+strings.Join(selectors, ","))
	}

	// Execute
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("⚓ Menjalankan: %s\n", strings.Join(cmd.Args, " "))
	return cmd.Run()
}
