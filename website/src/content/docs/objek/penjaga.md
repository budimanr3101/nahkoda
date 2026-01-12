---
title: Penjaga (DaemonSet)
description: Mengenal Penjaga, kru yang wajib ada di setiap jengkal kapal.
---

**Penjaga** adalah metafora untuk **DaemonSet** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Penjaga? Karena ia memiliki tugas khusus untuk memastikan setiap Mesin (Node) memiliki satu petugas yang berjaga. Penjaga biasanya digunakan untuk memantau keamanan, mencatat aktifitas mesin, atau fungsi pendukung lainnya yang harus ada di setiap unit mesin tanpa terkecuali.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Penjaga** | **DaemonSet** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Penjaga**:

- `nahkoda liat penjaga` - Melihat daftar petugas yang berjaga di setiap mesin.
- `nahkoda cek penjaga [nama]` - Melihat detail distribusi penjaga di seluruh armada mesin.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Melakukan update strategi penjaga.
- Menghapus penjaga secara parsial di mesin tertentu.
