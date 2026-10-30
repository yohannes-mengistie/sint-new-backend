# Sint New Go API

The Go service owns Sint New application data and business rules. Firebase
Authenticates users; the API verifies Firebase ID tokens. Cloudinary stores
uploaded photos. PostgreSQL stores users, products, stores, and price updates.

## 1. Install Go

Install Go from the official website:

<https://go.dev/dl/>

Verify:

```sh
go version
```

## 2. Start PostgreSQL

Docker is the easiest local setup:

```sh
cd /home/yohannes/Documents/sint_new_backend
docker compose up -d postgres
docker compose ps
```

The first startup creates the `sint_new` database and applies
`migrations/001_initial.sql`.

Check the database:

```sh
docker compose exec postgres pg_isready -U sint_new -d sint_new
docker compose exec postgres psql -U sint_new -d sint_new -c "\dt"
```

The local connection string is:

```text
postgres://sint_new:sint_new_dev_password@localhost:5433/sint_new?sslmode=disable
```

## 3. Configure Firebase Admin SDK

In Firebase Console:

1. Open Project settings.
2. Open the Service accounts tab.
3. Click Generate new private key.
4. Save the JSON outside both repositories.

Do not commit this JSON file.

## 4. Configure Cloudinary

Create a Cloudinary account and copy the API environment variable from its
dashboard. Keep the value only in the backend `.env` file.

```sh
cp .env.example .env
```

The API loads this `.env` file automatically at startup. Existing shell
environment variables take precedence over values from the file.

Set:

```env
HTTP_ADDR=:8082
DATABASE_URL=postgres://sint_new:sint_new_dev_password@localhost:5433/sint_new?sslmode=disable
GOOGLE_APPLICATION_CREDENTIALS=/absolute/path/to/firebase-service-account.json
CLOUDINARY_URL=cloudinary://API_KEY:API_SECRET@CLOUD_NAME
```

## 5. Download dependencies and start the API

```sh
go mod tidy
go run ./cmd/api
```

Expected log:

```text
Sint New API listening on :8082
```

Test:

```sh
curl http://localhost:8082/health
```

Expected response:

```json
{"status":"ok"}
```

Authenticated API endpoints:

```text
GET /api/v1/me
GET /api/v1/products
POST /api/v1/price-updates
POST /api/v1/uploads/photos
```

These endpoints require a Firebase ID token:

```http
Authorization: Bearer FIREBASE_ID_TOKEN
```

`GET /api/v1/me` synchronizes the verified Firebase user into PostgreSQL.
`GET /api/v1/products` reads products from PostgreSQL. The upload endpoint
stores an authenticated multipart photo in Cloudinary.

The API applies the numbered SQL migrations from `migrations/` at startup.
The sample data migration is intentionally not applied automatically; run it
manually only in local development.

## Deploy on Render

The repository includes [Dockerfile](./Dockerfile) and
[render.yaml](./render.yaml) for a Render web service and PostgreSQL database.
Create a new Blueprint in Render from this repository. Set these secret
environment variables in the Render dashboard:

- `CLOUDINARY_URL`
- `FIREBASE_SERVICE_ACCOUNT_JSON_B64`

Create the Firebase secret value locally without committing the JSON file:

```sh
base64 -w 0 /absolute/path/to/firebase-service-account.json
```

Paste the output into `FIREBASE_SERVICE_ACCOUNT_JSON_B64`. Render supplies
`PORT` automatically; the API binds to it in production. After deployment,
copy the Render HTTPS service URL and run Flutter with:

```sh
flutter run --dart-define=API_BASE_URL=https://your-service.onrender.com
```

To load local sample products into an existing database, run:

```sh
docker compose exec -T postgres psql -U sint_new -d sint_new \
  < migrations/002_seed_products.sql
```

## Stop PostgreSQL

```sh
docker compose down
```

To remove the local database volume as well:

```sh
docker compose down -v
```

The second command permanently deletes the local PostgreSQL data volume.

## Architecture

The service follows a dependency-inward clean architecture:

```text
cmd/api
  -> delivery/http
  -> application
  -> infrastructure
  -> external services
```

- `internal/domain` contains business entities and no framework code.
- `internal/application` contains use cases and ports (interfaces).
- `internal/infrastructure` contains Firebase, Cloudinary, and PostgreSQL adapters.
- `internal/delivery/http` contains routing, middleware, request parsing, and responses.
- `cmd/api` composes the concrete adapters and starts the server.

The application layer depends on interfaces, not on Firebase, Cloudinary,
PostgreSQL, or `net/http`. This makes use cases testable and keeps provider
changes isolated to infrastructure packages.
