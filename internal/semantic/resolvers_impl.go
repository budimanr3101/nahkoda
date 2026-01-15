package semantic

import (
	"nahkoda/internal/errors"
	"nahkoda/internal/parser"
)

// --- LIAT ---
type LiatResolver struct{}

func (r *LiatResolver) Resolve(ast parser.AST, intent *Intent) error {
	// Handle special objects
	if intent.Objek == "berita" || intent.Objek == "perbekalan" {
		intent.Filter = ""
		intent.IsDefaultFilter = false
		return nil
	}

	if ast.Kondisi != "" {
		intent.Kondisi = ast.Kondisi
		filter, ok := ResolveCondition(ast.Kondisi)
		if !ok {
			return errors.NewUnknownCondition(ast.Kondisi)
		}
		intent.Filter = filter
		intent.IsDefaultFilter = false
	} else {
		// Default filter
		if intent.Objek == "pod" {
			intent.Filter = "status=Running"
			intent.IsDefaultFilter = true
		}
	}
	return nil
}

// --- CEK ---
type CekResolver struct{}

func (r *CekResolver) Resolve(ast parser.AST, intent *Intent) error {
	if intent.Objek == "kesehatan" {
		return nil
	}
	if intent.Target == "" {
		return errors.NewMissingTarget(intent.Objek)
	}
	return nil
}

// --- HAPUS ---
type HapusResolver struct{}

func (r *HapusResolver) Resolve(ast parser.AST, intent *Intent) error {
	if intent.Target == "" {
		return errors.NewMissingTarget(intent.Objek)
	}
	intent.Filter = ""
	intent.IsDefaultFilter = false
	return nil
}

// --- PINDAH ---
type PindahResolver struct{}

func (r *PindahResolver) Resolve(ast parser.AST, intent *Intent) error {
	if intent.Objek != "kapal" {
		return errors.NewUnknownObject()
	}
	if intent.Target == "" {
		return errors.NewMissingTarget(intent.Objek)
	}
	return nil
}

// --- BACA ---
type BacaResolver struct{}

func (r *BacaResolver) Resolve(ast parser.AST, intent *Intent) error {
	if intent.Objek != "jurnal" {
		return errors.NewUnknownObject()
	}
	if intent.Target == "" {
		return errors.NewMissingTarget(intent.Objek)
	}
	intent.Filter = ""
	intent.IsDefaultFilter = false
	return nil
}

// --- MASUK ---
type MasukResolver struct{}

func (r *MasukResolver) Resolve(ast parser.AST, intent *Intent) error {
	// Logic from original Resolver:
	// If Objek is empty, default to "kru".
	// But `intent.Objek` is already set by generic logic before calling generic resolver?
	// Checking the original code:
	// > if intent.Objek == "" { intent.Objek = "kru" }
	// This mutation should ideally happen before, or here.
	// But `intent.Objek` comes from `ast.Objek` being mapped.
	// In the original, it was modifying `intent` based on specific logic.

	if intent.Objek == "" {
		intent.Objek = "kru"
	} else if intent.Objek != "kru" {
		return errors.NewUnknownObject()
	}

	if intent.Target == "" {
		return errors.NewMissingTarget("kru")
	}
	return nil
}

// --- BIKIN ---
type BikinResolver struct{}

func (r *BikinResolver) Resolve(ast parser.AST, intent *Intent) error {
	allowedBikin := map[string]bool{
		"namespace":  true,
		"pod":        true,
		"deployment": true,
		"service":    true,
		"ingress":    true,
		"configmap":  true,
		"secret":     true,
		"perbekalan": true,
	}
	if !allowedBikin[intent.Objek] {
		return errors.NewUnknownObject()
	}
	if intent.Target == "" {
		return errors.NewMissingTarget(intent.Objek)
	}
	return nil
}

// --- PANTAU ---
type PantauResolver struct{}

func (r *PantauResolver) Resolve(ast parser.AST, intent *Intent) error {
	if intent.Objek != "node" && intent.Objek != "pod" {
		return errors.NewUnknownObject()
	}
	return nil
}

// --- ATUR ---
type AturResolver struct{}

func (r *AturResolver) Resolve(ast parser.AST, intent *Intent) error {
	if intent.Objek != "deployment" {
		return errors.NewUnknownObject()
	}
	if intent.Target == "" {
		return errors.NewMissingTarget(intent.Objek)
	}
	if intent.Nilai == "" {
		return errors.New(errors.ErrInvalidSyntax, "jumlah replika harus ditentukan (contoh: ke 5)")
	}
	return nil
}

// --- TUKAR ---
type TukarResolver struct{}

func (r *TukarResolver) Resolve(ast parser.AST, intent *Intent) error {
	if intent.Objek != "deployment" && intent.Objek != "daemonset" {
		return errors.NewUnknownObject()
	}
	if intent.Target == "" {
		return errors.NewMissingTarget(intent.Objek)
	}
	return nil
}
