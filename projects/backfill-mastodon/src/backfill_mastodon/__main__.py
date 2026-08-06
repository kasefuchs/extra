from pathlib import Path
import typer
from rich.progress import Progress, TextColumn, BarColumn, TaskProgressColumn, MofNCompleteColumn, TimeRemainingColumn
from backfill_mastodon.model import Context

app = typer.Typer()


@app.command()
def main(config_path: Path = Path("config.yaml")):
    ctx = Context.load(config_path)
    with Progress(
        TextColumn("[bold blue]{task.description}"),
        BarColumn(),
        TaskProgressColumn(),
        MofNCompleteColumn(),
        TimeRemainingColumn(),
    ) as progress:
        try:
            local_client = ctx.config.local.client

            for remote in ctx.config.remotes:
                remote_client = remote.client

                for target in remote.targets:
                    key = f"{remote.api_url}:{target.account}"
                    processed = ctx.state.progress[key]
                    task = progress.add_task(f"Backfilling {target.account} from {remote.api_url}")

                    try:
                        account = remote_client.account_lookup(target.account)
                        progress.update(task, total=account.statuses_count, advance=processed)

                        while True:
                            last_id = ctx.state.max_id.get(key)
                            page = remote_client.account_statuses(account, limit=target.limit, max_id=last_id)
                            if not page:
                                break

                            for status in page:
                                progress.update(task, advance=1)
                                if status.url:
                                    local_client.search(status.url, resolve=True)

                                ctx.state.progress[key] += 1
                                ctx.state.max_id[key] = status.id

                    except Exception as e:
                        typer.secho(f"Failed to process {target.account}: {e}", fg=typer.colors.RED)

                    finally:
                        progress.update(task, completed=True)

        finally:
            ctx.save()


if __name__ == "__main__":
    app()
