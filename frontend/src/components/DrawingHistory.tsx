import React from 'react';
import DrawingPredictionCard from './DrawingPredictionCard';

interface Drawing {
  id: string;
  imageData: string;
  prediction?: string;
  confidence?: number;
  timestamp: Date;
}

interface DrawingHistoryProps {
  drawings: Drawing[];
  onClearHistory: () => void;
  onDeleteDrawing?: (id: string) => void;
  onDrawingClick?: (drawing: Drawing) => void;
}

const DrawingHistory: React.FC<DrawingHistoryProps> = ({ 
  drawings, 
  onClearHistory, 
  onDeleteDrawing,
  onDrawingClick 
}) => {
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
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold text-gray-900">Drawing History</h2>
          <p className="text-gray-600 mt-1">{drawings.length} drawing{drawings.length !== 1 ? 's' : ''}</p>
        </div>
        <button
          onClick={onClearHistory}
          className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors duration-200 flex items-center space-x-2"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
          <span>Clear All</span>
        </button>
      </div>

      {/* Grid of cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
        {drawings.map((drawing) => (
          <DrawingPredictionCard
            key={drawing.id}
            drawing={drawing}
            onDelete={onDeleteDrawing}
            onClick={onDrawingClick}
            showActions={!!onDeleteDrawing}
          />
        ))}
      </div>
    </div>
  );
};

export default DrawingHistory;