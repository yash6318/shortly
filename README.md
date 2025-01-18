# Shortly

This project is called Shortly, a URL shortener application with a frontend built using React and a backend built using Go. The backend uses PostgreSQL as the database.
## Features
- Shorten long URLs
- Redirect short URLs to their original destinations
- Health check endpoints
- Dockerized setup for easy deployment

## Setup

### Development Environment

1. **Clone the repository**:
   ```sh
   git clone https://github.com/yash6318/shortly.git
   cd shortly
   ```

2. **Create a `.env` file in the `url-shortener-backend` directory**:

    ```sh
    DB_HOST=db
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=yourpassword
    DB_NAME=shortener
    BASE_URL=http://localhost:8080
    ```

3. **Run the application using Docker Compose**:

    ```sh
    docker-compose up --build
    ```

4. **Access the application**:
    - Frontend: http://localhost:3000
    - Backend: http://localhost:8080

### Production Environment

1. **Create a `.env` file in the `url-shortener-backend` directory**:

    ```sh
    DB_HOST=db
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=yourpassword
    DB_NAME=shortener
    BASE_URL=http://localhost:8080
    ```

2. **Run the application using Docker Compose with production configuration**:

    ```sh
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
    ```

3. **Access the application**:
    - Frontend: http://your-production-url
    - Backend: http://your-production-url:8080

### Docker Hub Images

If you prefer to use pre-built Docker images from Docker Hub, you can use the following configuration:

```sh
version: "3.9"

services:
  app:
    image: yash6318/shortly-backend:latest
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - BASE_URL=http://localhost:8080
      - CORS_ALLOWED_ORIGINS=http://your-production-frontend.com
    depends_on:
      db:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3
  
  frontend:
    image: yash6318/shortly-frontend:latest
    ports:
      - "3000:3000"
    environment:
      - REACT_APP_API_URL=http://localhost:8080
    depends_on:
      app:
        condition: service_healthy
  
  db:
    image: postgres:15
    container_name: postgres-db
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    ports:
      - "5432:5432"
    volumes:
      - db-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "${DB_USER}", "-d", "${DB_NAME}"]
      interval: 10s
      retries: 5
      timeout: 5s

volumes:
  db-data:
```
This setup will allow users to run the application in both development and production environments using Docker Compose.