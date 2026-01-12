---
title: Mesin (Node)
description: Mengenal Mesin, sumber tenaga utama bagi kapal Anda.
---

**Mesin** adalah metafora untuk **Node** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Mesin? Karena Node adalah server fisik atau virtual yang menyediakan tenaga (CPU & RAM) agar Kru bisa bekerja. Jika Mesin mogok, maka Kru di dalamnya tidak akan bisa melakukan apa-apa.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Mesin** | **Node** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Mesin**:

- `nahkoda liat mesin` - Melihat semua mesin yang tersedia di kapal.
- `nahkoda cek mesin [nama]` - Melakukan inspeksi mendalam terhadap kondisi mesin.
- `nahkoda pantau mesin` - Melihat penggunaan beban (CPU/RAM) pada setiap mesin.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Melakukan *taint* atau *label* pada Mesin.
- Mematikan atau melakukan *drain* pada Mesin secara langsung.
- Melihat log sistem level OS pada Mesin tersebut.
