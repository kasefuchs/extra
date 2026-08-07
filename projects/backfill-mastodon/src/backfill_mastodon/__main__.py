import asyncio
from pathlib import Path
from async_typer import AsyncTyper
import typer
from rich.progress import Progress, TextColumn, BarColumn, TaskProgressColumn, MofNCompleteColumn, TimeRemainingColumn
from backfill_mastodon.model import Context

app = AsyncTyper()


@app.command()
async def main(config_path: Path = Path("config.yaml")):
    ctx = Context.load(config_path)
    semaphore = asyncio.Semaphore(ctx.config.concurrency)
    local_client = ctx.config.local.client

    with Progress(
        TextColumn("[bold blue]{task.description}"),
        BarColumn(),
        TaskProgressColumn(),
        MofNCompleteColumn(),
        TimeRemainingColumn(),
    ) as progress:
        try:
            for remote in ctx.config.remotes:
                remote_client = remote.client

                for target in remote.targets:
                    key = f"{remote.api_url}:{target.account}"
                    processed = ctx.state.progress[key]
                    task = progress.add_task(f"Backfilling {target.account} from {remote.api_url}")

                    async def process_status(status):
                        if status.url:
                            async with semaphore:
                                await asyncio.to_thread(local_client.search, status.url, resolve=True)

                        progress.update(task, advance=1)
                        ctx.state.progress[key] += 1

                    try:
                        account = remote_client.account_lookup(target.account)
                        progress.update(task, total=account.statuses_count, advance=processed)

                        while True:
                            last_id = ctx.state.max_id.get(key)
                            page = remote_client.account_statuses(
                                account,
                                limit=target.limit,
                                max_id=last_id,
                                exclude_replies=target.exclude_replies,
                                exclude_reblogs=target.exclude_reblogs,
                            )
                            if not page:
                                break

                            await asyncio.gather(*(process_status(s) for s in page))

                            ctx.state.max_id[key] = page[-1].id

                    except Exception as e:
                        typer.secho(f"Failed to process {target.account}: {e}", fg=typer.colors.RED)

                    finally:
                        progress.update(task, completed=True)

        finally:
            ctx.save()


if __name__ == "__main__":
    app()
