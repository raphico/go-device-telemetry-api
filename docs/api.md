# API Documentation

## Overview

This API allows registering devices, sending commands, and retrieving telemetry data.
All requests and responses are in **JSON**.
Base URL:

```
/api/v1
```

## Authentication

All endpoints (except `/auth/register` and `/auth/login`) require authentication via Bearer token:

```
Authorization: Bearer <JWT_TOKEN>
```

## **Auth Routes**

### 1. Register User

`POST /auth/register`

Registers a new user.

**Request Body**:

```json
{
  "username": "john",
  "email": "john@example.com",
  "password": "strongpassword123"
}
```

**Response** `201 Created`:

```json
{
  "id": "user-uuid",
  "username": "john",
  "email": "john@example.com"
}
```

**Errors**:

- `400 Bad Request` → invalid email/password
- `409 Conflict` → email already exists

### 2. Login

`POST /auth/login`

Logs in a user, returns **JWT access_token** in JSON and sets a **refresh_token** in secure, HTTP-only cookie.

**Request Body**:

```json
{
  "email": "john@example.com",
  "password": "strongpassword123"
}
```

**Response** `200 OK`:

```json
{
  "access_token": "<JWT_TOKEN>",
  "expires_in": 3600
}
```

**Errors**:

- `401 Unauthorized` → invalid credentials

### 3. Refresh Token

`POST /auth/refresh`

Uses the **refresh token cookie** to issue a new access token.

**Response** `200 OK`:

```json
{
  "access_token": "<NEW_JWT_TOKEN>",
  "expires_in": 3600
}
```

**Errors**:

- `401 Unauthorized` → missing, invalid, or expired refresh token
