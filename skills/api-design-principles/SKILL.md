---
name: api-design-principles
description: Guide for designing RESTful APIs following industry best practices. Use when designing new APIs, refactoring existing APIs, reviewing API specifications, or making architectural decisions about API structure. Covers REST principles, HTTP semantics, resource naming, error handling, versioning, security, and documentation.
---

# API Design Principles

## Core REST Principles

### Resource-Oriented Design

APIs should be designed around resources (nouns) not actions (verbs).

**Good:**
```
GET    /users          # List users
GET    /users/123      # Get specific user
POST   /users          # Create user
PUT    /users/123      # Update user
DELETE /users/123      # Delete user
```

**Avoid:**
```
GET    /getUsers
POST   /createUser
POST   /users/delete  # Use DELETE instead
```

### HTTP Method Semantics

Use HTTP methods correctly according to RFC 7231:

| Method | Safe | Idempotent | Purpose |
|--------|------|------------|---------|
| GET    | ✅   | ✅         | Retrieve resource representation |
| POST   | ❌   | ❌         | Create resource or trigger action |
| PUT    | ❌   | ✅         | Replace resource (full update) |
| PATCH  | ❌   | ❌         | Modify resource (partial update) |
| DELETE | ❌   | ✅         | Delete resource |

**Key rules:**
- GET must never modify server state
- POST for non-idempotent operations
- PUT replaces entire resource, PATCH modifies partially
- DELETE removes resource (can be soft-delete)

### URL Structure

**Best practices:**

1. **Use nouns, not verbs**
   ```
   ✅ /users
   ❌ /getUsers
   ```

2. **Use plural nouns for collections**
   ```
   ✅ /users/123/orders
   ❌ /user/123/order
   ```

3. **Use kebab-case for readability**
   ```
   ✅ /shipping-addresses
   ❌ /shippingAddresses
   ```

4. **Keep URLs shallow (max 3-4 levels)**
   ```
   ✅ /users/123/orders/456/items
   ❌ /orgs/depts/users/123/orders/items/456/details
   ```

5. **Avoid query parameters for resources**
   ```
   ✅ /users/123/orders
   ❌ /users?type=active&include=orders
   ```

### Resource Naming Conventions

```
/users                 # Collection
/users/123             # Specific resource
/users/123/orders      # Sub-collection
/users/123/orders/456  # Specific sub-resource
```

**Special cases:**
```
/archived-users        # Use hyphen for adjectives
/shipping-addresses    # Plural for multi-word nouns
/api/v1/users          # Version prefix
```

## Request/Response Design

### Request Body Guidelines

**POST (Create):**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "role": "admin"
}
```

**PUT (Full Update):**
```json
{
  "id": "123",
  "name": "John Doe",
  "email": "john.doe@example.com",  // Changed
  "role": "user"                     // Changed
}
```

**PATCH (Partial Update):**
```json
{
  "email": "john.doe@example.com"   // Only changed field
}
```

**Best practices:**
- PUT requires all fields (or use merge semantics)
- PATCH requires only changed fields
- Document merge behavior clearly
- Use content-type: application/merge-patch+json if needed

### Response Format

**Success responses:**
```json
{
  "data": {
    "id": "123",
    "name": "John Doe"
  },
  "meta": {
    "timestamp": "2026-01-09T10:00:00Z",
    "version": "1.0"
  }
}
```

**Collection responses:**
```json
{
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 20
  }
}
```

**Error responses:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": [
      {
        "field": "email",
        "reason": "Must be valid email address"
      }
    ]
  }
}
```

### Pagination

**Cursor-based (recommended for large datasets):**
```
GET /users?limit=20&cursor=abc123
```

Response:
```json
{
  "data": [...],
  "meta": {
    "next_cursor": "def456",
    "has_more": true
  }
}
```

**Offset-based (simpler, but less efficient):**
```
GET /users?page=1&per_page=20
```

Response:
```json
{
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 20,
    "total_pages": 5
  }
}
```

## Status Codes

Use HTTP status codes correctly:

### Success Codes
- `200 OK` - Successful GET, PUT, PATCH
- `201 Created` - Successful POST, include Location header
- `204 No Content` - Successful DELETE, PUT with no response body

### Client Error Codes (4xx)
- `400 Bad Request` - Invalid request body, parameters
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Valid auth but insufficient permissions
- `404 Not Found` - Resource doesn't exist
- `409 Conflict` - Resource state conflict (e.g., duplicate email)
- `422 Unprocessable Entity` - Semantic errors (validation failures)
- `429 Too Many Requests` - Rate limit exceeded

### Server Error Codes (5xx)
- `500 Internal Server Error` - Unexpected server error
- `503 Service Unavailable` - Service temporarily down

**Best practice:** Always return error details in response body for 4xx errors.

## Versioning

### Strategies

**URL Path Versioning (recommended):**
```
/api/v1/users
/api/v2/users
```

Pros:
- Clear version separation
- Easy to deprecate old versions
- CDN-friendly

Cons:
- Version in URL (some consider this REST-violating)

**Header Versioning:**
```
GET /users
Accept: application/vnd.myapi.v2+json
```

Pros:
- Clean URLs
- Backwards compatible

Cons:
- Harder to debug/test
- Not cache-friendly at URL level

**Query Parameter Versioning (not recommended):**
```
/users?version=2
```

Cons:
- Easy to miss
- Not cache-friendly

### Versioning Best Practices

1. **Start without version** for v1
   ```
   /users  # Assume v1
   ```

2. **Use semantic versioning**
   - MAJOR version for breaking changes
   - MINOR for additions (backwards compatible)
   - PATCH for bug fixes

3. **Deprecation policy**
   - Support at least 2 versions simultaneously
   - Communicate deprecation in response headers
   ```
   Deprecation: true
   Sunset: 2026-06-01
   ```

