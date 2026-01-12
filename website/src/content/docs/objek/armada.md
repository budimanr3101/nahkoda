---
title: Armada (Deployment)
description: Mengelola kelompok kru besar melalui Armada.
---

**Armada** adalah metafora untuk **Deployment** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Armada? Karena sebuah aplikasi biasanya terdiri dari banyak Kru yang identik. Armada memastikan jumlah Kru yang Anda inginkan selalu tersedia dan siap beroperasi. Jika satu Kru jatuh, Armada akan segera mengirim Kru pengganti.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Armada** | **Deployment** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Armada**:

- `nahkoda liat armada` - Melihat semua kelompok kerja yang sedang aktif.
- `nahkoda cek armada [nama]` - Melihat status replika dan strategi update.
- `nahkoda hapus armada [nama]` - Menarik seluruh kelompok kerja dari tugasnya.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Melakukan *Scaling* (menambah/mengurangi replika) melalui perintah teks (misal: `nahkoda atur armada ...`).
- Melakukan *Rollback* ke versi sebelumnya.
- Mengubah image container secara langsung melalui CLI Nahkoda.
