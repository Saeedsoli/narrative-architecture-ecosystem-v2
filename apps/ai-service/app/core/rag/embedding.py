# apps/ai-service/app/core/rag/embedding.py

from openai import OpenAI
from ..config import settings

class EmbeddingService:
    def __init__(self, api_key: str, model: str = "text-embedding-3-small"):
        self.client = OpenAI(api_key=api_key)
        self.model = model

    def create_embedding(self, text: str) -> list[float]:
        """
        Creates a vector embedding for a given text.
        """
        text = text.replace("\n", " ")
        response = self.client.embeddings.create(input=[text], model=self.model)
        return response.data[0].embedding

# Singleton instance
embedding_service = EmbeddingService(api_key=settings.OPENAI_API_KEY, model=settings.OPENAI_EMBEDDING_MODEL)