# API Documentation

## Overview

The IoT Device Telemetry API provides endpoints for managing users, devices, telemetry, and commands.
It follows a RESTful design, using JSON for request/response bodies.

- **Base URL**: `/api/v1`
- **Authentication**: JWT-based access tokens (short-lived) with refresh tokens (cookie-based).
- **Error Format**:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "error message"
  }
}
```

## Authentication

### Register User

**POST** `/auth/register`

Registers a new user.

**Request**:

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "StrongPassword123"
}
```

**Response** `201 Created`:

```json
{
  "id": "user-uuid",
  "username": "alice",
  "email": "alice@example.com"
}
```

### Login

**POST** `/auth/login`

Authenticates a user and returns an access token. Also sets a secure HTTP-only `refresh_token` cookie.

**Request**:

```json
{
  "email": "alice@example.com",
  "password": "StrongPassword123"
}
```

**Response** `200 OK`:

```json
{
  "access_token": "jwt-token",
  "expires_in": 3600
}
```

### Refresh Access Token

**POST** `/auth/refresh`

Refreshes the access token using the `refresh_token` cookie.

**Response** `200 OK`:

```json
{
  "access_token": "new-jwt-token",
  "expires_in": 3600
}
```

### Logout

**POST** `/auth/logout`

Revokes the refresh token and clears the cookie.

**Response** `204 No Content`

## Devices

### Create Device

**POST** `/devices`

**Request**:

```json
{
  "name": "Temperature Sensor",
  "device_type": "sensor",
  "status": "offline",
  "metadata": {
    "location": "greenhouse",
    "model": "esp32"
  }
}
```

**Response** `201 Created`:

```json
{
  "id": "device-uuid",
  "name": "Temperature Sensor",
  "device_type": "sensor",
  "status": "offline",
  "metadata": {
    "location": "greenhouse",
    "model": "esp32"
  }
}
```

### Get Device

**GET** `/devices/{device_id}`

**Response** `200 OK`:

```json
{
  "id": "device-uuid",
  "name": "Temperature Sensor",
  "device_type": "sensor",
  "status": "online",
  "metadata": {
    "location": "greenhouse",
    "model": "esp32"
  }
}
```

### List Devices

**GET** `/devices?limit=10&cursor=abcd123`

**Response** `200 OK`:

```json
[
  {
    "id": "device-uuid",
    "name": "Temperature Sensor",
    "device_type": "sensor",
    "status": "online",
    "metadata": {
      "location": "greenhouse",
      "model": "esp32"
    }
  }
]
```

**Pagination Meta (in headers or response meta)**:

```json
{
  "next_cursor": "abcd124",
  "limit": 10
}
```

### Update Device

**PATCH** `/devices/{device_id}`

**Request**:

```json
{
  "name": "Humidity Sensor",
  "device_type": "sensor",
  "metadata": {
    "location": "lab",
    "firmware": "1.2.0"
  }
}
```

**Response** `200 OK`:

```json
{
  "id": "device-uuid",
  "name": "Humidity Sensor",
  "device_type": "sensor",
  "status": "online",
  "metadata": {
    "location": "lab",
    "firmware": "1.2.0"
  }
}
```

## Telemetry

### Create Telemetry

**POST** `/devices/{device_id}/telemetry`

**Request**:

```json
{
  "telemetry_type": "environment",
  "payload": {
    "temperature": 22.5,
    "humidity": 60
  },
  "recorded_at": "2025-08-22T12:34:56Z"
}
```

**Response** `201 Created`:

```json
{
  "id": "telemetry-uuid",
  "telemetry_type": "environment",
  "payload": {
    "temperature": 22.5,
    "humidity": 60
  },
  "recorded_at": "2025-08-22T12:34:56Z"
}
```

### Get Device Telemetry

**GET** `/devices/{device_id}/telemetry?limit=10&cursor=abcd123`

**Response** `200 OK`:

```json
[
  {
    "id": "telemetry-uuid",
    "telemetry_type": "environment",
    "payload": {
      "temperature": 22.5,
      "humidity": 60
    },
    "recorded_at": "2025-08-22T12:34:56Z"
  }
]
```

**Pagination Meta**:

```json
{
  "next_cursor": "abcd124",
  "limit": 10
}
```

## Commands

### Create Command

**POST** `/devices/{device_id}/commands`

**Request**:

```json
{
  "command_name": "restart",
  "payload": { "delay": 5 }
}
```

**Response** `201 Created`:

```json
{
  "id": "command-uuid",
  "command_name": "restart",
  "payload": { "delay": 5 },
  "status": "pending"
}
```

### Get Device Commands

**GET** `/devices/{device_id}/commands?limit=10&cursor=abcd123`

**Response** `200 OK`:

```json
[
  {
    "id": "command-uuid",
    "command_name": "restart",
    "payload": { "delay": 5 },
    "status": "executed",
    "executed_at": "2025-08-22T12:40:00Z"
  }
]
```

**Pagination Meta**:

```json
{
  "next_cursor": "abcd124",
  "limit": 10
}
```

### Update Command Status

**PATCH** `/devices/{device_id}/commands/{command_id}`

**Request**:

```json
{
  "status": "executed",
  "executed_at": "2025-08-22T12:40:00Z"
}
```

**Response** `200 OK`:

```json
{
  "id": "command-uuid",
  "command_name": "restart",
  "payload": { "delay": 5 },
  "status": "executed",
  "executed_at": "2025-08-22T12:40:00Z"
}
```
