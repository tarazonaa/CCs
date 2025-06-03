import torch
import torch.nn.functional as F
import torchvision.models as models
import lightning as L
import numpy as np
import matplotlib.pyplot as plt

from matplotlib import colors
from matplotlib.backends.backend_agg import FigureCanvasAgg as FigureCanvas

class UNetConvRelu(torch.nn.Sequential):
    def __init__(self, in_channels=32, out_channels=32):
        super().__init__(
                torch.nn.Conv2d(in_channels,out_channels,
                                kernel_size=3, padding=1, bias=False),
                torch.nn.BatchNorm2d(out_channels),
                torch.nn.ReLU(inplace=True)
        )

class UNetDecoderBlock(torch.nn.Module):
    def __init__(self, in_channels, skip_channels, out_channels):
        super().__init__()
        self.conv1 = UNetConvRelu(in_channels + skip_channels, out_channels)
        self.conv2 = UNetConvRelu(out_channels, out_channels)

    def forward(self, x, skip):
        x = F.interpolate(x, size=skip.shape[2:], mode='bilinear', align_corners=False)
        x = torch.cat([x, skip], dim=1)
        x = self.conv1(x)
        x = self.conv2(x)
        return x


class Unet(torch.nn.Module):
    def __init__(self, in_channels=1, classes=11):
        super().__init__()
        base_model = models.resnet34(weights=models.ResNet34_Weights.IMAGENET1K_V1)

        if in_channels != 3:
            old_conv = base_model.conv1
            base_model.conv1 = torch.nn.Conv2d(
                    in_channels, 64,
                    kernel_size=7,
                    stride=2,
                    padding=3,
                    bias=False
            )

            with torch.no_grad():
                if in_channels == 1:
                    base_model.conv1.weight[:] = old_conv.weight.sum(dim=1, keepdim=True)
                else:
                    base_model.conv1.weight[:, :in_channels] = old_conv.weight[:, :in_channels]

        # Encoder
        self.encoder0 = torch.nn.Sequential(base_model.conv1,
                                            base_model.bn1,
                                            base_model.relu)
        self.encoder1 = torch.nn.Sequential(base_model.maxpool,
                                      base_model.layer1)
        self.encoder2 = base_model.layer2
        self.encoder3 = base_model.layer3
        self.encoder4 = base_model.layer4

        # Decoder
        self.decoder4 = UNetDecoderBlock(512, 256, 256)
        self.decoder3 = UNetDecoderBlock(256, 128, 128)
        self.decoder2 = UNetDecoderBlock(128, 64, 64)
        self.decoder1 = UNetDecoderBlock(64, 64, 32)
        self.decoder0 = UNetConvRelu(32, 32)

        # Final segmentation head
        self.segmentation_head = torch.nn.Conv2d(32, classes, kernel_size=1)

    def forward(self, x):
        # Encoder
        e0 = self.encoder0(x)
        e1 = self.encoder1(e0)
        e2 = self.encoder2(e1)
        e3 = self.encoder3(e2)
        e4 = self.encoder4(e3)

        # Decoder
        d4 = self.decoder4(e4, e3)
        d3 = self.decoder3(d4, e2)
        d2 = self.decoder2(d3, e1)
        d1 = self.decoder1(d2, e0)
        d0 = self.decoder0(F.interpolate(d1, scale_factor=2,
                                         mode='bilinear', align_corners=False))

        return self.segmentation_head(d0)

class UNetCE(L.LightningModule):
    def __init__(self):
        super().__init__()
        self.unet = Unet(in_channels=1, classes=11)
        self.loss_fn = torch.nn.CrossEntropyLoss()

    def forward(self, x):
        return self.unet(x)

    def inference(self, x, device=None):
        if device is None:
            device = torch.device("cuda" if torch.cuda.is_available() else "cpu")       
        
        self.to(device)
        self.eval()

        with torch.no_grad():
            logits = self.forward(x.to(device)) 
            pred = torch.nn.functional.softmax(logits[0], dim=0).argmax(dim=0)
            pred = pred.cpu().numpy()

        digit_colors = ['black', 'red', 'orange', 'yellow', 'lime',
                        'lightgreen', 'cyan', 'blue', 'indigo', 'purple',
                        'violet']
        
        cmap = colors.ListedColormap(digit_colors)
        bounds = np.linspace(-0.5, 10.5, 12)
        norm = colors.BoundaryNorm(bounds, cmap.N)

        fig, ax = plt.subplots(figsize=(4, 4), dpi=100)
        im = ax.imshow(pred, cmap=cmap, norm=norm)
        ax.axis('off')
        cbar = fig.colorbar(im, ax=ax, ticks=np.arange(11),
                            fraction=0.046, pad=0.04)
        cbar.set_ticklabels(
          ['bg'] + [str(i) for i in range(1, 10)] + ['0']
        )

        canvas = FigureCanvas(fig)
        canvas.draw()

        width, height = fig.canvas.get_width_height()
        img = np.frombuffer(canvas.tostring_argb(), dtype='uint8')
        img = img.reshape((height, width, 4))

        plt.close(fig)

        return img

ckpt_file = "semantic-segmentation-unet-cross-entropy-epoch=21.ckpt"
model = UNetCE.load_from_checkpoint(ckpt_file)
