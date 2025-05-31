
import React from 'react';
import type { Drawing } from '@/types/types';

interface DrawingHistoryProps {
  drawings: Drawing[];
  onClearHistory: () => void;
}

const DrawingHistory: React.FC<DrawingHistoryProps> = ({ drawings, onClearHistory }) => {
  if (drawings.length === 0) {
    return (
      <div className="text-center py-16">
        <div className="bg-white rounded-2xl shadow-lg p-12 max-w-md mx-auto">
          <div className="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">No drawings yet</h2>
          <p className="text-gray-600">Go to the Draw tab to create your first drawing!</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold text-gray-900">Drawing History</h2>
        <button
          onClick={onClearHistory}
          className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors duration-200"
        >
          Clear All
        </button>
      </div>
      

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
        {drawings.map((drawing) => (
          <div
            key={drawing.id}
            className="bg-white rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 transform hover:scale-105 p-4"
          >
            <div className="aspect-square mb-4">
              <img 
                src={drawing.imageData} 
                alt="Drawing"
                className="w-full h-full object-contain border border-gray-200 rounded-lg"
              />
            </div>
            
            <div className="space-y-2 text-sm">
              <div className="flex justify-between items-center">
                <span className="text-gray-500">Prediction:</span>
                <span className="font-semibold text-lg text-blue-600">
                  {drawing.prediction || '?'}
                </span>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="text-gray-500">Confidence:</span>
                <span className="font-medium">
                  {drawing.confidence 
                    ? `${(drawing.confidence * 100).toFixed(1)}%` 
                    : 'N/A'
                  }
                </span>
              </div>
              
              <div className="pt-2 border-t border-gray-100">
                <span className="text-xs text-gray-400">
                  {drawing.timestamp.toLocaleDateString()} {drawing.timestamp.toLocaleTimeString()}
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default DrawingHistory;