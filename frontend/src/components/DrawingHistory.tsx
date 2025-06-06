import React from 'react';
import { motion } from 'framer-motion';
import { Trash, ImageSquare } from '@phosphor-icons/react';
import { ImageMetadata } from '@/pages/Dashboard';
import axios from 'axios';


interface DrawingHistoryProps {
  imagesData: ImageMetadata[];
  onClearHistory: () => void;
  onDeleteDrawing?: (id: string) => void;
  // onDrawingClick?: (drawing: Drawing) => void;
}

const DrawingHistory: React.FC<DrawingHistoryProps> = ({
  imagesData,
  onClearHistory,
  onDeleteDrawing,
  // onDrawingClick
}) => {
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1,
      transition: { 
        when: "beforeChildren", 
        staggerChildren: 0.1 
      }
    }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: { 
      y: 0, 
      opacity: 1,
      transition: { duration: 0.3 }
    }
  };

  if (imagesData.length === 0) {
    return (
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="flex flex-col items-center justify-center py-16 text-center"
      >
        <div className="card-glass p-10 max-w-md mx-auto">
          <div className="w-16 h-16 bg-background rounded-full flex items-center justify-center mx-auto mb-6">
            <ImageSquare weight="fill" size={28} className="text-text-secondary" />
          </div>
          <h2 className="text-2xl font-semibold text-text-primary mb-3">No drawings yet</h2>
          <p className="text-text-secondary mb-6">Create your first drawing to see it appear here</p>
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={() => window.location.hash = ''}
            className="btn btn-primary mx-auto"
          >
            Start Drawing
          </motion.button>
        </div>
      </motion.div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-semibold text-text-primary">Drawing History</h2>
          <p className="text-text-secondary mt-1">
            {imagesData.length} drawing{imagesData.length !== 1 ? 's' : ''}
          </p>
        </div>
        <motion.button
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          onClick={onClearHistory}
          className="btn flex items-center space-x-2 bg-error/10 hover:bg-error/20 text-error border-none"
        >
          <Trash weight="bold" size={16} />
          <span>Clear All</span>
        </motion.button>
      </div>

      {/* Grid of cards */}
      <motion.div 
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6"
      >
        {imagesData.map( async (imgData, i) =>  {
            const fetchImageBlob = async (imageId: string) => {
              try {
              const res = await axios.get(`${import.meta.env.VITE_AUTH_ENDPOINT}/api/v1/images/blob/${imageId}`, {
                headers: {
                'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
                }

              })
              return URL.createObjectURL(new Blob([res.data], { type: 'image/png' }));
            } catch (error) {
              console.error('Error fetching image blob:', error);
              return imageId; // Fallback to imageId if fetch fails
            }
            }

          let sentImageSrc = `${import.meta.env.VITE_AUTH_ENDPOINT}/api/v1/images/blob/${imgData.sent_image_id}`;
          let receivedImageSrc = `${import.meta.env.VITE_AUTH_ENDPOINT}/api/v1/images/blob/${imgData.received_image_id}`;
          sentImageSrc = await fetchImageBlob(imgData.sent_image_id) || sentImageSrc;
          receivedImageSrc = await fetchImageBlob(imgData.received_image_id) || receivedImageSrc;
          return <motion.div key={i} variants={itemVariants}>
            <div className="card-glass p-4">
              <div className="flex flex-col space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  {/* Original Drawing */}
                  <div className="space-y-2">
                    <p className="text-sm font-medium text-text-secondary">Original</p>
                    <div className="bg-white rounded-lg p-2 shadow-sm">
                      <img 
                        src={sentImageSrc}
                        alt="Original drawing" 
                        className="w-full aspect-square object-contain"
                      />
                    </div>
                  </div>
                  
                  {/* Segmentation (currently same as original) */}
                  <div className="space-y-2">
                    <p className="text-sm font-medium text-text-secondary">Segmentation</p>
                    <div className="bg-white rounded-lg p-2 shadow-sm">
                      <img 
                        src={receivedImageSrc}
                        alt="Segmentation" 
                        className="w-full aspect-square object-contain"
                      />
                    </div>
                  </div>
                </div>

                {/* Prediction Info */}
                <div className="flex justify-between items-center">
                
                  {onDeleteDrawing && (
                    <motion.button
                      whileHover={{ scale: 1.05 }}
                      whileTap={{ scale: 0.95 }}
                      // onClick={() => onDeleteDrawing(drawing.id)}
                      className="text-error hover:text-error-dark transition-colors"
                    >
                      <Trash weight="bold" size={16} />
                    </motion.button>
                  )}
                </div>
              </div>
            </div>
          </motion.div>}
        )}
      </motion.div>
    </div>
  );
};

export default DrawingHistory;