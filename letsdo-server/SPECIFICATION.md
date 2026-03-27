# SPECIFICATION.md

## Project Overview
**Name:** letsdo
**Type:** API
**One-liner:** Simple todo app LetsDo - minimal backend structure with auth

## Domain Terms
| Term | Definition |
|------|------------|
| Todo | The primary entity in this system - todo that needs to be managed |
| Todo Status | The state of a todo (active, completed, archived, etc.) |
{{if .IncludeAuth}}| User | Authenticated user who can manage todos |
| Role | Permission level for users (admin, user, guest) |{{end}}

## Features

### Core (MVP)
- [ ] **List Todos**: View all todos with pagination and filtering
- [ ] **Create Todo**: Add new todo with validation
- [ ] **Update Todo**: Modify existing todo properties
- [ ] **Delete Todo**: Remove todo from the system
- [ ] **Get Todo Details**: View single todo with all information
{{if .IncludeAuth}}- [ ] **User Authentication**: Register, login, and manage user sessions
- [ ] **Authorization**: Role-based access control for todo operations{{end}}
{{if .IncludeS3}}- [ ] **File Upload**: Upload and manage files/images for todos{{end}}

### Future Enhancements
- [ ] **Bulk Operations**: Create/update/delete multiple todos at once
- [ ] **Export Data**: Download todos in CSV/JSON format
- [ ] **Webhooks**: Notify external systems of todo events
- [ ] **Audit Log**: Track all changes to todos
- [ ] **Search**: Full-text search across todo fields

## Data Model

### Todo
- `id`: UUID — unique identifier
- `name`: string — display name/title
- `description`: text — detailed description
- `status`: enum — current state (active, completed, archived)
- `created_at`: timestamp — when created
- `updated_at`: timestamp — last modified
{{if .IncludeAuth}}- `user_id`: UUID — owner/creator of the todo{{end}}
{{if .IncludeS3}}- `image_url`: string — uploaded image URL{{end}}
- **Relations:** {{if .IncludeAuth}}belongs to User{{end}}

{{if .IncludeAuth}}### User
- `id`: UUID — unique identifier
- `email`: string — primary identifier for login
- `password`: string (hashed) — bcrypt hashed password
- `role`: enum — permission level (admin, user)
- `created_at`: timestamp
- **Relations:** has many Todos{{end}}

## API Surface

### Todo Endpoints
| Method | Endpoint | Purpose | Auth |
|--------|----------|---------|------|
| GET | /api/todos | List all todos (paginated) | {{if .IncludeAuth}}Optional{{else}}None{{end}} |
| GET | /api/todos/:id | Get single todo | {{if .IncludeAuth}}Optional{{else}}None{{end}} |
| POST | /api/todos | Create new todo | {{if .IncludeAuth}}Required{{else}}None{{end}} |
| PUT | /api/todos/:id | Update todo | {{if .IncludeAuth}}Required{{else}}None{{end}} |
| DELETE | /api/todos/:id | Delete todo | {{if .IncludeAuth}}Required{{else}}None{{end}} |

{{if .IncludeAuth}}### Auth Endpoints
| Method | Endpoint | Purpose | Auth |
|--------|----------|---------|------|
| POST | /api/auth/register | Create new user account | None |
| POST | /api/auth/login | Authenticate and get JWT token | None |
| POST | /api/auth/refresh | Refresh JWT token | Required |
| GET | /api/auth/me | Get current user info | Required |{{end}}

{{if .IncludeS3}}### File Upload Endpoints
| Method | Endpoint | Purpose | Auth |
|--------|----------|---------|------|
| POST | /api/upload | Upload file to S3 | {{if .IncludeAuth}}Required{{else}}None{{end}} |
| GET | /api/files/:key | Get pre-signed URL for file | {{if .IncludeAuth}}Required{{else}}None{{end}} |{{end}}

### Request/Response Examples

#### Create Todo
**Request:**
```json
POST /api/todos
{
  "name": "Example Todo",
  "description": "This is an example todo",
  "status": "active"
}
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Example Todo",
  "description": "This is an example todo",
  "status": "active",
  "created_at": "2024-01-14T10:00:00Z",
  "updated_at": "2024-01-14T10:00:00Z"
}
```

## External Integrations

| Service | Purpose | Auth |
|---------|---------|------|
| PostgreSQL | Primary data storage | Connection string |
| Redis | Caching and session storage | Connection string |
{{if .IncludeS3}}| AWS S3 | File/image storage | API Key + Secret |{{end}}
{{if .IncludeAuth}}| JWT | Token-based authentication | Secret key |{{end}}

## Business Rules

### Todo Management
1. **Unique Names**: Todo names must be unique within the system
2. **Status Transitions**: Todos can only transition between valid states
3. **Validation**: All required fields must be provided when creating todo
{{if .IncludeAuth}}4. **Ownership**: Users can only modify their own todos (unless admin)
5. **Admin Override**: Admin users can manage all todos{{end}}

{{if .IncludeAuth}}### Authentication & Authorization
1. **Password Requirements**: Minimum 8 characters
2. **Token Expiry**: JWT tokens expire after 24 hours
3. **Refresh Tokens**: Can be used to get new access tokens
4. **Role Hierarchy**: Admin > User > Guest (if applicable){{end}}

{{if .IncludeS3}}### File Uploads
1. **File Size Limit**: Maximum 10MB per file
2. **Allowed Types**: Images only (jpg, png, gif, webp)
3. **Storage**: Files stored in S3 with unique keys
4. **Access**: Pre-signed URLs for secure access{{end}}

## Edge Cases

### Data Consistency
- **Concurrent Updates**: Use database transactions to prevent race conditions
- **Cascade Deletes**: Handle related data when deleting todo
{{if .IncludeAuth}}- **User Deletion**: Decide how to handle todos when user is deleted{{end}}

### Error Handling
- **Database Connection Loss**: Gracefully handle and retry connections
- **Invalid Input**: Return clear validation error messages
{{if .IncludeS3}}- **S3 Upload Failures**: Provide fallback or retry mechanism{{end}}
{{if .IncludeAuth}}- **Token Expiry During Request**: Handle gracefully with proper error codes{{end}}

### Performance
- **Large Dataset Queries**: Implement proper pagination
- **Cache Invalidation**: Clear cache when data changes
- **Database Indexing**: Add indexes for frequently queried fields

## Non-Functional Requirements

### Performance
- API response time < 200ms for simple queries
- Support 100+ concurrent users
- Database query optimization with proper indexes

### Security
- Input validation on all endpoints
- SQL injection prevention
{{if .IncludeAuth}}- Secure password storage with bcrypt
- JWT token security best practices{{end}}
{{if .IncludeS3}}- Secure file upload validation
- S3 bucket security policies{{end}}

### Scalability
- Horizontal scaling capability
- Database connection pooling
- Redis for caching frequently accessed data

### Monitoring
- Health check endpoint
- Structured logging
- Error tracking and alerting

## Testing Requirements

### Unit Tests
- Repository layer methods
- Service layer business logic
- Input validation functions
{{if .IncludeAuth}}- Authentication middleware{{end}}

### Integration Tests
- Full API endpoint flows
- Database operations
{{if .IncludeS3}}- File upload/download flows{{end}}
{{if .IncludeAuth}}- Authentication and authorization flows{{end}}

### Performance Tests
- Load testing for concurrent users
- Database query performance
- Cache hit/miss ratios

---

*This specification defines WHAT we're building. For HOW to implement it, see CLAUDE.md. For testing support, use AGENTS.md with the /test command.*