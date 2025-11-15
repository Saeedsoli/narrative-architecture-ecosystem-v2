# apps/ai-service/app/api/analyze.py

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from ..services.analysis_service import analysis_service

router = APIRouter()

class AnalyzeRequest(BaseModel):
    text: str
    context: str | None = None # e.g., "This is for the 'Plot' exercise"

class AnalyzeResponse(BaseModel):
    analysis: str
    # We don't return the context to the user

@router.post("/analyze", response_model=AnalyzeResponse)
async def analyze_text_endpoint(request: AnalyzeRequest):
    if not request.text or len(request.text.strip()) < 50:
        raise HTTPException(status_code=400, detail="Text is too short for a meaningful analysis.")
    
    try:
        result = analysis_service.analyze_text(request.text)
        return AnalyzeResponse(analysis=result["analysis"])
    except Exception as e:
        # Log the full error
        print(f"Analysis failed: {e}")
        raise HTTPException(status_code=500, detail="Failed to perform analysis.")