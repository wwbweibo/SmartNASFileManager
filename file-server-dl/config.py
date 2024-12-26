from infra.ollama import Config as OllamaConfig
import yaml

class Config:
    def __init__(self):
        self.ollama = OllamaConfig()
        self.nas_root_path = ""
    def from_yaml_file(self, file_path: str):
        with open(file_path, 'r') as f:
            config = yaml.safe_load(f)
            self.ollama = OllamaConfig()
            self.ollama.enabled = config['ollama']['enabled']
            self.ollama.model = config['ollama']['model']
            self.ollama.host = config['ollama']['host']
            self.ollama.port = config['ollama']['port']
            self.nas_root_path = config['nas_root_path']
        return self