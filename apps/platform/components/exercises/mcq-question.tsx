// apps/platform/components/exercises/mcq-question.tsx

'use client';

import { useState } from 'react';

interface Option {
  id: string;
  text: string;
}

interface McqQuestionProps {
  question: string;
  options: Option[];
  correctAnswerId: string;
  explanation: string;
  onAnswer: (isCorrect: boolean) => void;
}

export function McqQuestion({ question, options, correctAnswerId, explanation, onAnswer }: McqQuestionProps) {
  const [selectedOption, setSelectedOption] = useState<string | null>(null);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const handleSubmit = () => {
    if (!selectedOption) return;
    setIsSubmitted(true);
    const isCorrect = selectedOption === correctAnswerId;
    onAnswer(isCorrect);
  };

  const getOptionClasses = (optionId: string) => {
    if (!isSubmitted) {
      return selectedOption === optionId ? 'bg-blue-200' : 'hover:bg-gray-100';
    }
    if (optionId === correctAnswerId) {
      return 'bg-green-200 border-green-400';
    }
    if (optionId === selectedOption) {
      return 'bg-red-200 border-red-400';
    }
    return 'bg-gray-100';
  };

  return (
    <div className="p-6 border rounded-lg">
      <p className="font-semibold mb-4">{question}</p>
      <div className="space-y-2">
        {options.map((option) => (
          <button
            key={option.id}
            onClick={() => !isSubmitted && setSelectedOption(option.id)}
            disabled={isSubmitted}
            className={`w-full text-right p-3 border rounded-md transition-colors ${getOptionClasses(option.id)}`}
          >
            {option.text}
          </button>
        ))}
      </div>
      
      {!isSubmitted ? (
        <button
          onClick={handleSubmit}
          disabled={!selectedOption}
          className="mt-6 px-6 py-2 bg-blue-600 text-white rounded-md disabled:opacity-50"
        >
          ثبت پاسخ
        </button>
      ) : (
        <div className="mt-6 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
          <h4 className="font-bold">توضیحات:</h4>
          <p>{explanation}</p>
        </div>
      )}
    </div>
  );
}