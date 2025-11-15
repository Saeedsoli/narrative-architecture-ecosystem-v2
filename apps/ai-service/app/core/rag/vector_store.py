# apps/ai-service/app/core/rag/vector_store.py

from pinecone import Pinecone, ServerlessSpec
from ..config import settings

class VectorStore:
    def __init__(self, api_key: str, environment: str, index_name: str):
        self.pc = Pinecone(api_key=api_key, environment=environment)
        self.index_name = index_name
        self.index = None
        self._initialize_index()

    def _initialize_index(self):
        """
        Connects to an existing index or creates a new one if it doesn't exist.
        """
        if self.index_name not in self.pc.list_indexes().names():
            # Dimension for 'text-embedding-3-small' is 1536
            self.pc.create_index(
                name=self.index_name,
                dimension=1536,
                metric="cosine",
                spec=ServerlessSpec(
                    cloud="aws",
                    region="us-west-2"
                )
            )
        self.index = self.pc.Index(self.index_name)

    def upsert_vectors(self, vectors: list[tuple[str, list[float], dict]]):
        """
        Upserts vectors into the Pinecone index.
        Each vector is a tuple of (id, vector, metadata).
        """
        if not self.index:
            raise ConnectionError("Pinecone index not initialized.")
        self.index.upsert(vectors=vectors, namespace="book-content")

    def query(self, vector: list[float], top_k: int = 3) -> list[dict]:
        """
        Queries the index to find the most similar vectors.
        """
        if not self.index:
            raise ConnectionError("Pinecone index not initialized.")
            
        results = self.index.query(
            namespace="book-content",
            vector=vector,
            top_k=top_k,
            include_metadata=True
        )
        return [match['metadata'] for match in results['matches']]

# Singleton instance
vector_store = VectorStore(
    api_key=settings.PINECONE_API_KEY,
    environment=settings.PINECONE_ENVIRONMENT,
    index_name=settings.PINECONE_INDEX_NAME
)