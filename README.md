# Rinha de Backend – Solução em Go + PostgreSQL

## 1 · Problema

Implementar **4 endpoints** (`POST /pessoas`, `GET /pessoas/{id}`, `GET /pessoas?t=`, `GET /contagem-pessoas`) e suportar alta concorrência em um ambiente limitado a **1 ½ vCPU** e **3 GB RAM** distribuídos entre duas instâncias da API (Go), um Nginx (load-balancer) e um PostgreSQL.

---

## 2 · Resultado (Gatling)

| Métrica            | Valor          |
| ------------------ | -------------- |
| Requisições totais | **115 006**    |
| Duração do teste   | **3 min 25 s** |
| Erros (KO)         | **0**          |
| p50                | 1 ms           |
| **p75**            | **2 ms**       |
| **p95**            | **100 ms**     |
| p99                | 220 ms         |

---

![image](https://github.com/user-attachments/assets/862de4c2-8798-424b-95b0-7c70df90c689)


## 3 · Otimizações que fizeram diferença

| #   | Técnica                                          | Impacto                                                     |
| --- | ------------------------------------------------ | ----------------------------------------------------------- |
| 1   | **Pool de conexões `pgxpool`** (`MaxConns ≈ 15`) | Evitou _thrashing_, estabilizando a latência                |
| 2   | **Nginx `worker_connections 4096`**              | Suportou ~4 k conexões simultâneas por processo; p75 ≤ 2 ms |
| 3   | **Índice GiST `pg_trgm`** na coluna `searchable` | Busca substring até **50×** mais rápida                     |

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_people_search_trgm
  ON people
  USING gist (searchable gist_trgm_ops(siglen = 64));
```
