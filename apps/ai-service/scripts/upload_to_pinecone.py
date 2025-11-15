# apps/ai-service/scripts/seed_pinecone.py

import json
from app.core.rag.embedding import embedding_service
from app.core.rag.vector_store import vector_store
from ulid import ULID

def seed():
    # فرض می‌کنیم محتوای کتاب در یک فایل JSON است
    with open('book_content_chunks.json', 'r', encoding='utf-8') as f:
        chunks = json.load(f)

    vectors_to_upsert = []
    for chunk in chunks:
        # Create embedding for each chunk of text
        embedding = embedding_service.create_embedding(chunk['text'])
        
        # Prepare vector with metadata
        vector = (
            str(ULID()),  # Unique ID for the vector
            embedding,
            {
                "title": chunk['title'],
                "chapter": chunk['chapter'],
                "keywords": chunk['keywords'],
                # IMPORTANT: DO NOT STORE THE 'text' IN METADATA
            }
        )
        vectors_to_upsert.append(vector)
    
    # Upsert in batches
    batch_size = 100
    for i in range(0, len(vectors_to_upsert), batch_size):
        batch = vectors_to_upsert[i:i+batch_size]
        vector_store.upsert_vectors(batch)
        print(f"Upserted batch {i//batch_size + 1}")

if __name__ == "__main__":
    seed()