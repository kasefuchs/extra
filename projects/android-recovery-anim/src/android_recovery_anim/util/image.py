from typing import List, Tuple

from PIL import Image


from android_recovery_anim.model import Info


def extract_frames(img: Image.Image) -> List[Image.Image]:
    info = Info.from_image(img)
    if not info.is_valid:
        raise ValueError("Source image is not valid")

    img_bytes = img.convert("RGBA").tobytes()
    row_stride = len(img_bytes) // info.height

    extracted = []
    for i in range(info.frames):
        data = bytearray()
        for y in range(info.frame_height):
            row = y * info.frames + i
            start = row * row_stride
            end = start + row_stride
            data.extend(img_bytes[start:end])

        frame = Image.frombytes("RGBA", (info.width, info.frame_height), data)
        extracted.append(frame)

    return extracted
