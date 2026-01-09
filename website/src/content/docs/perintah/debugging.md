---
title: Debugging
description: Cara mendiagnosa masalah pada pods dan nodes.
---

Gunakan perintah debugging untuk melihat jurnal (logs) atau masuk langsung ke dalam kapal.

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda baca jurnal [pod]` | `kubectl logs [pod]` | Baca logs dari pod |
| `nahkoda masuk [pod]` | `kubectl exec -it [pod] -- /bin/sh` | Masuk ke container (shell) |
| `nahkoda cek kru [pod]` | `kubectl describe pod [pod]` | Lihat detail/status pod |
| `nahkoda cek mesin [node]` | `kubectl describe node [node]` | Lihat detail node |
