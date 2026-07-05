from __future__ import annotations

from typing import Self
from functools import cached_property
from dataclasses import dataclass

from PIL import Image


@dataclass
class Info:
    fps: int
    frames: int

    width: int
    height: int

    @cached_property
    def is_valid(self) -> bool:
        return self.height % self.frames == 0

    @cached_property
    def frame_height(self) -> int:
        return self.height // self.frames

    @classmethod
    def from_image(cls, img: Image.Image) -> Self:
        width, height = img.size

        fps = img.info.get("FPS", 20)
        frames = img.info.get("Frames")
        if not frames:
            raise ValueError("Could not find frame count in metadata")

        return cls(
            fps=int(fps),
            frames=int(frames),
            width=width,
            height=height,
        )
