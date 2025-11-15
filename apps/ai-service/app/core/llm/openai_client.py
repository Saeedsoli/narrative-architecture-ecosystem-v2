# apps/ai-service/app/core/llm/openai_client.py

import logging
from openai import OpenAI, APIError, RateLimitError
from tenacity import retry, stop_after_attempt, wait_exponential
from typing import List

from ..config import settings

logger = logging.getLogger(__name__)

class OpenAIClient:
    def __init__(self):
        self.client = OpenAI(api_key=settings.OPENAI_API_KEY)

    @retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, min=1, max=10))
    def generate_completion(self, system_prompt: str, user_prompt: str, max_tokens: int = 400) -> str:
        try:
            response = self.client.chat.completions.create(
                model=settings.OPENAI_CHAT_MODEL,
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_prompt}
                ],
                max_tokens=max_tokens,
                temperature=0.7,
            )
            return response.choices[0].message.content or ""
        except RateLimitError as e:
            logger.warning(f"OpenAI rate limit hit. Retrying... Error: {e}")
            raise
        except APIError as e:
            logger.error(f"OpenAI API error: {e}")
            raise
        except Exception as e:
            logger.exception("An unexpected error occurred while generating completion.")
            raise

    def create_embedding(self, text: str) -> List[float]:
        text = text.replace("\n", " ")
        response = self.client.embeddings.create(input=[text], model=settings.OPENAI_EMBEDDING_MODEL)
        return response.data[0].embedding

llm_client = OpenAIClient()