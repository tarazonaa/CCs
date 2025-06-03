import React from 'react';

interface Drawing {
  id: string;
  imageData: string;
  prediction?: string;
  confidence?: number;
  timestamp: Date;
}

interface DrawingPredictionCardProps {
  drawing: Drawing;
  onDelete?: (id: string) => void;
  onClick?: (drawing: Drawing) => void;
  showActions?: boolean;
}

const DrawingPredictionCard: React.FC<DrawingPredictionCardProps> = ({ 
  drawing, 
  onDelete, 
  onClick,
  showActions = false 
}) => {
  return (
    <div
      className="bg-white rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 transform hover:scale-105 p-4 cursor-pointer"
      onClick={() => onClick?.(drawing)}
    >
      
      <div className="aspect-square mb-4 relative group">
        <img 
          src={drawing.imageData} 
          alt="Drawing"
          className="w-full h-full object-contain border border-gray-200 rounded-lg"
        />
        
        
        {showActions && onDelete && (
          <button
            onClick={(e) => {
              e.stopPropagation();
              onDelete(drawing.id);
            }}
            className="absolute top-2 right-2 bg-red-500 hover:bg-red-600 text-white rounded-full w-6 h-6 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200"
          >
            ×
          </button>
        )}
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
          <div className="flex items-center space-x-2">
            <span className="font-medium">
              {drawing.confidence 
                ? `${(drawing.confidence * 100).toFixed(1)}%` 
                : 'N/A'
              }
            </span>
        
            {drawing.confidence && (
              <div className="w-16 h-2 bg-gray-200 rounded-full overflow-hidden">
                <div 
                  className={`h-full transition-all duration-300 ${
                    drawing.confidence > 0.8 ? 'bg-green-500' :
                    drawing.confidence > 0.6 ? 'bg-yellow-500' : 'bg-red-500'
                  }`}
                  style={{ width: `${drawing.confidence * 100}%` }}
                />
              </div>
            )}
          </div>
        </div>
        
        <div className="pt-2 border-t border-gray-100">
          <span className="text-xs text-gray-400">
            {drawing.timestamp.toLocaleDateString()} {drawing.timestamp.toLocaleTimeString()}
          </span>
        </div>
      </div>
    </div>
  );
};

export default DrawingPredictionCard;