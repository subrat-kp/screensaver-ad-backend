# Screensaver Ad Backend

A Go backend service using the Gin framework for managing screensaver advertisements and creative assets.

## Features

- RESTful API for creative asset management
- File upload endpoint with S3 integration (placeholder)
- Asset listing and retrieval
- Health check endpoint
- CORS support for frontend integration

## Tech Stack

- **Go 1.21**
- **Gin Web Framework**
- **AWS S3** (for file storage - to be configured)

## API Endpoints

### Health Check
```
GET /health
```
Returns the service health status and timestamp.

**Response:**
```json
{
  "status": "ok",
  "timestamp": 1234567890
}
```

### Upload Creative
```
POST /api/upload
```
Upload a new creative asset (image/video).

**Request:**
- Content-Type: `multipart/form-data`
- Body: `file` (file upload)

**Response:**
```json
{
  "message": "File uploaded successfully (placeholder)",
  "creative": {
    "id": "20231012150405",
    "name": "example.jpg",
    "file_url": "https://s3.amazonaws.com/placeholder/example.jpg",
    "content_type": "image/jpeg",
    "uploaded_at": "2023-10-12T15:04:05Z"
  }
}
```

### List Assets
```
GET /api/assets
```
Retrieve all uploaded creative assets.

**Response:**
```json
{
  "assets": [...],
  "count": 5
}
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

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

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Docker Support

### Build Docker Image
```bash
docker build -t screensaver-ad-backend .
```

### Run Docker Container
```bash
docker run -p 8080:8080 screensaver-ad-backend
```

## Development

### Project Structure
```
.
├── main.go          # Main application entry point
├── go.mod           # Go module dependencies
├── Dockerfile       # Docker configuration
├── .gitignore       # Git ignore rules
└── README.md        # Project documentation
```

### TODO

- [ ] Implement actual AWS S3 integration
- [ ] Add database for persistent storage
- [ ] Implement authentication/authorization
- [ ] Add unit tests
- [ ] Add logging middleware
- [ ] Implement rate limiting
- [ ] Add environment configuration
- [ ] Create deployment scripts

## Configuration

S3 configuration and other environment variables will be added in future updates.

## Contributing

Feel free to open issues and submit pull requests.

## License

MIT License
