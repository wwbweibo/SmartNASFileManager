class FileUnderstandingResult:
    def __init__(self, label: str, group: str, description: str) -> None:
        self.label = label
        self.group = group
        self.description = description
        self.ext = None

    def set_ext(self, ext: any) -> None:
        self.ext = ext

    def __str__(self) -> str:
        return f"Label: {self.label}\nGroup: {self.group}\nDescription: {self.description}\nExtension: {self.ext}"
    
    def __repr__(self) -> str:
        return f"Label: {self.label}\nGroup: {self.group}\nDescription: {self.description}\nExtension: {self.ext}"
    
    def to_dict(self) -> dict:
        return {
            'label': self.label,
            'group': self.group,
            'description': self.description,
            'extension': self.ext.to_dict() if self.ext is not None else None
        }