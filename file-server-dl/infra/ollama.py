import requests
import json

class Config:
    enabled: bool
    model: str
    host: str
    port: int

class OllamaClient: 
    def __init__(self, config: Config):
        self.config = config

    def request_ollama_generate(self, body: str, model: str = None, image: list[str] = None, format: dict = None) -> dict:
        '''
        请求ollama生成文本，返回生成的文本 
        '''
        req_body = {
            'prompt': body,
            'stream': False
        }
        if model:
            req_body['model'] = model
        else:
            req_body['model'] = self.config.model
        if image:
            req_body['image'] = image
        if format:
            req_body['format'] = format
        resp = requests.post(f'http://{self.config.host}:{self.config.port}/api/generate', json=req_body)
        return json.loads(resp.json()['response'])