4. **Document breaking changes clearly**

## Filtering, Sorting, and Search

### Filtering
```
GET /users?status=active&role=admin
```

Implement:
- Exact match: `?status=active`
- Multiple values: `?status=active,pending`
- Ranges: `?age[gte]=18&age[lte]=65`
- Booleans: `?verified=true`

### Sorting
```
GET /users?sort=created_at:desc,name:asc
```

Format: `sort=field:direction`
- Default direction: `asc`
- Multiple fields supported
- Prefix with `-` for desc: `?sort=-created_at`

### Searching
```
GET /users?q=john
```

Implement:
- Full-text search
- Partial matching on name/email
- Consider advanced search with POST body for complex queries

## Security Best Practices

### Authentication

**Use standard headers:**
```
Authorization: Bearer <token>
Authorization: ApiKey <key>
```

**Never send credentials in URL:**
```
❌ GET /users?token=abc123
❌ GET /users/abc123  # token in path
```

### Authorization

- Check permissions on every request
- Return 403 (not 404) for authorization failures
- Document required permissions per endpoint
- Use principle of least privilege

### Rate Limiting

Return rate limit info in headers:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1641753600
```

On limit exceeded:
```
429 Too Many Requests
X-RateLimit-Retry-After: 60
```

### Input Validation

Validate:
1. **Request body schema**
   - Required fields
   - Data types
   - String formats (email, UUID, etc.)
   - Value ranges

2. **Query parameters**
   - Type checking
   - Range validation
   - Allowed values (enums)

3. **Return detailed validation errors**
   ```json
   {
     "error": {
       "code": "VALIDATION_ERROR",
       "message": "Validation failed",
       "details": [
         {
           "field": "email",
           "message": "Invalid email format"
         }
       ]
     }
   }
   ```

## Documentation

### OpenAPI/Swagger Specification

Always maintain OpenAPI spec (YAML or JSON):

```yaml
openapi: 3.0.0
info:
  title: User Management API
  version: 1.0.0
paths:
  /users:
    get:
      summary: List users
      parameters:
        - name: page
          in: query
          schema:
            type: integer
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserList'
```

### Documentation Best Practices

1. **Provide clear examples**
   - Request examples
   - Response examples (success and error)
   - Use real-world data

2. **Document all fields**
   - Field types
   - Required vs optional
   - Constraints (length, range, format)
   - Default values

3. **Include error catalog**
   - All possible error codes
   - Error message format
   - Troubleshooting steps

4. **Keep docs in sync with code**
   - Use code generation from OpenAPI spec
   - Version docs with API
   - Provide changelog

## Performance Considerations

### Minimize Response Size

- Use field selection: `GET /users?fields=id,name`
- Compress responses: `Accept-Encoding: gzip`
- Paginate large collections
- Omit null fields (optional)

### Caching

Use HTTP caching headers:
```
Cache-Control: max-age=3600, public
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
Last-Modified: Wed, 09 Jan 2026 10:00:00 GMT
```

Support conditional requests:
```
If-None-Match: "33a64df551425fcc55e4d42a148795d9f25f89d4"
If-Modified-Since: Wed, 09 Jan 2026 10:00:00 GMT
```

Return `304 Not Modified` if not changed.

### Use HATEOAS (Hypermedia)

Include links in responses:
```json
{
  "data": {
    "id": "123",
    "name": "John Doe"
  },
  "links": {
    "self": "/users/123",
    "orders": "/users/123/orders",
    "avatar": "/users/123/avatar"
  }
}
```

## Common Pitfalls to Avoid

### 1. Inconsistent Naming
```
❌ /getUser
❌ /list-users
❌ /User
✅ /users
```

### 2. Wrong HTTP Methods
```
❌ GET /users/123/delete
❌ POST /users/123/delete
✅ DELETE /users/123
```

### 3. Returning 200 for Errors
```
❌ 200 OK with {"error": "Not found"}
✅ 404 Not Found with error details
```

### 4. Nested Resources Too Deep
```
❌ /orgs/depts/teams/users/123/posts/456/comments/789
✅ /comments/789 (use IDs for deep resources)
```

### 5. Missing Metadata
```
❌ {"data": [...]}  // No pagination info
✅ {"data": [...], "meta": {"total": 100, "page": 1}}
```

## Checklist for API Design

Review your API design against this checklist:

### Structure
- [ ] Uses nouns, not verbs in URLs
- [ ] Plural nouns for collections
- [ ] Kebab-case for multi-word names
- [ ] Consistent URL structure (max 3-4 levels)
- [ ] Clear resource hierarchy

### HTTP Methods
- [ ] GET for retrieval (never modifies state)
- [ ] POST for creation or non-idempotent actions
- [ ] PUT for full replacement
- [ ] PATCH for partial updates
- [ ] DELETE for deletion

### Status Codes
- [ ] 200 OK for successful GET/PUT/PATCH
- [ ] 201 Created for POST with Location header
- [ ] 204 No Content for DELETE
- [ ] Appropriate 4xx codes for client errors
- [ ] Appropriate 5xx codes for server errors

### Request/Response
- [ ] Clear request body schema
- [ ] Consistent response format
- [ ] Detailed error responses
- [ ] Pagination for collections
- [ ] Filtering and sorting support

### Security
- [ ] Authentication required
- [ ] Authorization checked
- [ ] Input validation
- [ ] Rate limiting
- [ ] HTTPS required

### Documentation
- [ ] OpenAPI/Swagger spec maintained
- [ ] Request/response examples
- [ ] Error documentation
- [ ] Changelog maintained
- [ ] Versioning strategy documented

---

**Last updated**: 2026-01-09
**Version**: 1.0
**Maintained by**: API Design Team
