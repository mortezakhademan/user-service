# simple user service
This project was created as sample for creating a user service.

---

## 🚀 Quick Start
1. Clone the repository:
   ```bash
   git clone https://github.com/mortezakhademan/simple-user-service.git
   cd simple-user-service
   ```

2. Create and configure the `.env` file and add these variables:
   1. HTTP_ADDRESS
   2. HTTP_PORT
   3. GRPC_PORT
   4. MONGO_URI
3. Run locally:

   ```bash
   go build -mod vendor -o main .
   ./main -env ./.env
   ```

   or run with Docker:

   ```bash
   docker-compose up -d
   ```

---

## ▶️ Running Locally

1. Build the project:

   ```bash
   go build -mod vendor -o main .
   ```
2. Run the project, passing the configs folder path with the `-config` parameter:

   ```bash
   ./main -env ./.env
   ```

---

## 🐳 Running with Docker

1. Mount the `.env` file as a volume in `docker-compose.yml`.
2. Set the MongoDB connection string in `.env`, for example:

   ```env
   MONGO_URI = "mongodb://localhost:27060/?authSource=admin"
   ```
3. Start the services:

   ```bash
   docker-compose up -d
   ```

---

## 🌐 APIs

1. **Swagger Documentation**
   Available at: [http://localhost:{HTTP_PORT}/api-docs](http://localhost:7020/api-docs)

2. **Get List of Users**
   Supports filters, sorting, pagination, and response type (JSON or Excel). Example:

   ```bash
   curl 'http://localhost:7020/users?page=1&pageSize=20&filters%5Bmobile%5D=09033596689' \
     -H 'accept: application/json' 
   ```
   