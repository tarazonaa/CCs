# Inference API

This directory contains a FastAPI service under `src/`, that serves a ML model
that solves a semantic segmentation task for handwritten number. The app serves
an inference POST endpoint (`/`) in `src/main.py`.

The model was done using Pytorch and Lightning using a UNet architecture. The 
actual neural network architecture 
can be seen in code in `src/model.py`. The trained model is loaded from a ![checkpoint file](https://drive.google.com/file/d/1DrBeBG18pWiuqKvqvap7tT5ktiN3CnFt)
and furhter information about it can be found in the ![slide deck](SemanticSegmentation.pdf).

## Deployment

The application was deployed to the local GPU Lab at our campus and a unit file to execute as a `systemd` service:
```
[Unit]
After=network.target

[Service]
User=$USER
WorkingDirectory=/path/to/this/directory
ExecStart=/path/to/this/directory/bin/uvicorn main:app --host 0.0.0.0 --port 8000
Restart=always

[Install]
WantedBy=multi-user.target
```

The unit file assumes that there is a venv that is used to start `uvicorn`.

Also, if you refer to the architecture diagram, you'll see that the GPU Lab is considered a VPC, this is due
to the fact that On-Prem, the Cloud and GPU Lab are in different VLANs and we had to configure a router to 
allow their communication.
