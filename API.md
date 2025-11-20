# API Documentation

Complete REST API documentation for the Weather Agency Label Tool.

## Base URL

- Development: `http://localhost:8080/api`
- Production: `https://your-domain.com/api`

## Content-Type

All API requests and responses use `application/json` content type unless otherwise specified.

## CORS

CORS is enabled for all origins in development. In production, configure as needed.

---

## Endpoints

### 1. Stations

#### Get All Stations

Retrieve list of all monitoring stations.

**Endpoint:** `GET /api/stations`

**Response:**
```json
[
  {
    "id": 1,
    "name": "北京市气象站",
    "longitude": 116.4074,
    "latitude": 39.9042
  },
  {
    "id": 2,
    "name": "上海市气象站",
    "longitude": 121.4737,
    "latitude": 31.2304
  }
]
```

**Status Codes:**
- `200 OK`: Success
- `500 Internal Server Error`: Database error

---

#### Get Nearest Station

Find the nearest monitoring station to given coordinates.

**Endpoint:** `GET /api/stations/nearest`

**Query Parameters:**
- `longitude` (required): Longitude coordinate (float)
- `latitude` (required): Latitude coordinate (float)

**Example:**
```
GET /api/stations/nearest?longitude=116.407396&latitude=39.904211
```

**Response:**
```json
{
  "id": 1,
  "name": "北京市气象站",
  "longitude": 116.4074,
  "latitude": 39.9042
}
```

**Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Invalid coordinates
- `500 Internal Server Error`: Database error

---

### 2. Images

#### Get All Images

Retrieve list of all images, sorted by annotation status (unannotated first).

**Endpoint:** `GET /api/images`

**Response:**
```json
[
  {
    "id": 1,
    "filename": "1700000000_fog_image.jpg",
    "filepath": "./uploads/1700000000_fog_image.jpg",
    "uploaded_at": "2024-01-01T12:00:00Z",
    "annotated": false
  },
  {
    "id": 2,
    "filename": "1700000001_ice_image.jpg",
    "filepath": "./uploads/1700000001_ice_image.jpg",
    "uploaded_at": "2024-01-01T12:30:00Z",
    "annotated": true
  }
]
```

**Status Codes:**
- `200 OK`: Success
- `500 Internal Server Error`: Database error

---

#### Get Image Details

Retrieve specific image with its annotation (if exists).

**Endpoint:** `GET /api/images/:id`

**Path Parameters:**
- `id` (required): Image ID (integer)

**Example:**
```
GET /api/images/1
```

