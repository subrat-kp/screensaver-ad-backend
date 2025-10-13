# Screensaver Ad Backend

A Go backend service using the Gin framework for managing screensaver advertisements and creative assets with PostgreSQL database and AWS S3 storage integration.

## Features

- RESTful API for creative asset management
- PostgreSQL database integration with GORM ORM
- AWS S3 file upload and storage
- Multipart form data handling for file uploads
- Asset metadata persistence
- Pagination support for asset listing
- Health check endpoint with database status
- CORS support for frontend integration
- Environment-based configuration with .env support

## Tech Stack

- **Go 1.23**
- **Gin Web Framework** - HTTP web framework
- **GORM** - ORM library for database operations
- **PostgreSQL** - Relational database for asset metadata
- **AWS SDK for Go** - S3 integration for file storage
- **godotenv** - Environment variable management
- **google/uuid** - UUID generation for unique file naming

## Architecture

```
.
├── config/
│   ├── database.go      # Database configuration and initialization
│   └── s3.go           # S3 client configuration
├── internal/
│   ├── controllers/
│   │   └── asset_controller.go  # HTTP request handlers
│   ├── models/
│   │   └── asset.go            # Asset data model
│   ├── repository/
│   │   └── asset_repository.go # Database operations
│   └── services/
│       ├── asset_service.go    # Business logic
│       └── s3_service.go       # S3 operations
├── main.go             # Application entry point
├── go.mod              # Go module dependencies
├── .env.example        # Environment variables template
├── Dockerfile          # Docker configuration
└── README.md           # Project documentation
```

## Database Setup

### Prerequisites

- PostgreSQL 12 or higher installed
- Database user with CREATE TABLE permissions

### 1. Create Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE screensaver_ad;

# Connect to the database
\c screensaver_ad
```

### 2. Database Schema

The application automatically creates the `asset_metadata` table on startup using GORM migrations:

```sql
CREATE TABLE asset_metadata (
    id SERIAL PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    s3_key VARCHAR(500) NOT NULL UNIQUE,
    s3_bucket VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'uploaded',
    uploaded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);
```

### 3. Environment Variables

Create a `.env` file based on `.env.example`:

```bash
# Development environment flag
GO_ENV=dev

# DATABASE
DB_HOST=your_host_url
DB_PORT=your_db_server_port
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_db_name
DB_SSLMODE=disable

# AWS S3
AWS_REGION=your_region
AWS_ACCESS_KEY_ID=your_access_key_id
AWS_SECRET_ACCESS_KEY=your_secret_access_key
AWS_S3_BUCKET=your_s3_bucket_name
```

### 4. AWS S3 Setup

1. Create an S3 bucket
2. Create an `input/` folder inside the bucket
3. Configure IAM user with S3 permissions:
   - `s3:PutObject`
   - `s3:GetObject`
   - `s3:DeleteObject`
4. Generate access keys for the IAM user
5. Add credentials to `.env` file

## API Endpoints

### Health Check

```
GET /health
```

Returns the service health status and database connection status.

**Response:**
```json
{
  "status": "ok",
  "database": "connected"
}
```

### Create Asset with File Upload

```
POST /api/assets
```

Upload a new creative asset (image or video) to S3 and create a database record.

**Request:**
- Content-Type: `multipart/form-data`
- Fields:
  - `file` (required): The image or video file to upload
  - `name` (optional): Custom name for the asset (defaults to original filename)

**Supported File Types:**
- Images: JPEG, JPG, PNG, GIF, WebP
- Videos: MP4, MPEG, QuickTime, AVI, WebM

**Maximum File Size:** 32 MB

**Example using curl:**
```bash
curl -X POST http://localhost:8080/api/assets \
  -F "file=@/path/to/image.jpg" \
  -F "name=my-creative-asset"
