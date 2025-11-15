# apps/ai-service/app/main.py

from fastapi import FastAPI, Depends, HTTPException, Security
from fastapi.security import APIKeyHeader
import logging

from .api import analyze
from .core.config import settings

# --- Logger Setup ---
logging.basicConfig(level=settings.LOG_LEVEL.upper(), format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

# --- API Key Security ---
API_KEY_NAME = "X-API-KEY"
api_key_header = APIKeyHeader(name=API_KEY_NAME, auto_error=False)

async def get_api_key(api_key_header: str = Security(api_key_header)):
    if api_key_header == settings.API_KEY:
        return api_key_header
    else:
        raise HTTPException(status_code=403, detail="Could not validate credentials")

# --- App Initialization ---
app = FastAPI(
    title="Narrative Architecture AI Service",
    version="1.0.0",
    dependencies=[Depends(get_api_key)]
)

# --- Routers ---
app.include_router(analyze.router, prefix="/v1")

@app.get("/health", tags=["Health"])
async def health_check():
    return {"status": "ok"}

@app.on_event("startup")
async def startup_event():
    logger.info("AI Service is starting up...")
    # Можно добавить проверку подключения к Pinecone/OpenAI