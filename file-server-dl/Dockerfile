from python:3.10
WORKDIR /app
ADD . .
RUN pip install -r requirements.txt
# RUN https_proxy=http://192.168.163.65:10080 pip install -r requirements.txt
CMD ["python", "main.py"]