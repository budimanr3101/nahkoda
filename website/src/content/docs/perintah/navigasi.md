---
title: Navigasi (Context)
description: Cara berpindah antar cluster dan context.
---

Gunakan perintah navigasi untuk melihat dan berpindah antar kapal (cluster/context) yang tersedia.

| Perintah Nahkoda | Ekuivalen Kubectl | Fungsi |
| :--- | :--- | :--- |
| `nahkoda liat kapal` | `kubectl config get-contexts` | List semua cluster/context |
| `nahkoda pindah kapal [nama]` | `kubectl config use-context [nama]` | Pindah cluster aktif |

:::tip
Analogi: **Kapal** dalam Nahkoda adalah representasi dari Kubernetes Context atau Cluster.
:::
