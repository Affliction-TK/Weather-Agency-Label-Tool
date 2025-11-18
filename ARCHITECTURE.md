# System Architecture

This document describes the architecture and design decisions of the Weather Agency Label Tool.

## Overview

The Weather Agency Label Tool is a full-stack web application for annotating weather monitoring images. It follows a traditional three-tier architecture:

```
┌─────────────────────────────────────────────────────────┐
│                      Frontend (Svelte)                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Image List   │  │ Annotation   │  │ Upload Tab   │  │
│  │ Component    │  │ Form         │  │              │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │ HTTP/REST API
                         ▼
┌─────────────────────────────────────────────────────────┐
│                   Backend (Go + Gorilla Mux)            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ API Handlers │  │ Business     │  │ Static File  │  │
│  │              │  │ Logic        │  │ Serving      │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │ SQL Queries
                         ▼
┌─────────────────────────────────────────────────────────┐
│                     Database (MySQL)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │  stations    │  │   images     │  │ annotations  │  │
│  │              │  │              │  │              │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                  File System (uploads/)                 │
│               Image files stored locally                │
└─────────────────────────────────────────────────────────┘
```

## Components

### Frontend (Svelte)

**Technology:** Svelte 5 + Vite

**Purpose:** User interface for image annotation

**Key Components:**

1. **App.svelte** - Main application shell
   - Manages application state
   - Handles routing between tabs
   - Coordinates data flow between components

2. **ImageList.svelte** - Sidebar image list
   - Displays thumbnails with annotation status
   - Provides image selection
   - Shows visual feedback for current selection

3. **AnnotationForm.svelte** - Main annotation interface
   - Image preview
   - Form inputs for all annotation fields
   - Station auto-suggestion based on coordinates
   - Form validation and submission

4. **UploadTab.svelte** - Image upload interface
   - Drag-and-drop file upload
   - Multiple file selection
   - Upload progress indication

**Data Flow:**
```
User Action → Component Event → App State Update → API Call → UI Update
```

**State Management:**
- Component-level state using Svelte stores
- Props for parent-child communication
- Events for child-parent communication

### Backend (Go)

**Technology:** Go 1.24 + Gorilla Mux

**Purpose:** REST API server, business logic, and file serving

**Key Components:**

1. **HTTP Server**
   - Gorilla Mux router for flexible routing
   - CORS middleware for cross-origin requests
   - Static file serving for frontend and images

2. **API Handlers**
   - `getStations()` - List all monitoring stations
   - `getNearestStation()` - Find nearest station by coordinates
   - `getImages()` - List all images with status
   - `getImage()` - Get specific image with annotation
   - `createAnnotation()` - Create or update annotation
   - `uploadImage()` - Handle image upload
   - `serveImage()` - Serve image files

3. **Business Logic**
   - Haversine distance calculation for station matching
   - Image file management
   - Annotation upsert logic (create or update)

4. **Database Layer**
   - SQL prepared statements for security
   - Connection pooling for performance
   - Transaction handling for data consistency

**Design Patterns:**
- Handler functions follow standard http.HandlerFunc signature
- Separation of concerns (handlers, logic, data access)
- RESTful API design principles

### Database (MySQL)

**Technology:** MySQL 8.0 with utf8mb4 encoding

**Purpose:** Persistent data storage

**Schema Design:**

1. **stations** table
   - Stores monitoring station information
   - Indexed on coordinates for spatial queries
   - Pre-populated with 150 stations

2. **images** table
   - Tracks uploaded images and their status
   - Filename includes timestamp to prevent collisions
   - Boolean flag for annotation status

3. **annotations** table
   - Stores annotation data for each image
   - Foreign keys enforce referential integrity
   - Unique constraint on image_id (one annotation per image)
   - Automatic timestamp tracking

**Relationships:**
```
stations (1) ──< (N) annotations
images   (1) ──< (1) annotations
```

**Indexes:**
- Primary keys on all tables
- Foreign key indexes for joins
- Coordinate index on stations for spatial queries
- Annotated status index on images for filtering

### File System

**Purpose:** Store uploaded image files

**Structure:**
```
uploads/
├── 1700000000_fog_image.jpg
├── 1700000001_ice_image.jpg
└── 1700000002_snow_image.jpg
```

**Naming Convention:**
- `{unix_timestamp}_{original_filename}`
- Prevents filename collisions
- Preserves original filename for reference

## Request Flow

### 1. Loading the Application

```
Browser → GET / → Backend → Serve index.html → Load Svelte App
                                              → GET /api/stations
                                              → GET /api/images
                                              → Load first unannotated image
```

### 2. Annotating an Image

```
User fills form → Submit → POST /api/annotations → Validate data
                                                  → Upsert annotation
                                                  → Mark image as annotated
                                                  → Return success
                         ← Response ← Update UI
                                   → Auto-load next image
```

### 3. Uploading Images

```
User selects files → Upload → POST /api/upload → Validate file
                                                → Generate filename
                                                → Save to disk
                                                → Insert to database
                              ← Response ← Update image list
```

### 4. Finding Nearest Station

