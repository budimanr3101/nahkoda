---
title: Geladak (Namespace)
description: Mengatur area kerja dan isolasi Kru di atas Geladak.
---

**Geladak** adalah metafora untuk **Namespace** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Geladak? Karena dalam kapal besar, tidak semua kru bekerja di tempat yang sama. Ada geladak mesin, geladak hunian, dan geladak navigasi. Geladak memisahkan area kerja agar tidak terjadi tabrakan antar tim dan mempermudah pengorganisasian kapal yang sangat luas.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Geladak** | **Namespace** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Geladak**:

- `nahkoda liat geladak` - Melihat semua area yang ada di atas kapal.
- `nahkoda bikin geladak [nama]` - Membuka area/lantai baru di kapal.
- `nahkoda liat kru di geladak [nama]` - Mencari kru yang sedang bekerja di area tertentu.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Menghapus Geladak yang masih memiliki Kru aktif (safety switch).
- Mengatur Resource Quota per Geladak.
- Berpindah namespace default secara permanen.
