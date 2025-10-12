# Screensaver Ad Backend

A Go backend service using the Gin framework for managing screensaver advertisements and creative assets with PostgreSQL database support.

## Features

- RESTful API for creative asset management
- PostgreSQL database integration for asset metadata persistence
- File upload endpoint with S3 integration (placeholder)
- Asset status tracking (processing/processed)
- Automatic status checking via S3 processed output folder
- Asset listing and retrieval
- Health check endpoint with database status
- CORS support for frontend integration
- Backward compatibility with in-memory storage (fallback mode)

## Tech Stack

- **Go 1.21**
- **Gin Web Framework**
- **PostgreSQL** (for persistent asset metadata storage)
- **AWS S3** (for file storage - to be configured)
- **lib/pq** (PostgreSQL driver)
- **AWS SDK for Go** (S3 integration)

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

The application automatically creates the `asset_metadata` table on startup with the following schema:

```sql
CREATE TABLE IF NOT EXISTS asset_metadata (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    file_url TEXT NOT NULL,
    content_type VARCHAR(100),
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'processing'
);
```

### 3. Environment Variables

Configure the following environment variables to connect to your PostgreSQL database:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=screensaver_ad

# AWS S3 Configuration (optional)
export AWS_REGION=us-east-1
export S3_BUCKET=your-bucket-name
export S3_PROCESSED_FOLDER=processed/
```

### 4. Database Connection

The application connects to PostgreSQL using the environment variables. If the database is not available, it falls back to in-memory storage mode.

**Connection String Format:**
```
host=<DB_HOST> port=<DB_PORT> user=<DB_USER> password=<DB_PASSWORD> dbname=<DB_NAME> sslmode=disable
```

### 5. Default Values

If environment variables are not set, the following defaults are used:

- `DB_HOST`: localhost
- `DB_PORT`: 5432
- `DB_USER`: postgres
- `DB_PASSWORD`: postgres
- `DB_NAME`: screensaver_ad

## API Endpoints

### Health Check

```
GET /health
```

Returns the service health status, timestamp, and database connection status.

**Response:**
```json
{
  "status": "ok",
  "timestamp": 1234567890,
  "database": "connected"
}
```

### Upload Creative

```
POST /api/upload
```

Upload a new creative asset (image/video). Creates a database record with status="processing".

**Request:**
- Content-Type: `multipart/form-data`
- Body: `file` (file upload)

**Response:**
```json
{
  "message": "File uploaded successfully",
  "creative": {
    "id": "20231012150405",
    "name": "example.jpg",
    "file_url": "https://s3.amazonaws.com/placeholder/example.jpg",
    "content_type": "image/jpeg",
    "uploaded_at": "2023-10-12T15:04:05Z",
    "status": "processing"
  }
}
```

### List All Assets

```
GET /api/assets
```

Retrieve all creative assets from the database.

**Response:**
```json
{
  "assets": [
    {
      "id": "20231012150405",
      "name": "example.jpg",
      "file_url": "https://s3.amazonaws.com/placeholder/example.jpg",
      "content_type": "image/jpeg",
      "uploaded_at": "2023-10-12T15:04:05Z",
      "status": "processing"
    }
  ],
  "count": 1
}
```

### Get Asset Status

```
GET /api/assets/:id/status
```

Get the current status of an asset. If status is "processing", checks S3 processed output folder and updates status to "processed" if file exists.

**Response:**
```json
{
  "status": "processed",
  "file_url": "https://s3.amazonaws.com/placeholder/example.jpg",
  "message": "Asset example.jpg is processed"
}
```

**Status Values:**
- `processing`: Asset is being processed
- `processed`: Asset processing is complete and file is available in S3

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher (optional, falls back to in-memory mode)
- AWS credentials configured (for S3 integration)

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
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=screensaver_ad
```

5. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Docker Support

### Build Docker Image

```bash
docker build -t screensaver-ad-backend .
```

### Run with Docker Compose (with PostgreSQL)

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
  screensaver-ad-backend
```

## Development

### Project Structure

```
.
├── main.go          # Main application entry point
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── Dockerfile       # Docker configuration
├── .gitignore       # Git ignore rules
└── README.md        # Project documentation
```

### Database Schema Details

**asset_metadata table:**
- `id` (VARCHAR(255), PRIMARY KEY): Unique identifier for the asset
- `name` (VARCHAR(500), NOT NULL): Original filename
- `file_url` (TEXT, NOT NULL): S3 URL or storage location
- `content_type` (VARCHAR(100)): MIME type of the file
- `uploaded_at` (TIMESTAMP, NOT NULL): Upload timestamp
- `status` (VARCHAR(50), NOT NULL, DEFAULT 'processing'): Current processing status

### Features Implemented

- ✅ PostgreSQL database integration
- ✅ Asset metadata table with status tracking
- ✅ Status endpoint with S3 processed output checking
- ✅ Automatic status updates (processing → processed)
- ✅ Backward compatibility with in-memory storage
- ✅ Health check with database status

### TODO

- [ ] Implement actual AWS S3 file upload
- [ ] Add file validation (size, type)
- [ ] Implement authentication/authorization
- [ ] Add unit tests
- [ ] Add logging middleware
- [ ] Implement rate limiting
- [ ] Add pagination for asset listing
- [ ] Create deployment scripts
- [ ] Add database migrations
- [ ] Implement soft delete functionality

## Configuration

### Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | PostgreSQL host | localhost | No |
| `DB_PORT` | PostgreSQL port | 5432 | No |
| `DB_USER` | Database username | postgres | No |
| `DB_PASSWORD` | Database password | postgres | No |
| `DB_NAME` | Database name | screensaver_ad | No |
| `AWS_REGION` | AWS region for S3 | us-east-1 | Yes (for S3) |
| `S3_BUCKET` | S3 bucket name | - | Yes (for S3) |
| `S3_PROCESSED_FOLDER` | S3 processed output folder | processed/ | No |

## Frontend Integration

### Example: Fetching Asset Status

When an asset page is opened in the frontend, use the status endpoint to check if processing is complete:

```javascript
// Example JavaScript code for frontend
async function checkAssetStatus(assetId) {
  try {
    const response = await fetch(`http://localhost:8080/api/assets/${assetId}/status`);
    const data = await response.json();
    
    if (data.status === 'processed') {
      console.log('Asset is ready:', data.file_url);
      // Display the processed asset
    } else if (data.status === 'processing') {
      console.log('Asset is still processing, checking again in 5 seconds...');
      // Poll again after a delay
      setTimeout(() => checkAssetStatus(assetId), 5000);
    }
  } catch (error) {
    console.error('Failed to fetch asset status:', error);
  }
}

// Usage
checkAssetStatus('20231012150405');
```

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

5. The application will fall back to in-memory storage if database connection fails

### Common Errors

**Error: "pq: database 'screensaver_ad' does not exist"**
- Solution: Create the database using `CREATE DATABASE screensaver_ad;`

**Error: "pq: password authentication failed"**
- Solution: Verify `DB_PASSWORD` environment variable matches your PostgreSQL user password

**Error: "connection refused"**
- Solution: Ensure PostgreSQL is running and accessible on the specified host/port

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
