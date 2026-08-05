#!/usr/bin/env python3

import typer
from functools import cached_property
from pathlib import Path
from urllib.request import urlretrieve
from pydantic import BaseModel, ConfigDict, TypeAdapter
from enum import IntEnum
from os.path import join
import requests

app = typer.Typer()


class StickerFormat(IntEnum):
    PNG = 1
    APNG = 2
    LOTTIE = 3
    GIF = 4

    @property
    def extension(self) -> str:
        return {
            StickerFormat.PNG: "png",
            StickerFormat.APNG: "png",
            StickerFormat.LOTTIE: "json",
            StickerFormat.GIF: "gif",
        }[self]


class Sticker(BaseModel):
    model_config = ConfigDict(extra="ignore")

    type: int
    tags: str
    name: str
    id: str
    guild_id: str
    format_type: StickerFormat
    description: str
    available: bool
    asset: str

    @cached_property
    def url(self) -> str:
        return f"https://cdn.discordapp.com/stickers/{self.id}.{self.format_type.extension}?lossless=true"

    @cached_property
    def filename(self) -> str:
        return f"{self.name}.{self.format_type.extension}"


StickerList = TypeAdapter(list[Sticker])


@app.command()
def main(
    input_file: Path = typer.Argument(exists=True, dir_okay=False),
    output_dir: Path = typer.Option(Path("./stickers/"), "-o", "--output"),
):
    with open(input_file) as f:
        for sticker in StickerList.validate_json(f.read()):
            response = requests.get(sticker.url)
            with open(join(output_dir, sticker.filename), "wb") as w:
                w.write(response.content)


if __name__ == "__main__":
    app()
