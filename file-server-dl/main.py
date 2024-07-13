from magika import Magika
import socketserver
from flask import Flask, request
from service.file_interferer import interfer_file_type


app = Flask(__name__)
# m = Magika()
@app.route('/api/v1/file/interfer', methods=['POST'])
def FileTypeInterfer():
    data = request.get_json()
    result = interfer_file_type(data['path'])
    return {
        'type': result[0],
        'group': result[1],
        'description': result[2]
    }


if __name__ == "__main__":
    # start an  http server here
    app.run(host='0.0.0.0', port='8081', debug=True)