**Response:**
```json
{
  "image": {
    "id": 1,
    "filename": "1700000000_fog_image.jpg",
    "filepath": "./uploads/1700000000_fog_image.jpg",
    "uploaded_at": "2024-01-01T12:00:00Z",
    "annotated": true
  },
  "annotation": {
    "id": 1,
    "image_id": 1,
    "category": "大雾",
    "severity": "中度",
    "observation_time": "2024-01-01T08:00:00Z",
    "location": "北京市朝阳区",
    "longitude": 116.407396,
    "latitude": 39.904211,
    "station_id": 1,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**Note:** If image has no annotation, `annotation` field will be `null`.

**Status Codes:**
- `200 OK`: Success
- `404 Not Found`: Image not found
- `500 Internal Server Error`: Database error

---

### 3. Annotations

#### Create or Update Annotation

Save or update annotation for an image.

**Endpoint:** `POST /api/annotations`

**Request Body:**
```json
{
  "image_id": 1,
  "category": "大雾",
  "severity": "中度",
  "observation_time": "2024-01-01T08:00:00Z",
  "location": "北京市朝阳区",
  "longitude": 116.407396,
  "latitude": 39.904211,
  "station_id": 1
}
```

**Field Descriptions:**
- `image_id`: ID of the image being annotated (required, integer)
- `category`: Weather phenomenon (required, enum: "大雾", "结冰", "积涝")
- `severity`: Severity level (required, enum: "无", "轻度", "中度", "重度")
- `observation_time`: Time of observation (required, ISO 8601 datetime string)
- `location`: Location name (required, string)
- `longitude`: Longitude coordinate (required, float)
- `latitude`: Latitude coordinate (required, float)
- `station_id`: Monitoring station ID (required, integer)

**Response:**
```json
{
  "id": 1,
  "image_id": 1,
  "category": "大雾",
  "severity": "中度",
  "observation_time": "2024-01-01T08:00:00Z",
  "location": "北京市朝阳区",
  "longitude": 116.407396,
  "latitude": 39.904211,
  "station_id": 1,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

**Behavior:**
- If annotation exists for the image, it will be updated
- If no annotation exists, a new one will be created
- Image's `annotated` flag is automatically set to `true`

**Status Codes:**
- `201 Created`: Annotation created/updated successfully
- `400 Bad Request`: Invalid request body
- `500 Internal Server Error`: Database error

---

### 4. Upload

#### Upload Image

Upload a new image file.

**Endpoint:** `POST /api/upload`

**Content-Type:** `multipart/form-data`

**Form Data:**
- `image` (required): Image file (JPEG, PNG, GIF)

**Example (cURL):**
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "image=@/path/to/image.jpg"
```

**Response:**
```json
{
  "id": 3,
  "filename": "1700000002_image.jpg",
  "filepath": "./uploads/1700000002_image.jpg",
  "uploaded_at": "2024-01-01T13:00:00Z",
  "annotated": false
}
```

**Notes:**
- Max file size: 32MB
- Filename is automatically generated with timestamp prefix
- File is stored in `./uploads/` directory
- Uploaded image is immediately added to database

**Status Codes:**
- `201 Created`: Image uploaded successfully
- `400 Bad Request`: No file provided or invalid file
- `500 Internal Server Error`: File system or database error

---

### 5. Image Files

#### Serve Image File

Retrieve the actual image file.

**Endpoint:** `GET /images/:filename`

**Path Parameters:**
- `filename` (required): Image filename (string)

**Example:**
```
GET /images/1700000000_fog_image.jpg
```

**Response:**
- Content-Type: `image/jpeg`, `image/png`, etc.
- Binary image data

**Status Codes:**
- `200 OK`: Image found and served
- `404 Not Found`: Image file not found

---

## Data Models

### Station
```typescript
{
  id: number
  name: string
  longitude: number  // Decimal degrees
  latitude: number   // Decimal degrees
}
```

### Image
```typescript
{
  id: number
  filename: string
  filepath: string
  uploaded_at: string  // ISO 8601 datetime
  annotated: boolean
}
```

### Annotation
```typescript
{
  id: number
  image_id: number
  category: "大雾" | "结冰" | "积涝"
  severity: "无" | "轻度" | "中度" | "重度"
  observation_time: string  // ISO 8601 datetime
  location: string
  longitude: number  // Decimal degrees
  latitude: number   // Decimal degrees
  station_id: number
  created_at: string  // ISO 8601 datetime
  updated_at: string  // ISO 8601 datetime
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

**Common Error Status Codes:**
- `400 Bad Request`: Invalid request parameters or body
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

---

## Rate Limiting

Currently, no rate limiting is implemented. Consider adding rate limiting in production using Nginx or application-level middleware.

---

## Authentication

Currently, no authentication is implemented. For production use, consider adding:
- JWT-based authentication
- OAuth 2.0
- API keys
- Session-based authentication

---

## Examples

### JavaScript (Fetch API)

```javascript
// Get all stations
const stations = await fetch('http://localhost:8080/api/stations')
  .then(res => res.json());

// Get nearest station
const nearest = await fetch(
  'http://localhost:8080/api/stations/nearest?longitude=116.4&latitude=39.9'
).then(res => res.json());

// Save annotation
const annotation = await fetch('http://localhost:8080/api/annotations', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    image_id: 1,
    category: '大雾',
    severity: '中度',
    observation_time: new Date().toISOString(),
    location: '北京市',
    longitude: 116.4,
    latitude: 39.9,
    station_id: 1
  })
}).then(res => res.json());

// Upload image
const formData = new FormData();
formData.append('image', fileInput.files[0]);
const uploaded = await fetch('http://localhost:8080/api/upload', {
  method: 'POST',
  body: formData
}).then(res => res.json());
```

### Python (Requests)

```python
import requests

# Get all stations
stations = requests.get('http://localhost:8080/api/stations').json()

# Get nearest station
params = {'longitude': 116.4, 'latitude': 39.9}
nearest = requests.get(
    'http://localhost:8080/api/stations/nearest',
    params=params
).json()

# Save annotation
annotation_data = {
    'image_id': 1,
    'category': '大雾',
    'severity': '中度',
    'observation_time': '2024-01-01T08:00:00Z',
    'location': '北京市',
    'longitude': 116.4,
    'latitude': 39.9,
    'station_id': 1
}
annotation = requests.post(
    'http://localhost:8080/api/annotations',
    json=annotation_data
).json()

# Upload image
files = {'image': open('fog_image.jpg', 'rb')}
uploaded = requests.post(
    'http://localhost:8080/api/upload',
    files=files
).json()
```

### cURL

```bash
# Get all stations
curl http://localhost:8080/api/stations

# Get nearest station
curl "http://localhost:8080/api/stations/nearest?longitude=116.4&latitude=39.9"

# Save annotation
curl -X POST http://localhost:8080/api/annotations \
  -H "Content-Type: application/json" \
  -d '{
    "image_id": 1,
    "category": "大雾",
    "severity": "中度",
    "observation_time": "2024-01-01T08:00:00Z",
    "location": "北京市",
    "longitude": 116.4,
    "latitude": 39.9,
    "station_id": 1
  }'

# Upload image
curl -X POST http://localhost:8080/api/upload \
  -F "image=@fog_image.jpg"
```

---

## Performance Considerations

- Database connections are pooled (max 10 open, 5 idle)
- Stations are loaded from database each request (consider caching)
- Image files are served directly (consider CDN for production)
- No pagination on image list (add if dataset grows large)

---

## Future Enhancements

Potential API improvements:
- Pagination for image and station lists
- Filtering and search capabilities
- Bulk operations (upload multiple, annotate multiple)
- Image compression and thumbnail generation
- Export annotations (JSON, CSV)
- Statistics and analytics endpoints
- WebSocket for real-time updates
- Audit log for annotation changes
