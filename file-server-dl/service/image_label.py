import cn_clip.clip as clip
from cn_clip.clip import load_from_name, available_models
from PIL import Image
import torch

device = 'cuda' if torch.cuda.is_available() else 'cpu'

model, preprocess = load_from_name("ViT-B-16", device="mps", download_root="./")
text_labels = []
with open("label_cn.txt") as f:
    text_labels = f.readlines()
# remove whitespace characters like `\n` at the end of each line
text_labels = [x.strip() for x in text_labels]
labels = clip.tokenize(text_labels).to("mps")
with torch.no_grad():
    labels = clip.tokenize(text_labels).to("mps")
    text_feature = model.encode_text(labels)
    text_feature /= text_feature.norm(dim=-1, keepdim=True)

def label_image(path: str) -> list[dict]:
    image = preprocess(Image.open(path)).unsqueeze(0).to(device)
    with torch.no_grad():
        image_features = model.encode_image(image)
        image_features /= image_features.norm(dim=-1, keepdim=True)
    similarity = (100.0 * image_features @ text_feature.T).softmax(dim=-1)
    # find top 5 labels with highest similarity
    values, indices = similarity[0].topk(5)
    # print("Label and confidence:")
    predictions = []
    for value, index in zip(values, indices):
        predictions.append({
            'label': text_labels[index],
            'confidence': f"{100 * value.item():.2f}%"
        })
    return predictions