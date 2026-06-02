# Asynchronous Image Processing Flow

## Overview

The project uses an asynchronous image processing pipeline for uploaded post images.

Instead of processing images directly during the HTTP request lifecycle, the API uploads original images to object storage and creates background processing tasks through AWS SQS.

An AWS Lambda worker consumes processing events, generates thumbnails, validates image metadata, and stores processed metadata in PostgreSQL.

This architecture improves:

- API response time
- scalability
- retry safety
- fault tolerance
- frontend performance
- storage efficiency

---

# Architecture Flow

Client Upload  
↓  
POST /posts  
↓  
API validates image  
↓  
Original image uploaded to S3-compatible storage  
↓  
post_images row created in PostgreSQL  
↓  
processing_status = pending  
↓  
Database transaction committed  
↓  
SQS message published  
↓  
AWS Lambda triggered  
↓  
Image atomically claimed for processing  
↓  
Image downloaded from storage  
↓  
Image decoded and validated  
↓  
Thumbnail generated  
↓  
Thumbnail uploaded to storage  
↓  
Database metadata updated  
↓  
processing_status = completed

---

# Upload Flow

## 1. Image Upload

The client uploads images using multipart/form-data.

The API validates:

- file size
- MIME type
- allowed image formats

Supported formats:

- JPEG
- PNG
- WEBP
- GIF
- HEIC
- HEIF

---

## 2. Original Image Storage

Original images are uploaded to S3-compatible object storage.

The API stores only image metadata in PostgreSQL.

Example object key:

```text
customers/{customerID}/posts/{sessionID}/{fileID}.jpg
```
---

## 3. Database Insert

A record is created in guest.post_images.

Important fields:

| Field | Description |
|---|---|
| object_key | Original image storage key |
| processing_status | Current processing state |
| thumbnail_key | Generated thumbnail storage key |
| width | Original image width |
| height | Original image height |
| processed_at | Processing completion timestamp |
| failure_reason | Processing error details |

Default status: `pending`

Database constraints additionally enforce:

- valid processing states
- positive image dimensions
- metadata integrity

---

# Event Publishing

After the database transaction successfully commits, the API publishes an SQS message.

Example payload:

```json
{
  "image_id": "uuid",
  "s3_key": "customers/.../image.jpg"
}
```
This decouples HTTP requests from expensive image processing operations.

---

# Lambda Image Processing

## 1. SQS Consumption

AWS Lambda is triggered automatically by SQS.

The Lambda:

- parses the message
- validates payload
- processes messages independently
- reports partial batch failures back to SQS

This allows failed messages to retry without reprocessing successful messages in the same batch.

---

## 2. Atomic Processing Claim

To avoid duplicate image processing during concurrent Lambda executions or SQS retries, the worker atomically claims images for processing.

Only images in pending state can transition to processing.

This guarantees that only one worker processes a specific image.

Processing state transitions:

```text
pending
→ processing
→ completed
```
or:

```text 
pending 
→ processing 
→ failed
```

State transitions are enforced both in application logic and PostgreSQL constraints.

---

## 3. Image Validation

Lambda downloads the original image and validates:

- width
- height
- image integrity

Invalid or corrupted images are marked as failed.

---

## 4. Thumbnail Generation

The system generates resized thumbnails using Go image processing libraries.

Main libraries used:

- image
- image/jpeg
- image/png
- github.com/nfnt/resize

Generated thumbnails are uploaded separately.

Example thumbnail key:

```text
posts/thumbnails/{sessionID}/{fileID}.jpg
```

---

## 5. Metadata Update

After successful processing, Lambda updates:

- processing_status
- thumbnail_key
- width
- height
- processed_at

Example completed state:

```text
processing_status = completed
```
---

# Failure Handling

If processing fails:
`processing_status = failed`
and:
`failure_reason`
is stored for debugging and observability.

The original image remains stored even if thumbnail generation fails.

---

# Retry Safety

The pipeline is designed to be retry-safe.

Key mechanisms:

- SQS retries failed messages automatically
- Lambda uses atomic processing claims
- duplicate processing is prevented
- partial batch failure responses avoid retrying successful messages

This makes the system resilient to:

- temporary infrastructure failures
- Lambda crashes
- network interruptions
- duplicate queue deliveries

---

# API Metadata Exposure

API responses expose processed image metadata.

Example:

```json
{
  "images": [
    {
      "objectKey": "customers/.../image.jpg",
      "url": "presigned-original-url",
      "thumbnailKey": "posts/thumbnails/.../thumb.jpg",
      "thumbnailURL": "public-thumbnail-url",
      "processingStatus": "completed",
      "width": 1200,
      "height": 800
    }
  ]
}
```
Frontend applications can use this metadata to:

- display loading states
- show thumbnails immediately
- react to processing status
- optimize image rendering

---

# Storage Access Strategy

## Original Images

Original uploaded images are private.

The API generates presigned URLs for temporary secure access.

Benefits:

- secure access control
- temporary authorization
- no direct public exposure

---

## Thumbnails

Generated thumbnails are public and returned through direct storage URLs.

Benefits:

- faster frontend loading
- CDN compatibility
- improved caching
- lower bandwidth usage

This architecture allows large originals to remain protected while lightweight thumbnails are optimized for public delivery.

---

# Benefits of the Architecture

## Performance

Heavy image processing is removed from the synchronous HTTP request lifecycle.

---

## Scalability

SQS and Lambda allow independent horizontal scaling of image processing workers.

---

## Reliability

Queue-based retries improve processing reliability.

---

## Fault Tolerance

Failures are isolated from user-facing API requests.

---

## Better Frontend UX

Frontend applications can:

- display thumbnails quickly
- show processing indicators
- progressively load images
- reduce mobile bandwidth usage

---

# Main Components

| Component | Responsibility |
|---|---|
| API Service | Upload images and publish processing events |
| PostgreSQL | Store image metadata and processing state |
| PostgreSQL Constraints | Enforce metadata integrity |
| S3-compatible storage | Store original images and thumbnails |
| AWS SQS | Queue asynchronous image processing tasks |
| AWS Lambda | Execute background image processing |
| Image Processing Service | Decode, validate, resize, and upload thumbnails |

---

# Key Backend Concepts Used

- asynchronous processing
- event-driven architecture
- queue-based workers
- atomic work claiming
- retry-safe processing
- object storage architecture
- presigned URLs
- image optimization
- metadata-driven APIs
- serverless processing
- partial batch failure handling
- distributed worker coordination