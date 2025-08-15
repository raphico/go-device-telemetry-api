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
  "username": "john0x0",
  "email": "john@example.com",
  "password": "strongpassword123"
}
```

**Response** `201 Created`:

```json
{
  "id": "user-uuid",
  "name": "John Snow",
  "email": "john@example.com"
}
```
