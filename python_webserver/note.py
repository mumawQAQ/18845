from typing import Dict


class Note:
    def __init__(self, id: int, title: str, content: str):
        self.id = id
        self.title = title
        self.content = content

    def to_dict(self) -> Dict:
        return {"id": self.id, "title": self.title, "content": self.content}
