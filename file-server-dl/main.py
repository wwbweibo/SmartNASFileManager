from magika import Magika
from flask import Flask, request
from service.file_understanding import FileUnderstanding

fileUnderstanding = FileUnderstanding()

app = Flask(__name__)
# m = Magika()
@app.route('/api/v1/file/understanding', methods=['POST'])
def FileTypeInterfer():
    data = request.get_json()
    return fileUnderstanding.understand(data['path'])

if __name__ == "__main__":
    # start an  http server here
    app.run(host='0.0.0.0', port='8081', debug=True)