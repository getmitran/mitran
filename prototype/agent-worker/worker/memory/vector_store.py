"""Vector memory store using ChromaDB for semantic similarity search."""

import time
from typing import Optional

import chromadb


class VectorMemory:
    def __init__(self, persist_dir: str):
        self.client = chromadb.PersistentClient(path=persist_dir)
        self.facts = self.client.get_or_create_collection("facts")
        self.episodes = self.client.get_or_create_collection("episodes")
        self.corrections = self.client.get_or_create_collection("corrections")

    def add_fact(self, key: str, value: str) -> None:
        self.facts.upsert(
            ids=[key],
            documents=[value],
            metadatas=[{"key": key, "timestamp": time.time()}],
        )

    def add_episode(self, content: str, category: str = "general") -> None:
        doc_id = f"ep_{int(time.time() * 1000)}"
        self.episodes.add(
            ids=[doc_id],
            documents=[content],
            metadatas=[{"category": category, "timestamp": time.time()}],
        )

    def add_correction(self, original: str, corrected: str, reason: str) -> None:
        doc_id = f"cor_{int(time.time() * 1000)}"
        self.corrections.add(
            ids=[doc_id],
            documents=[f"{original} -> {corrected}"],
            metadatas=[{
                "original": original,
                "corrected": corrected,
                "reason": reason,
                "timestamp": time.time(),
            }],
        )

    def search(self, query: str, n_results: int = 5, collection: Optional[str] = None) -> list[dict]:
        targets = (
            [self._get_collection(collection)]
            if collection
            else [self.facts, self.episodes, self.corrections]
        )
        results = []
        for col in targets:
            if col.count() == 0:
                continue
            res = col.query(query_texts=[query], n_results=min(n_results, col.count()))
            for i, doc in enumerate(res["documents"][0]):
                results.append({
                    "document": doc,
                    "metadata": res["metadatas"][0][i],
                    "distance": res["distances"][0][i] if res.get("distances") else None,
                    "collection": col.name,
                })
        results.sort(key=lambda x: x.get("distance") or float("inf"))
        return results[:n_results]

    def list_facts(self) -> list[dict]:
        if self.facts.count() == 0:
            return []
        res = self.facts.get()
        return [
            {"key": m["key"], "value": d}
            for d, m in zip(res["documents"], res["metadatas"])
        ]

    def list_episodes(self) -> list[dict]:
        if self.episodes.count() == 0:
            return []
        res = self.episodes.get()
        return [
            {"content": d, "category": m["category"], "timestamp": m["timestamp"]}
            for d, m in zip(res["documents"], res["metadatas"])
        ]

    def delete_fact(self, key: str) -> bool:
        try:
            self.facts.delete(ids=[key])
            return True
        except Exception:
            return False

    def stats(self) -> dict:
        return {
            "facts": self.facts.count(),
            "episodes": self.episodes.count(),
            "corrections": self.corrections.count(),
        }

    def _get_collection(self, name: str):
        return {"facts": self.facts, "episodes": self.episodes, "corrections": self.corrections}[name]
