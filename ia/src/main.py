import io
import base64
import torch
import numpy as np

from fastapi import FastAPI
from fastapi import FastAPI, File, UploadFile
from fastapi.responses import JSONResponse
from PIL import Image

from model import model

app = FastAPI()

def preprocess_image(file: UploadFile) -> torch.Tensor:
    image = Image.open(file.file).convert("L")  # convert to grayscale
    image_tensor = torch.from_numpy(np.array(image)).unsqueeze(0).unsqueeze(0).float() / 255.0
    return image_tensor

def image_to_base64(img_array: np.ndarray) -> str:
    img = Image.fromarray(img_array.astype(np.uint8))
    buffered = io.BytesIO()
    img.save(buffered, format="PNG")
    return base64.b64encode(buffered.getvalue()).decode()

@app.post("/")
async def segment_image(image: UploadFile = File(...)):
    try:
        x = preprocess_image(image)
        segmentation_rgb = model.inference(x)  # shape: (H, W, 3)
        img_b64 = image_to_base64(segmentation_rgb)
        return JSONResponse(content={"segmentation_base64": img_b64})
    
    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})
