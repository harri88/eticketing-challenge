FROM python:3.11-slim

WORKDIR /app

# Copy requirements if they exist, otherwise create minimal setup
COPY requirements.txt* ./

RUN if [ -f requirements.txt ]; then pip install --no-cache-dir -r requirements.txt; else pip install --no-cache-dir psycopg2-binary flask flask-cors; fi

COPY . .

# Create a simple placeholder app if main.py doesn't exist
RUN if [ ! -f main.py ]; then echo 'from flask import Flask; app = Flask(__name__); @app.route("/health", methods=["GET"]); def health(): return {"status": "ok"}, 200; if __name__ == "__main__": app.run(host="0.0.0.0", port=8000)' > main.py; fi

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/health')"

CMD ["python", "main.py"]
