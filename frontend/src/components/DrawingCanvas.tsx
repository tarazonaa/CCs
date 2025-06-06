import React, { useRef, useEffect, useState } from 'react';
import { useSnackbar } from 'notistack';
import { useTranslation } from 'react-i18next';
import axios from 'axios';

interface DrawingCanvasProps {
  onDrawingComplete: (imageData: string) => void;
}

const DrawingCanvas: React.FC<DrawingCanvasProps> = ({ onDrawingComplete }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [currBase64Img, setCurrBase64Img] = useState<string>("");
  const [context, setContext] = useState<CanvasRenderingContext2D | null>(null);
  const { enqueueSnackbar } = useSnackbar();
  const { t } = useTranslation();

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
    
    let canvasImg = canvasRef.current.toDataURL();
    onDrawingComplete(canvasImg);

    const originalCanvas = canvasRef.current;

    // Create a resized 112x112 canvas
    const resizedCanvas = document.createElement('canvas');
    resizedCanvas.width = 112;
    resizedCanvas.height = 112;

    const ctx = resizedCanvas.getContext('2d');
    if (!ctx) {
      console.error("Failed to get context for resized canvas");
      return;
    }

    // Draw the original image resized
    ctx.drawImage(originalCanvas, 0, 0, 112, 112);

    // Get the pixel data
    const imageData = ctx.getImageData(0, 0, 112, 112);
    const data = imageData.data;

    // Convert to grayscale and invert colors
    for (let i = 0; i < data.length; i += 4) {
      const r = data[i];
      const g = data[i + 1];
      const b = data[i + 2];

      // Grayscale: average the RGB values
      const gray = (r + g + b) / 3;

      // Invert grayscale value (255 - gray)
      const inverted = 255 - gray;

      data[i] = inverted;     // Red
      data[i + 1] = inverted; // Green
      data[i + 2] = inverted; // Blue
      // Alpha stays the same (data[i + 3])
    }

    // Put the processed image back to the canvas
    ctx.putImageData(imageData, 0, 0);

    // Export as JPEG and upload
    resizedCanvas.toBlob(async (blob) => {
      if (!blob) {
        console.error("Failed to convert processed canvas to Blob.");
        return;
      }

      const formData = new FormData();
      formData.append('image', blob, 'processed.jpg');

      try {
        const response = await axios.post('https://10.49.12.47:8443/api/v1/inference', formData);
        setCurrBase64Img(response.data.segmentation_base64); // Assuming the response contains the image URL
      } catch (error: any) {
        console.error('Upload error:', error?.response?.data || error.message);
      } finally {
        const uploadFormData = new FormData();
        uploadFormData.append('original_image', blob, 'original.jpg');
        // Turn the base64 string into a Blob for upload
        const base64Response = await fetch(`data:image/png;base64,${currBase64Img}`);
        const base64Blob = await base64Response.blob();
        uploadFormData.append('inference_image', base64Blob, 'inference.jpg');
        await axios.post(`${import.meta.env.VITE_AUTH_ENDPOINT}/api/v1/images`, uploadFormData, {
          headers: {
            'Content-Type': 'multipart/form-data',  
            'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
          },
        });
        enqueueSnackbar(t('drawing_saved'), { variant: 'success' });
      }
    }, 'image/jpeg');
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
    {
      // Returned image is a base64 string
      currBase64Img && (
        <div className="mt-6">
          <h3 className="text-lg font-semibold mb-2">Returned Image:</h3>
          <img src={`data:image/png;base64,${currBase64Img}`} alt="Processed Drawing" className="border rounded-lg shadow-md" />
        </div>
      )
    }
    </div>
  );
};

export default DrawingCanvas;
