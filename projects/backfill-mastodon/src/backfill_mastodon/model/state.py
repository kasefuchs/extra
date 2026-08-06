from collections import defaultdict
from pydantic import BaseModel, ConfigDict, Field
from mastodon.types_base import IdType


class State(BaseModel):
    model_config = ConfigDict(arbitrary_types_allowed=True)

    progress: dict[str, int] = Field(default_factory=lambda: defaultdict(int))
    max_id: dict[str, IdType | None] = Field(default_factory=dict)
