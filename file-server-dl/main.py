from magika import Magika
from flask import Flask, request
from service.file_understanding import FileUnderstanding
import logging

logging.basicConfig(level=logging.DEBUG)

fileUnderstanding = FileUnderstanding()
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

if __name__ == "__main__":
    # start an  http server here
    app.run(host='0.0.0.0', port='8081', debug=False, use_reloader=False)