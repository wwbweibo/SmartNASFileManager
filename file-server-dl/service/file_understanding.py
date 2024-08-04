from  magika import Magika
from pathlib import Path
from models.file import FileUnderstandingResult
from service.image_understanding import ImageUnderstanding
import logging

class FileUnderstanding:
    def __init__(self):
        self.magika = Magika()
        self.image_understanding = ImageUnderstanding()

    def understand(self, path: str) -> FileUnderstandingResult:
        result = self.magika.identify_path(Path(path))
        file_understanding = FileUnderstandingResult(result.output.ct_label, result.output.group, result.output.description)
        if file_understanding.group == 'image':
            file_understanding.set_ext(self.image_understanding.understand(path))
        logging.info("File Understanding Result: %s", file_understanding)
        return file_understanding