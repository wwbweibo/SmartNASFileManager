import logging
import torch
import base64
from io import BytesIO
import cn_clip.clip as clip
from cn_clip.clip import load_from_name
from lavis.models import load_model_and_preprocess
from models.image import ImageUnderstandingResult, ImageLabel
from PIL import Image
import importlib.util as importutil
import numpy as np
from transformers import AutoModel, AutoImageProcessor
from infra.milvus import conn as milvus_conn
from infra.ollama import OllamaClient, Config as OllamaConfig
from config import Config

class ImageUnderstanding:
    def __init__(self, config: Config):
        # check if torch_directml is available
        if importutil.find_spec("torch_directml") is None:
            self.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        else:
            import torch_directml
            self.device = torch_directml.device()
        self.config = config
        self.clip_model = None
        self.clip_preprocess = None
        self.text_labels = None
        self.text_feature = None
        self.caption_model = None
        self.caption_vis_processors = None
        self.milvus_conn = milvus_conn
        if not config.ollama.enabled:
            # if ollama is not enabled, init local clip and caption model
            self.__init_clip_model__()
            self.__init_caption_model__()
        self.__init_embedding_model__()

    def __init_clip_model__(self):
        self.clip_model, self.clip_preprocess = load_from_name("ViT-B-16", device=self.device, download_root="./")
        text_labels = []
        with open("label_cn.txt") as f:
            text_labels = f.readlines()
        # remove whitespace characters like `\n` at the end of each line
        text_labels = [x.strip() for x in text_labels]
        self.text_labels = text_labels
        labels = clip.tokenize(text_labels).to(self.device)
        with torch.no_grad():
            labels = clip.tokenize(text_labels).to(self.device)
            text_feature = self.clip_model.encode_text(labels)
            text_feature /= text_feature.norm(dim=-1, keepdim=True)
            self.text_feature = text_feature

    def __init_caption_model__(self):
        self.caption_model, self.caption_vis_processors, _  = load_model_and_preprocess('blip_caption', model_type='base_coco', is_eval=True,  device=self.device)

    def __init_embedding_model__(self):
        self.embedding_processor = AutoImageProcessor.from_pretrained("google/vit-base-patch16-224")
        self.embedding_model = AutoModel.from_pretrained("google/vit-base-patch16-224").to(self.device)

    def label_image(self, path: str) -> list[ImageLabel]:
        if self.clip_model is None:
            self.__init_clip_model__()
        pil_image = Image.open(path)
        # 检查图像通道，如果为4通道，去掉alpha通道
        if pil_image.mode == "RGBA":
            pil_image = pil_image.convert("RGB")
        image = self.clip_preprocess(pil_image).unsqueeze(0).to(self.device)
        with torch.no_grad():
            image_features = self.clip_model.encode_image(image)
            image_features /= image_features.norm(dim=-1, keepdim=True)
        similarity = (100.0 * image_features @ self.text_feature.T).softmax(dim=-1)
        # find top 5 labels with highest similarity
        values, indices = similarity[0].topk(5)
        # print("Label and confidence:")
        predictions = []
        for value, index in zip(values, indices):
            predictions.append(ImageLabel(self.text_labels[index], value.item()))
        return predictions

    def caption_image(self, path: str) -> str:
        if self.caption_model is None:
            self.__init_caption_model__()
        pil_image = Image.open(path)
        # 检查图像通道，如果为4通道，去掉alpha通道
        if pil_image.mode == "RGBA":
            pil_image = pil_image.convert("RGB")
        return self.caption_model.generate({"image": self.caption_vis_processors['eval'](pil_image).unsqueeze(0).to(self.device)})[0]

    def image_embedding(self, path: str) -> np.ndarray:
        '''
        计算图像特征，并存入milvus
        '''
        image = Image.open(path)
        inputs = self.embedding_processor(image, return_tensors="pt").to(self.device)
        outputs = self.embedding_model(**inputs)
        embedding = outputs.pooler_output.cpu().detach().numpy().flatten()
        return embedding

    def image_understand_with_local_model(self, path: str):
        labels = self.label_image(path)
        logging.info("Image Labels: %s", labels)
        caption = self.caption_image(path)
        logging.info("Image Caption: %s", caption)
        return labels, caption
    
    def image_understand_with_ollama(self, path: str):
        image = Image.open(path)
        # resize long side to 1024
        width, height = image.size
        if width > height:
            if width > 1024:
                height = int(1024 * height / width)
                width = 1024
        if width <= height:
            if height > 1024:
                width = int(1024 * width / height)
                height = 1024
        image = image.resize((width, height))
        # 写入到byte数组
        buffered = BytesIO()
        image.save(buffered, format="JPEG")
        bts = buffered.getvalue()
        b64str = base64.b64encode(bts).decode("utf-8")
        prompt = '''You are an experienced art critic and photographer who specializes in evaluating works of art using simple and beautiful language.
Now, please use a short paragraph to describe the content of the picture you saw, and use this paragraph as the 'caption' in your answer.
After that, you are asked to give 3-5 words that summarize the image in a high level and are used to label the image you saw, these words will be used as 'tags' in your answer.
Finally, you will need to rate the image from four perspectives: 'Composition', 'Light and Shadow', 'Color' and 'Idea of the Work'. You need to rate the image from four perspectives: 'composition', 'light and shadow', 'color' and 'ideas', and give a final score of 0-10 on a scale of 0.1. The four scores and the overall rating will be used together as the 'scores' in your answer, and you will also be given a reason for why you gave the scores you gave from the four perspectives mentioned above, which will be used as the 'reason' for your answer.
Your answer needs to use the json format as a return, if you are not sure what you are seeing, please just return the empty Json object, e.g. '{}'. The content in your answer MUST be in 'Chinese', including 'caption', 'tags' and 'reason', any Non-Chinese answer will be considered as an invalid answer.
Your answer should not contain any subjective personal pronouns, e.g. 'I', 'we' etc. When you think you need to use them, please use words such as 'audience', 'others' etc. instead.'''
        format = {
        "type": "object",
        "properties": {
            "caption": {
                "type": "string"
            },
            "tags": {
                "type": "array",
                "items": {
                    "type": "string"
                }
            },
            "score": {
                "type": "array",
                "items": {
                    "type": "number",
                    "minimum": 0,
                    "maximum": 10
                }
            },
            "reason": {
                "type": "string"
            }
        },
        "required": ["caption", "tags", "score", "reason"]
    }
        result = OllamaClient(config=self.config.ollama).request_ollama_generate(body=prompt, image=[b64str], format=format)
        # 由ollama输出的标签没有置信度
        return [ImageLabel(x, 0.0) for x in result['tags']], result['caption']

    def understand(self, path: str) -> ImageUnderstandingResult:
        if self.config.ollama.enabled: 
            # using ollama
            labels, caption = self.image_understand_with_ollama(path)
        else:
            # using local model
            labels, caption = self.image_understand_with_local_model(path)
        embedding = self.image_embedding(path)
        self.milvus_conn.insert(embedding, path)
        return ImageUnderstandingResult(labels, caption)
    
    def image_similarity(self, path: str) -> list[dict]:
        embedding = self.image_embedding(path)
        records = self.milvus_conn.search_by_vec(embedding)
        results = []
        for record in records:
            results.append({"path": record.id, "score": record.distance})
        return results