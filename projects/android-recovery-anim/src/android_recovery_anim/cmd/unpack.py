from pathlib import Path

import typer
from PIL import Image

from android_recovery_anim.util import extract_frames

app = typer.Typer()


@app.command()
def unpack(
    input_file: Path = typer.Argument(exists=True, dir_okay=False),
    output_dir: Path = typer.Option(Path("./frames/"), "-o", "--output"),
):
    with Image.open(input_file) as img:
        frames = extract_frames(img)

    output_dir.mkdir(parents=True, exist_ok=True)
    for idx, frame in enumerate(frames):
        frame.save(output_dir / f"{idx:03d}.png")

    typer.secho(f"Successfully extracted {len(frames)} frames", fg=typer.colors.GREEN, bold=True)