```
User enters coordinates → Auto-trigger → GET /api/stations/nearest
                                       → Calculate distances
                                       → Return nearest station
                        ← Response ← Auto-select station
```

## Design Decisions

### Why Go?

- Fast compilation and execution
- Excellent concurrency support (though not heavily used here)
- Strong standard library for HTTP and database
- Simple deployment (single binary)
- Type safety and reliability

### Why Svelte?

- Reactive by default (simplifies state management)
- No virtual DOM (better performance)
- Smaller bundle size than React/Vue
- Less boilerplate code
- Excellent developer experience

### Why MySQL?

- Reliable and well-tested
- Good performance for this use case
- Wide hosting support
- Strong ACID guarantees
- Familiar SQL syntax

### Why REST API?

- Simple and well-understood
- Easy to test and debug
- No complex state management
- Good for CRUD operations
- Language-agnostic

### Local File Storage

**Pros:**
- Simple implementation
- No third-party dependencies
- Fast access
- No additional costs

**Cons:**
- Not horizontally scalable
- No built-in redundancy
- Manual backup required

**Future:** Could migrate to object storage (S3, MinIO) if needed.

## Security Considerations

### Current Implementation

1. **SQL Injection Prevention**
   - Prepared statements used throughout
   - No string concatenation for queries

2. **File Upload Security**
   - File size limits (32MB)
   - Filename sanitization
   - Separate storage from code

3. **CORS**
   - Configured for development
   - Should be restricted in production

4. **Input Validation**
   - Type checking on API parameters
   - Enum validation for categories and severities

### Production Recommendations

1. Add authentication (JWT, OAuth)
2. Add authorization (role-based access)
3. Rate limiting on API endpoints
4. HTTPS/TLS encryption
5. Content Security Policy headers
6. Regular security audits
7. Dependency updates
8. Database user with minimal privileges
9. File upload virus scanning
10. Request logging and monitoring

## Performance Considerations

### Current Optimizations

1. **Database Connection Pooling**
   - Max 10 open connections
   - 5 idle connections maintained

2. **Indexing**
   - All foreign keys indexed
   - Frequently queried fields indexed

3. **Frontend Build**
   - Vite production build
   - Code splitting
   - Asset optimization

### Scaling Strategies

**Vertical Scaling:**
- Increase server resources
- Optimize database queries
- Add database caching (Redis)

**Horizontal Scaling:**
- Load balancer (Nginx)
- Multiple backend instances
- Shared database
- Object storage for images
- CDN for image delivery

**Database Optimization:**
- Query optimization
- Index optimization
- Partitioning for large tables
- Read replicas for scaling reads

## Monitoring and Logging

### Current Approach

- Basic console logging
- Go's built-in HTTP logging
- MySQL query logs

### Production Recommendations

1. **Application Logging**
   - Structured logging (JSON)
   - Log levels (DEBUG, INFO, WARN, ERROR)
   - Log rotation
   - Centralized log aggregation

2. **Metrics**
   - Request count and latency
   - Error rates
   - Database query performance
   - System resources (CPU, memory, disk)

3. **Monitoring Tools**
   - Prometheus for metrics
   - Grafana for visualization
   - Alert Manager for alerts
   - Uptime monitoring

4. **Health Checks**
   - `/health` endpoint
   - Database connectivity check
   - Disk space check
   - Memory usage check

## Deployment Architecture

### Development
```
Local Machine → Go Server (8080) → Local MySQL (3306)
              → Vite Dev Server (5173)
```

### Production (Recommended)
```
Internet → Nginx (80/443) → Go Server (8080) → MySQL (3306)
        → SSL/TLS              ↓
        → Load Balancer        └→ File System
        → Static Assets
```

## Future Enhancements

### Short Term
1. Image thumbnails for faster loading
2. Pagination for large image lists
3. Search and filter functionality
4. Export annotations to CSV/JSON
5. Batch operations

### Medium Term
1. User authentication and authorization
2. Audit trail for annotations
3. Image comparison view
4. Statistics dashboard
5. API versioning

### Long Term
1. Machine learning integration (auto-detection)
2. Real-time collaboration
3. Mobile app
4. Advanced analytics
5. Integration with weather systems

## Testing Strategy

### Current Tests
- Unit tests for core functions (Haversine distance)
- Build verification

### Recommended Testing
1. **Unit Tests**
   - Business logic functions
   - Utility functions
   - Database operations (with test DB)

2. **Integration Tests**
   - API endpoint tests
   - Database integration
   - File upload/download

3. **E2E Tests**
   - User workflows
   - Browser automation (Playwright/Selenium)
   - Critical paths

4. **Performance Tests**
   - Load testing
   - Stress testing
   - Concurrent user testing

## Maintenance

### Regular Tasks
1. Database backups (daily)
2. Log rotation and cleanup
3. Dependency updates
4. Security patches
5. Disk space monitoring
6. Performance monitoring

### Tools
- Makefile for common tasks
- Backup scripts
- Monitoring dashboards
- Alert notifications

---

For more details on specific aspects:
- API details: [API.md](API.md)
- Deployment: [DEPLOYMENT.md](DEPLOYMENT.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)
