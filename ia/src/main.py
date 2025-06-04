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
    if image.size != (112, 112):
        return JSONResponse(
            status_code=400,
            content={"error": "Image must be 112x112 pixels."},
        )
    image_tensor = (
        torch.from_numpy(np.array(image)).unsqueeze(0).unsqueeze(0).float() / 255.0
    )
    return image_tensor


def image_to_base64(img_array: np.ndarray) -> str:
    if img_array.shape[-1] == 4:
        # ARGB to RGB conversion, drop alpha and reorder channels
        img_array = img_array[..., [1, 2, 3]]
    img = Image.fromarray(img_array.astype(np.uint8))
    buffered = io.BytesIO()
    img = img.convert("RGB")  # Ensure image is in RGB format
    # save image to buffer
    img.save(buffered, "PNG")
    return base64.b64encode(buffered.getvalue()).decode()


@app.get("/")
def read_root():
    return {
        "message": "Welcome to the Image Segmentation API @CC. To make an inference, send a POST request with a 112x112 image file for inference."
    }


@app.post("/")
async def segment_image(image: UploadFile = File(...)):

    try:
        x = preprocess_image(image)
        if isinstance(x, JSONResponse):
            return x  # Return error response if preprocessing failed
        segmentation_rgb = model.inference(x)  # shape: (H, W, 3)
        img_b64 = image_to_base64(segmentation_rgb)
        return JSONResponse(content={"segmentation_base64": img_b64})

    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})
