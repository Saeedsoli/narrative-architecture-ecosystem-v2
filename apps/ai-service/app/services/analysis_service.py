# apps/ai-service/app/services/analysis_service.py

from ..core.rag.retriever import retriever
from ..core.llm.openai_client import llm_client
from ..core.llm.prompts import Prompts

class AnalysisService:
    def analyze_text(self, user_text: str) -> dict:
        # 1. Retrieve relevant context from the book (via RAG)
        book_context = retriever.retrieve_relevant_context(user_text)

        # 2. (Optional) Find relevant platform links (e.g., from Elasticsearch)
        # For now, we'll use a placeholder.
        platform_links = [
            {"title": "مقاله مرتبط: ساختار پیرنگ", "url": "/articles/plot-structure"}
        ]

        # 3. Create the prompt for the LLM
        user_prompt = Prompts.get_analysis_user_prompt(
            user_text=user_text,
            book_context=book_context,
            platform_links=platform_links
        )

        # 4. Call the LLM
        response_content = llm_client.generate(
            system_prompt=Prompts.ANALYSIS_SYSTEM_PROMPT,
            user_prompt=user_prompt
        )

        # 5. Format the final response
        # Here you might parse the LLM's response into a structured JSON
        return {
            "analysis": response_content,
            "retrieved_context": book_context, # For logging/debugging
        }

# Singleton instance
analysis_service = AnalysisService()