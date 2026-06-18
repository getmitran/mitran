FROM python:3.11-slim
WORKDIR /app
COPY . .
RUN pip install --no-cache-dir -e .
EXPOSE 8888
CMD ["python", "-m", "flask", "--app", "worker.app", "run", "--host", "0.0.0.0", "--port", "8888"]
