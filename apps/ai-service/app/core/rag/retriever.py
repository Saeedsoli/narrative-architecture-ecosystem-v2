# apps/ai-service/app/core/rag/retriever.py

from .embedding import embedding_service
from .vector_store import vector_store

class Retriever:
    def __init__(self, embed_service, vec_store):
        self.embedding_service = embed_service
        self.vector_store = vec_store

    def retrieve_relevant_context(self, query: str, top_k: int = 3) -> str:
        """
        Takes a user query, creates an embedding, queries the vector store,
        and returns a formatted context string.
        """
        # 1. Create embedding for the user's query
        query_embedding = self.embedding_service.create_embedding(query)

        # 2. Query Pinecone for similar documents
        results = self.vector_store.query(vector=query_embedding, top_k=top_k)

        # 3. Format the results into a context string for the LLM
        # We only provide metadata, not the full text, to prevent leaks.
        context_str = "Relevant concepts from 'Narrative Architecture' book:\n"
        for i, result in enumerate(results):
            context_str += f"{i+1}. Concept: {result.get('title')}\n"
            context_str += f"   - Chapter: {result.get('chapter', 'N/A')}\n"
            context_str += f"   - Keywords: {', '.join(result.get('keywords', []))}\n\n"
        
        return context_str

# Singleton instance
retriever = Retriever(embed_service=embedding_service, vec_store=vector_store)