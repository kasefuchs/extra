#!/usr/bin/env python3

import typer
from mastodon import Mastodon

app = typer.Typer()


@app.command()
def main(
    api_url: str = typer.Option(
        prompt=typer.style("Enter Mastodon instance URL", fg=typer.colors.YELLOW),
    ),
    app_name: str | None = None,
    scopes: list[str] = ["read", "write"],
    client_id: str | None = None,
    client_secret: str | None = None,
    website: str | None = None,
    open_browser: bool = True,
):
    if not (client_id and client_secret):
        prompt_text = typer.style("Enter application name", fg=typer.colors.YELLOW)
        client_name = app_name or typer.prompt(prompt_text)

        client_id, client_secret = Mastodon.create_app(
            client_name=client_name,
            scopes=scopes,
            website=website,
            api_base_url=api_url,
        )

    client = Mastodon(
        client_id=client_id,
        client_secret=client_secret,
        api_base_url=api_url,
    )

    auth_url = client.auth_request_url(scopes=scopes)

    typer.secho(f"Visit URL: {auth_url}", fg=typer.colors.BLUE)

    if open_browser:
        typer.secho("Opening browser...", fg=typer.colors.CYAN)
        typer.launch(auth_url)

    code_prompt = typer.style("Enter authorization code", fg=typer.colors.YELLOW)
    code = typer.prompt(code_prompt)
    access_token = client.log_in(code=code.strip(), scopes=scopes)

    typer.secho(access_token, fg=typer.colors.GREEN, bold=True)


if __name__ == "__main__":
    app()
