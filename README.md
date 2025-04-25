# Embedding Service

This microservice provides vector embeddings for user queries as part of the RAG pipeline powering **MattBot** on [matthewpsimons.com](https://matthewpsimons.com).

It wraps a local build of `llama-embedding` and exposes a lightweight HTTP API for downstream services (such as `llm-orchestrator`) to request embeddings.

---

## Purpose

The Embedding Service transforms raw text queries into dense vector representations that can be used to retrieve semantically relevant chunks from a vector database (Qdrant).

---

## System Role in RAG Pipeline

```text
user query → llm-orchestrator → embedding service → vector → llm-orchestrator → Qdrant
```

---

## Features

- Uses `llama-embedding` for local, self-hosted vector generation
- CPU-based embedding (fast for small queries)
- Simple JSON HTTP API
- Suitable for proof-of-concept or lightweight workloads

---

## API

### POST `/embed`

**Body:**
```json
{
  "text": "What are Matt's thoughts on platform enablement?"
}

```

**Response:**
```json
{
  "vector": [0.031, -0.204, ...]
}
```

---

## Development Notes

- The Docker image **should be built on the host** to take advantage of host-specific optimizations from the `llama-embedding` binary.
- This service currently runs in **CPU mode**. It is performant enough for short queries and prototyping.
- **Longer queries work**, but response times increase noticeably. For production or batch-scale workloads, GPU acceleration is recommended.

---

## Future Improvements

- Add GPU support for high-throughput embedding
- Support batch embedding
- Improve observability (timing logs, error reporting)

---

## License

© 2025 Matthew Simons. All rights reserved.

---

## Contact

Questions or suggestions? Reach out via [matthewpsimons.com](https://matthewpsimons.com).

