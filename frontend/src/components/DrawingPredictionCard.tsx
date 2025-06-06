import React from 'react';
import { motion } from 'framer-motion';
import { Trash, ArrowsOut } from '@phosphor-icons/react';

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
  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  };

  return (
    <motion.div
      whileHover={{ y: -5, scale: 1.02 }}
      transition={{ duration: 0.2 }}
      className="card-glass overflow-visible"
    >
      <div 
        className="relative aspect-square cursor-pointer group"
        onClick={() => onClick?.(drawing)}
      >
        <img 
          src={drawing.imageData} 
          alt="Drawing"
          className="w-full h-full object-contain bg-white rounded-t-apple"
        />
        
        <div className="absolute inset-0 bg-black/0 group-hover:bg-black/5 transition-colors duration-200 rounded-t-apple flex items-center justify-center opacity-0 group-hover:opacity-100">
          <motion.div
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            className="w-10 h-10 rounded-full bg-white/80 backdrop-blur-sm flex items-center justify-center shadow-apple"
          >
            <ArrowsOut weight="bold" size={18} className="text-primary" />
          </motion.div>
        </div>
        
        {showActions && onDelete && (
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={(e) => {
              e.stopPropagation();
              onDelete(drawing.id);
            }}
            className="absolute -top-2 -right-2 w-8 h-8 bg-error text-white rounded-full flex items-center justify-center shadow-apple opacity-0 group-hover:opacity-100 transition-opacity duration-200 z-10"
          >
            <Trash weight="bold" size={14} />
          </motion.button>
        )}
      </div>
        
      <div className="p-4 space-y-3">
        <div className="flex justify-between items-center">
          <span className="text-text-secondary text-sm">Prediction</span>
          <span className="font-semibold text-xl text-primary">
            {drawing.prediction || '?'}
          </span>
        </div>
        
        <div className="space-y-1">
          <div className="flex justify-between items-center text-sm">
            <span className="text-text-secondary">Confidence</span>
            <span className="font-medium">
              {drawing.confidence 
                ? `${(drawing.confidence * 100).toFixed(0)}%` 
                : 'N/A'
              }
            </span>
          </div>
          
          <div className="w-full h-1.5 bg-background rounded-full overflow-hidden">
            <div 
              className={`h-full transition-all duration-300 ${
                !drawing.confidence ? 'bg-gray-300' :
                drawing.confidence > 0.8 ? 'bg-success' :
                drawing.confidence > 0.6 ? 'bg-warning' : 'bg-error'
              }`}
              style={{ width: `${drawing.confidence ? drawing.confidence * 100 : 0}%` }}
            />
          </div>
        </div>
        
        <div className="pt-2 border-t border-border/30">
          <span className="text-xs text-text-secondary">
            {formatDate(drawing.timestamp)}
          </span>
        </div>
      </div>
    </motion.div>
  );
};

export default DrawingPredictionCard;