---
title: Monitoring (Pods & Nodes)
description: Cara memantau kesehatan kru dan mesin.
---

Gunakan perintah monitoring untuk melihat status kru (pods) dan mesin (nodes) Anda.

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda liat kru` | `kubectl get pods -A` | List semua pod di semua namespace |
| `nahkoda liat kru di geladak [ns]` | `kubectl get pods -n [ns]` | List pod di namespace tertentu |
| `nahkoda liat kru rusak` | `kubectl get pods ... \| grep -v Running` | Cari pod yang error/crash |
| `nahkoda liat mesin` | `kubectl get nodes` | List worker nodes |
| `nahkoda liat berita` | `kubectl get events --sort-by=...` | Lihat event cluster terbaru |

:::note
Secara default, `liat kru` tanpa kondisi akan menampilkan semua pod di semua namespace.
:::
