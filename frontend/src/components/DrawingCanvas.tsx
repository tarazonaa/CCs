import React, { useRef, useEffect, useState } from 'react';

interface DrawingCanvasProps {
  onDrawingComplete: (imageData: string) => void;
}

const DrawingCanvas: React.FC<DrawingCanvasProps> = ({ onDrawingComplete }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [context, setContext] = useState<CanvasRenderingContext2D | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (canvas) {
      const ctx = canvas.getContext('2d');
      if (ctx) {
        ctx.fillStyle = 'white';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        ctx.strokeStyle = 'black';
        ctx.lineWidth = 8;
        ctx.lineCap = 'round';
        ctx.lineJoin = 'round';
        setContext(ctx);
      }
    }
  }, []);

  const startDrawing = (e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!context) return;
    
    setIsDrawing(true);
    const rect = canvasRef.current!.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
    context.beginPath();
    context.moveTo(x, y);
  };

  const draw = (e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isDrawing || !context) return;
    
    const rect = canvasRef.current!.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
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
  };

  const saveDrawing = () => {
    if (!canvasRef.current) return;
    
    const imageData = canvasRef.current.toDataURL();
    onDrawingComplete(imageData);
  };

  return (
    <div className="flex flex-col items-center space-y-8">
      {/* Canvas Container */}
      <div className="bg-white rounded-2xl shadow-lg p-6">
        <canvas
          ref={canvasRef}
          width={280}
          height={280}
          className="border-2 border-gray-300 rounded-lg cursor-crosshair hover:border-blue-400 transition-colors duration-200"
          onMouseDown={startDrawing}
          onMouseMove={draw}
          onMouseUp={stopDrawing}
          onMouseLeave={stopDrawing}
        />
      </div>
    
      <div className="flex space-x-4">
        <button
          onClick={clearCanvas}
          className="bg-red-600 hover:bg-red-700 text-white px-6 py-3 rounded-lg font-semibold transition-all duration-200 transform hover:scale-105 shadow-lg"
        >
          Clear Canvas
        </button>
        <button
          onClick={saveDrawing}
          className="bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-lg font-semibold transition-all duration-200 transform hover:scale-105 shadow-lg"
        >
          Save Drawing
        </button>
      </div>
      
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 max-w-md text-center">
        <p className="text-blue-800 text-sm">
          Draw a digit (0-9) in the canvas above, then click "Save Drawing" to add it to your history.
        </p>
      </div>
    </div>
  );
};

export default DrawingCanvas;