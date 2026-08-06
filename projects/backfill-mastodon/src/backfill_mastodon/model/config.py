from functools import cached_property
from pathlib import Path
from pydantic import BaseModel, Field
from pydantic_settings import BaseSettings, SettingsConfigDict, PydanticBaseSettingsSource
from mastodon import Mastodon


class TargetConfig(BaseModel):
    account: str
    limit: int = 40


class InstanceConfig(BaseModel):
    api_url: str = ""
    access_token: str = ""

    @cached_property
    def client(self) -> Mastodon:
        return Mastodon(
            api_base_url=self.api_url,
            access_token=self.access_token or None,
        )


class RemoteConfig(InstanceConfig):
    targets: list[TargetConfig] = Field(default_factory=list)


class Config(BaseSettings):
    model_config = SettingsConfigDict(
        extra="ignore",
        env_prefix="BACKFILL_MASTODON_",
        env_nested_delimiter="__",
    )

    local: InstanceConfig = Field(default_factory=InstanceConfig)
    remotes: list[RemoteConfig] = Field(default_factory=list)
    state_path: Path = Path("state.json")

    @classmethod
    def settings_customise_sources(
        cls,
        settings_cls: type[BaseSettings],
        init_settings: PydanticBaseSettingsSource,
        env_settings: PydanticBaseSettingsSource,
        dotenv_settings: PydanticBaseSettingsSource,
        file_secret_settings: PydanticBaseSettingsSource,
    ) -> tuple[PydanticBaseSettingsSource, ...]:
        return env_settings, init_settings, file_secret_settings
