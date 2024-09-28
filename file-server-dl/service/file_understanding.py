from  magika import Magika
from pathlib import Path
from models.file import FileUnderstandingResult
from service.image_understanding import ImageUnderstanding
import logging

class FileUnderstanding:
    def __init__(self, config: any):
        self.magika = Magika()
        self.image_understanding = ImageUnderstanding()
        self.config = config

    def understand(self, path: str) -> FileUnderstandingResult:
        result = self.magika.identify_path(Path(self.config['nas_root_path']  + path))
        file_understanding = FileUnderstandingResult(result.output.ct_label, result.output.group, result.output.description)
        logging.info("File Understanding: %s", file_understanding)
        if file_understanding.group == 'image':
            file_understanding.set_ext(self.image_understanding.understand( self.config['nas_root_path'] + path))
        logging.info("File Understanding Result: %s", file_understanding)
        return file_understanding