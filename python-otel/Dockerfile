FROM python:3.13.7-alpine3.22

WORKDIR /app

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY *.py .

EXPOSE 5000 8080

CMD opentelemetry-instrument flask run -h 0.0.0.0 -p 5000
