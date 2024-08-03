class ImageLabel:
    def __init__(self, label: str, confidence: float):
        self.label = label
        self.confidence = confidence

    def __str__(self):
        return f"{self.label}: {self.confidence:.2f}%"
    
    def __repr__(self):
        return f"{self.label}: {self.confidence:.2f}%"
    
    def to_dict(self):
        return {
            'label': self.label,
            'confidence': f"{self.confidence:.2f}%"
        }

class ImageUnderstandingResult:
    def __init__(self, labels: list[ImageLabel], caption: str):
        self.labels = labels
        self.caption = caption

    def __str__(self):
        return f"Labels: {self.labels}\nCaption: {self.caption}"
    
    def __repr__(self):
        return f"Labels: {self.labels}\nCaption: {self.caption}"
    
    def to_dict(self):
        return {
            'labels': [label.to_dict() for label in self.labels],
            'caption': self.caption
        }