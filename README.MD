
# 🌤️ Weather API

A simple API that allows users to subscribe for weather updates by email. Supports selecting a city and frequency of updates.

---

## 🚀 Getting Started

### 📦 Requirements

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## 🛠️ Build and Run the API

1. **Build the Docker images**:
   ```bash
   docker-compose build
   ```

2. **Start the containers in the background**:
   ```bash
   docker-compose up -d
   ```

   The API will be accessible at:

   ```
   http://localhost:PORT  # Replace PORT with the actual port used by your API (e.g., 8000 or 3000)
   ```

---

## 🧪 Example Request(local)

**Subscribe to weather updates:**

```bash
curl -X 'POST' \
  'http://localhost:PORT/api/subscribe' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "your_email@example.com",
    "city": "Kyiv",
    "frequency": "hourly"
  }'
```

---

## 🛑 Stopping the API

To stop and remove the containers:

```bash
docker-compose down
```

---

## 🧹 Removing Database Volume (if needed)

1. **List volumes**:
   ```bash
   docker volume ls
   ```

2. **Remove the volume (example)**:
   ```bash
   docker volume rm awesomeproject_postgres-data
   ```

Replace `awesomeproject_postgres-data` with your actual volume name if it differs.

