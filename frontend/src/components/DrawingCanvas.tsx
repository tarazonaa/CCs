import type React from "react";
import { useRef, useEffect, useState } from "react";
import { useSnackbar } from "notistack";
import { useTranslation } from "react-i18next";
import axios from "axios";
import { motion } from "framer-motion";
import {
  Trash,
  FloppyDisk,
  Info,
  DotsNineIcon as Grid,
} from "@phosphor-icons/react";

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
    ctx.strokeStyle = "rgba(0, 0, 0, 0.1)";
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
      const ctx = canvas.getContext("2d");
      if (ctx) {
        ctx.fillStyle = "white";
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        ctx.strokeStyle = "#000000";
        ctx.lineWidth = 8;
        ctx.lineCap = "round";
        ctx.lineJoin = "round";
        setContext(ctx);
        setCanvasCleared(true);
        if (showGrid) {
          drawGrid(ctx);
        }
      }
    }
  }, [showGrid]);

  const getCoordinates = (
    e:
      | React.MouseEvent<HTMLCanvasElement>
      | React.TouchEvent<HTMLCanvasElement>
      | TouchEvent
      | MouseEvent,
  ): { x: number; y: number } => {
    if (!canvasRef.current) return { x: 0, y: 0 };

    const rect = canvasRef.current.getBoundingClientRect();
    const scaleX = canvasRef.current.width / rect.width;
    const scaleY = canvasRef.current.height / rect.height;

    let clientX: number;
    let clientY: number;

    if ("touches" in e) {
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

  const startDrawing = (
    e:
      | React.MouseEvent<HTMLCanvasElement>
      | React.TouchEvent<HTMLCanvasElement>,
  ) => {
    if (!context) return;
    setIsDrawing(true);
    setCanvasCleared(false);

    const { x, y } = getCoordinates(e);

    context.beginPath();
    context.moveTo(x, y);
  };

  const draw = (
    e:
      | React.MouseEvent<HTMLCanvasElement>
      | React.TouchEvent<HTMLCanvasElement>,
  ) => {
    if (!isDrawing || !context) return;

    if ("touches" in e) {
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
    context.fillStyle = "white";
    context.fillRect(0, 0, canvasRef.current.width, canvasRef.current.height);
    if (showGrid) {
      drawGrid(context);
    }
    setCanvasCleared(true);
  };

  const saveDrawing = () => {
    if (!canvasRef.current || canvasCleared) return;

    const canvasImg = canvasRef.current.toDataURL();
    onDrawingComplete(canvasImg);

    const originalCanvas = canvasRef.current;
    const resizedCanvas = document.createElement("canvas");
    resizedCanvas.width = 112;
    resizedCanvas.height = 112;

    const ctx = resizedCanvas.getContext("2d");
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
      formData.append("image", blob, "processed.jpg");

      try {
        const response = await axios.post(
          "https://10.49.12.47:8443/api/v1/inference",
          formData,
        );
        setCurrBase64Img(response.data.segmentation_base64); // Assuming the response contains the image URL
      } catch (error: unknown) {
        if (error instanceof Error) {
          console.error("Upload error:", error.message);
        } else {
          console.error("Upload error: ", error?.response?.data);
        }
      } finally {
        const uploadFormData = new FormData();
        uploadFormData.append("original_image", blob, "original.jpg");
        // Convert the base64 image to a Blob
        const byteString = atob(currBase64Img);
        const ab = new ArrayBuffer(byteString.length);
        const ia = new Uint8Array(ab);
        for (let i = 0; i < byteString.length; i++) {
          ia[i] = byteString.charCodeAt(i);
        }
        const base64Blob = new Blob([ab], { type: "image/png" });
        uploadFormData.append("inference_image", base64Blob, "inference.jpg");
        return axios.post(
            `${import.meta.env.VITE_API_URL}/api/v1/images`,
            uploadFormData,
            {
              headers: {
                "Content-Type": "multipart/form-data",
                Authorization: `Bearer ${localStorage.getItem("access_token")}`,
              },
            }
          ).then(() => {
          // Tell parent to update the history
          onDrawingComplete(currBase64Img);
          enqueueSnackbar(t("drawing_saved"), { variant: "success" });
        })
        .catch(error => {
          console.warn("Image upload failed:", error.message);
          enqueueSnackbar(t("error_saving_drawing"), { variant: "error" });
        });
      }
    }, "image/jpeg");
  };

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const handleTouchMove = (e: TouchEvent) => {
      if (isDrawing) {
        e.preventDefault();
      }
    };

    canvas.addEventListener("touchmove", handleTouchMove, { passive: false });

    return () => {
      canvas.removeEventListener("touchmove", handleTouchMove);
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
            <h2 className="text-xl font-semibold text-text-primary">
              {t("draw_title")}
            </h2>
          </div>
          <div
            className="relative bg-white rounded-apple shadow-apple-md overflow-hidden mx-auto"
            style={{ width: "280px", height: "280px" }}
          >
            <canvas
              ref={canvasRef}
              width={280}
              height={280}
              style={{ width: "280px", height: "280px" }}
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
                <p className="text-text-tertiary text-sm italic">
                  {t("draw_here")}
                </p>
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
                  context.fillStyle = "white";
                  context.fillRect(
                    0,
                    0,
                    canvasRef.current!.width,
                    canvasRef.current!.height,
                  );
                  if (!showGrid) {
                    drawGrid(context);
                  }
                }
              }}
              className="btn flex items-center space-x-2 bg-surface-secondary hover:bg-surface-secondary-hover text-text-primary"
            >
              <Grid weight="bold" size={16} />
              <span>{t(showGrid ? "hide_grid" : "show_grid")}</span>
            </motion.button>
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={clearCanvas}
              className="btn flex items-center space-x-2 bg-surface-secondary hover:bg-surface-secondary-hover text-text-primary"
            >
              <Trash weight="bold" size={16} />
              <span>{t("clear")}</span>
            </motion.button>
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={saveDrawing}
              disabled={canvasCleared}
              className={`btn flex items-center space-x-2 ${
                canvasCleared
                  ? "bg-surface-disabled text-text-disabled cursor-not-allowed"
                  : "bg-green-500 text-white hover:bg-green-600"
              }`}
            >
              <FloppyDisk weight="bold" size={16} />
              <span>{t("save")}</span>
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
            <h3 className="text-lg font-semibold text-text-primary">
              {t("drawing_tips_title")}
            </h3>
          </div>

          <ul className="space-y-3 text-text-secondary">
            <li className="flex items-start space-x-2">
              <span className="w-5 h-5 rounded-full bg-primary/10 text-primary text-xs flex items-center justify-center mt-0.5">
                1
              </span>
              <span>{t("tip_1")}</span>
            </li>
            <li className="flex items-start space-x-2">
              <span className="w-5 h-5 rounded-full bg-primary/10 text-primary text-xs flex items-center justify-center mt-0.5">
                2
              </span>
              <span>{t("tip_2")}</span>
            </li>
            <li className="flex items-start space-x-2">
              <span className="w-5 h-5 rounded-full bg-primary/10 text-primary text-xs flex items-center justify-center mt-0.5">
                3
              </span>
              <span>{t("tip_3")}</span>
            </li>
          </ul>

          <div className="mt-6 p-4 rounded-lg bg-primary/5 border border-primary/10">
            <p className="text-text-secondary text-sm">
              {t("post_save_message")}
            </p>
          </div>
        </div>
      </motion.div>
      {currBase64Img && (
        <div className="mt-6">
          <h3 className="text-lg font-semibold mb-2">{t("returned_image")}</h3>
          <img
            src={`data:image/png;base64,${currBase64Img}`}
            alt="Processed Drawing"
            className="border rounded-lg shadow-md"
          />
        </div>
      )}
    </div>
  );
};

export default DrawingCanvas;

