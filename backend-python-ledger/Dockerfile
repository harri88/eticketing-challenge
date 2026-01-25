# Stage 1: Build dependencies
FROM python:3.11-slim as builder

WORKDIR /app

# Prevent Python from writing .pyc files and enable unbuffered logging
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Install build essentials for database drivers if needed
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

# Install dependencies to a local folder
COPY requirements.txt .
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

# Stage 2: Final lightweight image
FROM python:3.11-slim

WORKDIR /app

# Install runtime library for PostgreSQL (libpq)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq5 \
    && rm -rf /var/lib/apt/lists/*

# Copy installed dependencies from builder
COPY --from=builder /install /usr/local

# Copy application code
COPY . .

# Create a non-privileged user for security
RUN adduser --disabled-password --gecos "" ledgeruser
USER ledgeruser

EXPOSE 8000

# Run FastAPI using Uvicorn
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]