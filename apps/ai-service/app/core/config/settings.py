# apps/ai-service/app/core/config/settings.py

from pydantic import BaseSettings, Field

class Settings(BaseSettings):
    # Server
    HOST: str = "0.0.0.0"
    PORT: int = 8000
    LOG_LEVEL: str = "info"

    # OpenAI
    OPENAI_API_KEY: str = Field(..., env="OPENAI_API_KEY")
    OPENAI_EMBEDDING_MODEL: str = "text-embedding-3-small"
    OPENAI_CHAT_MODEL: str = "gpt-4-turbo-preview"

    # Pinecone
    PINECONE_API_KEY: str = Field(..., env="PINECONE_API_KEY")
    PINECONE_ENVIRONMENT: str = "gcp-starter"
    PINECONE_INDEX_NAME: str = "narrative-arch-book"
    
    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"

settings = Settings()