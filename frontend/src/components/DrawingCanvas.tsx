import React, { useRef, useEffect, useState } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { 
  Trash, 
  FloppyDisk,
  Info,
  DotsNineIcon as Grid
} from '@phosphor-icons/react';

interface DrawingCanvasProps {
  onDrawingComplete: (imageData: string) => void;
}

const DrawingCanvas: React.FC<DrawingCanvasProps> = ({ onDrawingComplete }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [context, setContext] = useState<CanvasRenderingContext2D | null>(null);
  const [canvasCleared, setCanvasCleared] = useState(true);
  const [lastPos, setLastPos] = useState({ x: 0, y: 0 });
  const [showGrid, setShowGrid] = useState(false);

  const drawGrid = (ctx: CanvasRenderingContext2D) => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const w = canvas.width;
    const h = canvas.height;
    const cellWidth = w / 3;
    const cellHeight = h / 3;

    ctx.save();
    ctx.strokeStyle = 'rgba(0, 0, 0, 0.1)';
    ctx.lineWidth = 1;

    // Draw vertical lines
    for (let i = 1; i < 3; i++) {
      ctx.beginPath();
      ctx.moveTo(cellWidth * i, 0);
      ctx.lineTo(cellWidth * i, h);
      ctx.stroke();
    }

    // Draw horizontal lines
    for (let i = 1; i < 3; i++) {
      ctx.beginPath();
      ctx.moveTo(0, cellHeight * i);
      ctx.lineTo(w, cellHeight * i);
      ctx.stroke();
    }
    ctx.restore();
  };

  useEffect(() => {
    const canvas = canvasRef.current;
    if (canvas) {
      const ctx = canvas.getContext('2d');
      if (ctx) {
        ctx.fillStyle = 'white';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        ctx.strokeStyle = '#000000';
        ctx.lineWidth = 8;
        ctx.lineCap = 'round';
        ctx.lineJoin = 'round';
        setContext(ctx);
        setCanvasCleared(true);
        if (showGrid) {
          drawGrid(ctx);
        }
      }
    }
  }, [showGrid]);

  const getCoordinates = (
    e: React.MouseEvent<HTMLCanvasElement> | React.TouchEvent<HTMLCanvasElement> | TouchEvent | MouseEvent
  ): { x: number, y: number } => {
    if (!canvasRef.current) return { x: 0, y: 0 };
    
    const rect = canvasRef.current.getBoundingClientRect();
    const scaleX = canvasRef.current.width / rect.width;
    const scaleY = canvasRef.current.height / rect.height;
    
    let clientX: number, clientY: number;
    
    if ('touches' in e) {
      if (e.touches.length === 0) {
        return lastPos;
      }
      clientX = e.touches[0].clientX;
      clientY = e.touches[0].clientY;
    } else {
      clientX = e.clientX;
      clientY = e.clientY;
    }
    
    const x = (clientX - rect.left) * scaleX;
    const y = (clientY - rect.top) * scaleY;
    
    setLastPos({ x, y });
    
    return { x, y };
  };

  const startDrawing = (e: React.MouseEvent<HTMLCanvasElement> | React.TouchEvent<HTMLCanvasElement>) => {
    if (!context) return;
    setIsDrawing(true);
    setCanvasCleared(false);
    
    const { x, y } = getCoordinates(e);
    
    context.beginPath();
    context.moveTo(x, y);
  };

  const draw = (e: React.MouseEvent<HTMLCanvasElement> | React.TouchEvent<HTMLCanvasElement>) => {
    if (!isDrawing || !context) return;
    
    if ('touches' in e) {
      e.preventDefault();
    }
    
    const { x, y } = getCoordinates(e);
    
    context.lineTo(x, y);
    context.stroke();
  };

  const stopDrawing = () => {
    if (!context) return;
    setIsDrawing(false);
    context.closePath();
  };

  const clearCanvas = () => {
    if (!context || !canvasRef.current) return;
    context.fillStyle = 'white';
    context.fillRect(0, 0, canvasRef.current.width, canvasRef.current.height);
    if (showGrid) {
      drawGrid(context);
    }
    setCanvasCleared(true);
  };

  const saveDrawing = () => {
    if (!canvasRef.current || canvasCleared) return;

    let canvasImg = canvasRef.current.toDataURL();
    onDrawingComplete(canvasImg);

    const originalCanvas = canvasRef.current;
    const resizedCanvas = document.createElement('canvas');
    resizedCanvas.width = 112;
    resizedCanvas.height = 112;

    const ctx = resizedCanvas.getContext('2d');
    if (!ctx) {
      console.error("Failed to get context for resized canvas");
      return;
    }

    ctx.drawImage(originalCanvas, 0, 0, 112, 112);
    const imageData = ctx.getImageData(0, 0, 112, 112);
    const data = imageData.data;

    for (let i = 0; i < data.length; i += 4) {
      const r = data[i];
      const g = data[i + 1];
      const b = data[i + 2];
      const gray = (r + g + b) / 3;
      const inverted = 255 - gray;
      data[i] = inverted;
      data[i + 1] = inverted;
      data[i + 2] = inverted;
    }

    ctx.putImageData(imageData, 0, 0);

    resizedCanvas.toBlob(async (blob) => {
      if (!blob) {
        console.error("Failed to convert processed canvas to Blob.");
        return;
      }

      const formData = new FormData();
      formData.append('image', blob, 'processed.jpg');

      try {
        const response = await axios.post('https://10.49.12.47:8443/api/v1/inference', formData);
        console.log('Upload success:', response.data);
      } catch (error: any) {
        console.error('Upload error:', error?.response?.data || error.message);
      }
    }, 'image/jpeg');
  };

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    
    const handleTouchMove = (e: TouchEvent) => {
      if (isDrawing) {
        e.preventDefault();
      }
    };
    
    canvas.addEventListener('touchmove', handleTouchMove, { passive: false });
    
    return () => {
      canvas.removeEventListener('touchmove', handleTouchMove);
    };
  }, [isDrawing]);

  return (
    <div className="flex flex-col items-center gap-8">
      {/* Canvas Container */}
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <div className="card-glass p-6">
          <div className="mb-6 text-center">
            <h2 className="text-xl font-semibold text-text-primary">Draw a Digit (0-9)</h2>
          </div>
          <div className="relative bg-white rounded-apple shadow-apple-md overflow-hidden mx-auto" style={{ width: '280px', height: '280px' }}>
            <canvas
              ref={canvasRef}
              width={280}
              height={280}
              style={{ width: '280px', height: '280px' }}
              className="cursor-crosshair drawing-canvas touch-none"
              onMouseDown={startDrawing}
              onMouseMove={draw}
              onMouseUp={stopDrawing}
              onMouseLeave={stopDrawing}
              onTouchStart={startDrawing}
              onTouchMove={draw}
              onTouchEnd={stopDrawing}
            />
            {canvasCleared && (
              <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                <p className="text-text-tertiary text-sm italic">Draw here</p>
              </div>
            )}
          </div>

          <div className="flex justify-center mt-6 space-x-4">
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={() => {
                setShowGrid(!showGrid);
                if (context) {
                  context.fillStyle = 'white';
                  context.fillRect(0, 0, canvasRef.current!.width, canvasRef.current!.height);
                  if (!showGrid) {
                    drawGrid(context);
                  }
                }
              }}
              className="btn flex items-center space-x-2 bg-surface-secondary hover:bg-surface-secondary-hover text-text-primary"
            >
              <Grid weight="bold" size={16} />
              <span>{showGrid ? 'Hide Grid' : 'Show Grid'}</span>
            </motion.button>
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={clearCanvas}
              className="btn flex items-center space-x-2 bg-surface-secondary hover:bg-surface-secondary-hover text-text-primary"
            >
              <Trash weight="bold" size={16} />
              <span>Clear</span>
            </motion.button>
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={saveDrawing}
              disabled={canvasCleared}
              className={`btn flex items-center space-x-2 ${
                canvasCleared 
                  ? 'bg-surface-disabled text-text-disabled cursor-not-allowed' 
                  : 'bg-green-500 text-white hover:bg-green-600'
              }`}
            >
              <FloppyDisk weight="bold" size={16} />
              <span>Save</span>
            </motion.button>
          </div>
        </div>
      </motion.div>

      {/* Tips Card */}
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
        className="w-full max-w-md"
      >
        <div className="card-glass p-6">
          <div className="flex items-center space-x-3 mb-4">
            <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
              <Info weight="bold" size={20} className="text-primary" />
            </div>
            <h3 className="text-lg font-semibold text-text-primary">Drawing Tips</h3>
          </div>
          
          <ul className="space-y-3 text-text-secondary">
            <li className="flex items-start space-x-2">
              <span className="inline-block w-5 h-5 rounded-full bg-primary/10 text-primary text-xs flex items-center justify-center mt-0.5">1</span>
              <span>Draw a clear digit (0-9) in the center of the canvas</span>
            </li>
            <li className="flex items-start space-x-2">
              <span className="inline-block w-5 h-5 rounded-full bg-primary/10 text-primary text-xs flex items-center justify-center mt-0.5">2</span>
              <span>Make it large enough to fill most of the drawing area</span>
            </li>
            <li className="flex items-start space-x-2">
              <span className="inline-block w-5 h-5 rounded-full bg-primary/10 text-primary text-xs flex items-center justify-center mt-0.5">3</span>
              <span>Click "Save" to add your drawing to your history</span>
            </li>
          </ul>

          <div className="mt-6 p-4 rounded-lg bg-primary/5 border border-primary/10">
            <p className="text-text-secondary text-sm">
              After saving, your drawing will be processed by our AI to recognize the digit you've drawn.
            </p>
          </div>
        </div>
      </motion.div>
    </div>
  );
};

export default DrawingCanvas;