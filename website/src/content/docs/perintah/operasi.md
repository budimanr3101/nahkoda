---
title: Operasi & Metrik
description: Cara membuat resource dan memantau performa.
---

Beberapa perintah operasi dan metrik yang didukung oleh Nahkoda.

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda bikin geladak [nama]` | `kubectl create namespace [nama]` | Buat namespace baru |
| `nahkoda bikin kru [nama]` | `kubectl run [nama] --image=nginx...` | Buat pod baru (nginx) |
| `nahkoda pantau kru` | `kubectl top pod` | Lihat penggunaan resource pod |
| `nahkoda pantau mesin` | `kubectl top node` | Lihat penggunaan resource node |
| `nahkoda hapus kru [nama]` | `kubectl delete pod [nama]` | Hapus pod |
