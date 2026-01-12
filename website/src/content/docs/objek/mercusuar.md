---
title: Mercusuar (Ingress)
description: Pemandu trafik dari luar samudra menuju Pelabuhan.
---

**Mercusuar** adalah metafora untuk **Ingress** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Mercusuar? Karena Mercusuar berdiri di titik terluar kapal/pulau untuk memberikan arahan bagi kapal-kapal asing (trafik internet) agar tahu ke mana mereka harus menuju. Mercusuar mengatur rute berdasarkan nama domain atau jalur tertentu menuju Pelabuhan yang tepat.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Mercusuar** | **Ingress** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Mercusuar**:

- `nahkoda liat mercusuar` - Melihat semua gerbang masuk dari dunia luar.
- `nahkoda cek mercusuar [nama]` - Melihat aturan (rules) rute yang ada di mercusuar.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Mengonfigurasi sertifikat SSL (HTTPS) melalui CLI.
- Menambahkan rute (path) baru secara interaktif.
- Mengelola Ingress Controller.
