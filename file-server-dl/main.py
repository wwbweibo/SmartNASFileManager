from magika import Magika
from flask import Flask, request
from service.file_understanding import FileUnderstanding
from config import Config
import logging
import yaml

logging.basicConfig(level=logging.INFO)
config = Config().from_yaml_file("config.yaml")
print("===============")
print(config.__dict__)
print(config.ollama.__dict__)
fileUnderstanding = FileUnderstanding(config=config)
logging.info("File Understanding Service Started")

app = Flask(__name__)
# m = Magika()
@app.route('/api/v1/file/understanding', methods=['POST'])
def FileTypeInterfer():
    data = request.get_json()
    logging.info("File Understanding Request: %s", data['path'])
    try:
        result = fileUnderstanding.understand(data['path'])
        return result.to_dict()
    except Exception as e:
        logging.error("Error: %s", e)
        return {"error": str(e)}

@app.route('/api/v1/image/similar', methods=['POST'])
def ImageSimilar():
    data = request.get_json()
    logging.info("Image Similar Request: %s", data['path'])
    try:
        result = fileUnderstanding.similar(data['path'])
        return result.to_dict()
    except Exception as e:
        logging.error("Error: %s", e)
        return {"error": str(e)}


if __name__ == "__main__":
    # start an  http server here
    app.run(host='0.0.0.0', port='8081', debug=False, use_reloader=False)