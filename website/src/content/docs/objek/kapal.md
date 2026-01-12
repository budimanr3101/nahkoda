---
title: Kapal (Context/Cluster)
description: Berpindah antar armada besar (Cluster) melalui perintah Kapal.
---

**Kapal** adalah metafora untuk **Context/Cluster** dalam Kubernetes.

### ⚓ Filosofi
Kenapa disebut Kapal? Karena Anda mungkin memiliki kendali atas lebih dari satu kapal atau cluster. Terkadang Anda sedang di atas Kapal Produksi, terkadang di Kapal Pengembangan. Perintah ini membantu Anda mengetahui di mana posisi Anda sekarang dan berpindah ke kapal lain dengan cepat.

### 📋 Ekuivalen Kubernetes
| Nahkoda | Kubernetes |
| :--- | :--- |
| **Kapal** | **Context / Cluster** |

### 🚀 Perintah Populer
Berikut adalah beberapa perintah yang bisa Anda jalankan terhadap **Kapal**:

- `nahkoda liat kapal` - Melihat semua kapal yang bisa Anda nahkodai.
- `nahkoda pindah kapal [nama]` - Berpindah kendali ke kapal lain secara instan.

### ⚠️ Yang Belum Bisa Dilakukan
Saat ini, Nahkoda belum mendukung:
- Menambahkan Kapal (cluster) baru secara manual.
- Mengedit konfigurasi kubeconfig.
- Menghapus Kapal dari daftar navigasi.
