from  magika import Magika
from pathlib import Path

m = Magika()

def interfer_file_type(path: str) -> list[str]:
	result = m.identify_path(Path(path))
	return [result.output.ct_label, result.output.group, result.output.description]