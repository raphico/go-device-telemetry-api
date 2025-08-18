# API Overview

## Authentication

All requests must be authenticated with a **Bearer Token**.
Include it in the `Authorization` header:

```http
Authorization: Bearer <access_token>
```

- Tokens are obtained via the **Auth API** (see Users & Auth section).
- If a request is missing or has an invalid token → return **401 Unauthorized**.

---

## Headers

Every request must include:

| Header            | Value                        | Required |
| ----------------- | ---------------------------- | -------- |
| `Authorization`   | `Bearer <token>`             | ✅       |
| `Content-Type`    | `application/json`           | ✅       |
| `Accept`          | `application/json`           | ✅       |
| `Idempotency-Key` | UUID for safe retries (POST) | Optional |

---

## Response Format

All responses follow a **standard envelope**:

```json
{
  "data": {...},
  "error": null
}
```

If an error occurs:

```json
{
  "data": null,
  "error": {
    "code": "DEVICE_NOT_FOUND",
    "message": "The requested device does not exist."
  }
}
```

- `code` → machine-readable error identifier
- `message` → human-readable explanation

---

## Errors

Common error codes across all endpoints:

| Code               | HTTP | Meaning                             |
| ------------------ | ---- | ----------------------------------- |
| `UNAUTHORIZED`     | 401  | Missing or invalid token            |
| `FORBIDDEN`        | 403  | Authenticated but not allowed       |
| `NOT_FOUND`        | 404  | Resource not found                  |
| `VALIDATION_ERROR` | 400  | Request payload invalid             |
| `CONFLICT`         | 409  | Resource already exists / duplicate |
| `INTERNAL_ERROR`   | 500  | Unexpected server error             |

---

## Pagination

Collection endpoints use cursor-based pagination (industry-standard at scale):

Request:

```http
GET /devices?limit=20&cursor=eyJpZCI6IjEyMyJ9
```

Response:

```json
{
  "data": [
    { "id": "dev_123", "name": "Temp Sensor" },
    { "id": "dev_124", "name": "Motor Controller" }
  ],
  "meta": {
    "next_cursor": "eyJpZCI6IjEyNCJ9",
    "has_more": true
  },
  "error": null
}
```

---

👉 This **Overview** will be referenced in **every section**.

Do you want me to move next into **Devices API** (create/list/get/disable), or would you prefer I flesh out **Users & Auth API** first (since Devices depends on having users)?