```

**Response:**
```json
{
  "message": "Asset created successfully",
  "asset": {
    "id": 1,
    "file_name": "my-creative-asset",
    "file_size": 123456,
    "content_type": "image/jpeg",
    "s3_key": "input/my-creative-asset_a1b2c3d4.jpg",
    "s3_bucket": "screensaver-creatives",
    "status": "uploaded",
    "uploaded_at": "2025-10-13T10:40:00Z",
    "created_at": "2025-10-13T10:40:00Z",
    "updated_at": "2025-10-13T10:40:00Z"
  }
}
```

### List All Assets

```
GET /api/assets?limit=10&offset=0
```

Retrieve all creative assets from the database with pagination.

**Query Parameters:**
- `limit` (optional): Number of assets to return (default: 10, max: 100)
- `offset` (optional): Number of assets to skip (default: 0)

**Response:**
```json
{
  "assets": [
    {
      "id": 1,
      "file_name": "my-creative-asset",
      "file_size": 123456,
      "content_type": "image/jpeg",
      "s3_key": "input/my-creative-asset_a1b2c3d4.jpg",
      "s3_bucket": "screensaver-creatives",
      "status": "uploaded",
      "uploaded_at": "2025-10-13T10:40:00Z",
      "created_at": "2025-10-13T10:40:00Z",
      "updated_at": "2025-10-13T10:40:00Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### Get Single Asset

```
GET /api/assets/:id
```

Retrieve a specific asset by ID.

**Response:**
```json
{
  "id": 1,
  "file_name": "my-creative-asset",
  "file_size": 123456,
  "content_type": "image/jpeg",
  "s3_key": "input/my-creative-asset_a1b2c3d4.jpg",
  "s3_bucket": "screensaver-creatives",
  "status": "uploaded",
  "uploaded_at": "2025-10-13T10:40:00Z",
  "created_at": "2025-10-13T10:40:00Z",
  "updated_at": "2025-10-13T10:40:00Z"
}
```

### Update Asset

```
PUT /api/assets/:id
```

Update an existing asset's metadata.

**Request Body:**
```json
{
  "file_name": "updated-name"
}
```

**Response:**
```json
{
  "id": 1,
  "file_name": "updated-name",
  "file_size": 123456,
  "content_type": "image/jpeg",
  "s3_key": "input/my-creative-asset_a1b2c3d4.jpg",
  "s3_bucket": "screensaver-creatives",
  "status": "uploaded",
  "uploaded_at": "2025-10-13T10:40:00Z",
  "created_at": "2025-10-13T10:40:00Z",
  "updated_at": "2025-10-13T10:45:00Z"
}
```

### Delete Asset

```
DELETE /api/assets/:id
```

Delete an asset from the database (soft delete).

**Response:**
```json
{
  "message": "Asset deleted successfully"
}
```

### Update Asset Status

```
PATCH /api/assets/:id/status
```

Update the processing status of an asset.

**Request Body:**
```json
{
  "status": "processed"
}
```

**Valid Status Values:**
- `uploaded` - Asset has been uploaded to S3 (default status on creation)
- `processed` - Asset has been processed and is ready for use

**Example using curl:**
```bash
curl -X PATCH http://localhost:8080/api/assets/1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "processed"}'
```

**Response:**
```json
{
  "message": "Asset status updated successfully",
  "asset": {
    "id": 1,
    "file_name": "my-creative-asset",
    "file_size": 123456,
    "content_type": "image/jpeg",
    "s3_key": "input/my-creative-asset_a1b2c3d4.jpg",
    "s3_bucket": "screensaver-creatives",
    "status": "processed",
    "uploaded_at": "2025-10-13T10:40:00Z",
    "created_at": "2025-10-13T10:40:00Z",
    "updated_at": "2025-10-13T11:50:00Z"
  }
}
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 12 or higher
- AWS account with S3 access
- AWS credentials configured

### Installation

1. Clone the repository:
```bash
git clone https://github.com/subrat-kp/screensaver-ad-backend.git
cd screensaver-ad-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database (see Database Setup section)

4. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your actual values
```

5. Run the application:
```bash
GO_ENV=dev go run main.go
```

The server will start on `http://localhost:8080`

### Development Mode

The application loads `.env` file only when `GO_ENV` is set to `dev` or `development`:

```bash
GO_ENV=dev go run main.go
```

For production, set environment variables directly without using `.env` file.

## Docker Support

### Build Docker Image

```bash
docker build -t screensaver-ad-backend .
```

### Run with Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: screensaver_ad
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: screensaver_ad
      DB_SSLMODE: disable
      AWS_REGION: us-east-1
      AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
      AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}
      AWS_S3_BUCKET: screensaver-creatives
    depends_on:
      - postgres

volumes:
  postgres_data:
```

Run with:
```bash
docker-compose up
```

### Run Docker Container (standalone)

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=screensaver_ad \
  -e DB_SSLMODE=disable \
  -e AWS_REGION=us-east-1 \
  -e AWS_ACCESS_KEY_ID=your_key \
  -e AWS_SECRET_ACCESS_KEY=your_secret \
  -e AWS_S3_BUCKET=screensaver-creatives \
  screensaver-ad-backend
```

## Configuration

### Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GO_ENV` | Environment mode (dev/production) | - | No |
| `DB_HOST` | PostgreSQL host | localhost | Yes |
| `DB_PORT` | PostgreSQL port | 5432 | Yes |
| `DB_USER` | Database username | postgres | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | screensaver_ad | Yes |
| `DB_SSLMODE` | SSL mode (disable/require) | disable | No |
| `AWS_REGION` | AWS region for S3 | - | Yes |
| `AWS_ACCESS_KEY_ID` | AWS access key | - | Yes |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | - | Yes |
| `AWS_S3_BUCKET` | S3 bucket name | - | Yes |

## Features Implemented

- ✅ PostgreSQL database integration with GORM
- ✅ AWS S3 file upload and storage
- ✅ Multipart form data handling
- ✅ Asset metadata persistence
- ✅ File type validation (images and videos)
- ✅ Unique file naming with UUID
- ✅ Pagination support
- ✅ CRUD operations for assets
- ✅ Health check endpoint
- ✅ CORS support
- ✅ Environment-based configuration
- ✅ Graceful S3 initialization (optional)
- ✅ Transaction rollback on failures
- ✅ Asset status tracking (uploaded/processed)
- ✅ Status update endpoint

## Security Best Practices

1. **Never commit credentials** - `.env` file is in `.gitignore`
2. **Use IAM roles** in production instead of access keys
3. **Rotate credentials regularly**
4. **Use minimal IAM permissions** - only grant necessary S3 access
5. **Enable SSL/TLS** for database connections in production
6. **Validate file types and sizes** before upload
7. **Use presigned URLs** for temporary file access

## Troubleshooting

### Database Connection Issues

If the application fails to connect to PostgreSQL:

1. Check if PostgreSQL is running:
```bash
psql -U postgres -h localhost
```

2. Verify environment variables are set correctly

3. Check PostgreSQL logs for connection errors

4. Ensure the database user has proper permissions

5. Verify SSL mode matches your PostgreSQL configuration

### S3 Upload Issues

**Error: "S3 client is not initialized"**
- Solution: Verify all AWS environment variables are set correctly

**Error: "failed to upload to S3: AccessDenied"**
- Solution: Check IAM user permissions for S3 bucket access

**Error: "failed to upload to S3: NoSuchBucket"**
- Solution: Ensure the S3 bucket exists and the name is correct

### Common Errors

**Error: "pq: database 'screensaver_ad' does not exist"**
- Solution: Create the database using `CREATE DATABASE screensaver_ad;`

**Error: "pq: password authentication failed"**
- Solution: Verify `DB_PASSWORD` environment variable matches your PostgreSQL user password

**Error: "connection refused"**
- Solution: Ensure PostgreSQL is running and accessible on the specified host/port

**Error: "invalid file type"**
- Solution: Only images (JPEG, PNG, GIF, WebP) and videos (MP4, MPEG, QuickTime, AVI, WebM) are supported

## Contributing

Feel free to open issues and submit pull requests.

### Development Guidelines

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License
