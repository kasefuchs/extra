from pathlib import Path

import typer
from rich import print
from dataclasses import asdict
from PIL import Image

from android_recovery_anim.model import Info

app = typer.Typer()


@app.command()
def info(
    input_file: Path = typer.Argument(exists=True, dir_okay=False),
    as_json: bool = typer.Option(False, "-j", "--json"),
):
    with Image.open(input_file) as img:
        info = Info.from_image(img)

    if as_json:
        print(asdict(info))
    else:
        typer.echo(f"Valid: {info.is_valid}")
        typer.echo(f"Resolution: {info.width}x{info.height} px")
        typer.echo(f"Frame size: {info.width}x{info.frame_height} px")
        typer.echo(f"Frames: {info.frames}")
        typer.echo(f"FPS: {info.fps}")
