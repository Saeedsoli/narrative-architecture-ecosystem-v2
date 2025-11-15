# apps/ai-service/app/core/nlp/preprocess.py

from hazm import Normalizer, sent_tokenize

class TextPreprocessor:
    def __init__(self):
        self.normalizer = Normalizer()

    def normalize(self, text: str) -> str:
        """
        Applies basic normalization like fixing unicode characters,
        standardizing punctuation, and removing unnecessary whitespaces.
        """
        return self.normalizer.normalize(text)

    def split_into_sentences(self, text: str) -> list[str]:
        """
        Splits a block of text into individual sentences.
        """
        return sent_tokenize(text)

# Singleton instance
preprocessor = TextPreprocessor()