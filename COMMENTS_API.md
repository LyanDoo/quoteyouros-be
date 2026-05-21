# Comments API Documentation

## Overview

The Comments API provides a complete system for managing blog post comments with support for:
- **Top-level comments** on blog posts
- **Nested replies** to existing comments
- **Star ratings** (1-5 scale, optional)
- **Anonymous and authenticated** comments
- **Spam detection** and moderation (admin-only)

---

## Base URL

```
http://localhost:8030/api
```

---

## Authentication

**Protected endpoints** require a JWT bearer token:

```
Authorization: Bearer <jwt_token>
```

Obtain a token by logging in or registering through the Auth endpoints.

---

## Endpoints

### 1. Get Blog Post Comments

Retrieve all top-level comments for a blog post (public).

```
GET /blog/{blog_id}/comments?page=1&limit=10
```

**Parameters:**
- `page` (query, optional): Page number (default: 1)
- `limit` (query, optional): Items per page, max 100 (default: 10)

**Response (200 OK):**
```json
{
  "data": {
    "comments": [
      {
        "id": "uuid-string",
        "blog_post_id": "uuid-string",
        "reply_to_comment_id": null,
        "author_name": "John Doe",
        "author_email": "john@example.com",
        "content": "Great blog post!",
        "rating": 5,
        "is_spam": false,
        "created_at": "2026-05-21 14:30:45",
        "updated_at": "2026-05-21 14:30:45"
      }
    ],
    "total": 42,
    "page": 1,
    "limit": 10,
    "total_pages": 5
  },
  "message": "Successfully retrieved comments"
}
```

**Error Response (400/404/500):**
```json
{
  "error": "Blog post not found",
  "status": 404
}
```

---

### 2. Get Comment Replies

Retrieve all replies to a specific comment (public).

```
GET /comments/{comment_id}/replies?page=1&limit=10
```

**Parameters:**
- `comment_id` (path, required): ID of parent comment
- `page` (query, optional): Page number (default: 1)
- `limit` (query, optional): Items per page, max 100 (default: 10)

**Response (200 OK):**
```json
{
  "data": {
    "comments": [
      {
        "id": "uuid-string",
        "blog_post_id": "uuid-string",
        "reply_to_comment_id": "parent-comment-uuid",
        "author_name": "Jane Smith",
        "author_email": null,
        "content": "I completely agree!",
        "rating": null,
        "is_spam": false,
        "created_at": "2026-05-21 15:45:20",
        "updated_at": "2026-05-21 15:45:20"
      }
    ],
    "total": 3,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  },
  "message": "Successfully retrieved replies"
}
```

---

### 3. Create Comment (Top-level)

Post a new comment on a blog post (public).

```
POST /blog/{blog_id}/comments
Content-Type: application/json
```

**Request Body:**
```json
{
  "author_name": "John Doe",
  "author_email": "john@example.com",
  "content": "This is a thoughtful comment about the blog post.",
  "rating": 5
}
```

**Field Requirements:**
- `author_name` (required, string): Commenter's name
- `author_email` (optional, string): Commenter's email (omit for anonymous)
- `content` (required, string): Comment text
- `rating` (optional, integer): Star rating 1-5
- `reply_to_comment_id` (omitted for top-level comments)

**Response (201 Created):**
```json
{
  "data": {
    "id": "new-comment-uuid",
    "blog_post_id": "blog-uuid",
    "reply_to_comment_id": null,
    "author_name": "John Doe",
    "author_email": "john@example.com",
    "content": "This is a thoughtful comment about the blog post.",
    "rating": 5,
    "is_spam": false,
    "created_at": "2026-05-21 16:00:00",
    "updated_at": "2026-05-21 16:00:00"
  },
  "message": "Comment created successfully"
}
```

---

### 4. Reply to Comment

Post a reply to an existing comment (public).

```
POST /blog/{blog_id}/comments
Content-Type: application/json
```

**Request Body:**
```json
{
  "author_name": "Jane Smith",
  "author_email": null,
  "content": "I completely agree with your point!",
  "rating": null,
  "reply_to_comment_id": "parent-comment-uuid"
}
```

**Field Requirements:**
- `author_name` (required, string): Replier's name
- `author_email` (optional, string): Replier's email (omit for anonymous)
- `content` (required, string): Reply text
- `rating` (optional, integer): Star rating 1-5
- `reply_to_comment_id` (required, string): UUID of parent comment

**Response (201 Created):**
```json
{
  "data": {
    "id": "new-reply-uuid",
    "blog_post_id": "blog-uuid",
    "reply_to_comment_id": "parent-comment-uuid",
    "author_name": "Jane Smith",
    "author_email": null,
    "content": "I completely agree with your point!",
    "rating": null,
    "is_spam": false,
    "created_at": "2026-05-21 16:05:00",
    "updated_at": "2026-05-21 16:05:00"
  },
  "message": "Comment created successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Parent comment not found",
  "status": 400
}
```

---

### 5. Update Comment (Admin Only)

Modify a comment's rating, content, or mark as spam (protected).

```
PUT /comments/{comment_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "content": "Updated comment text",
  "rating": 4,
  "is_spam": false
}
```

**Field Requirements:**
- `content` (optional, string): Updated comment text
- `rating` (optional, integer): Updated rating 1-5
- `is_spam` (optional, boolean): Mark/unmark as spam

**Response (200 OK):**
```json
{
  "data": {
    "id": "comment-uuid",
    "blog_post_id": "blog-uuid",
    "reply_to_comment_id": null,
    "author_name": "John Doe",
    "author_email": "john@example.com",
    "content": "Updated comment text",
    "rating": 4,
    "is_spam": false,
    "created_at": "2026-05-21 14:30:45",
    "updated_at": "2026-05-21 16:15:00"
  },
  "message": "Comment updated successfully"
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "Missing or invalid authorization token",
  "status": 401
}
```

---

### 6. Delete Comment (Admin Only)

Remove a comment and all its replies (protected).

```
DELETE /comments/{comment_id}
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**
```json
{
  "message": "Comment deleted successfully"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Comment not found",
  "status": 404
}
```

---

## Usage Examples

### JavaScript/TypeScript (Fetch API)

**Create a Top-level Comment:**
```javascript
const response = await fetch('http://localhost:8030/api/blog/blog-id-123/comments', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    author_name: 'Alice Johnson',
    author_email: 'alice@example.com',
    content: 'Excellent insights in this post!',
    rating: 5
  })
});

const result = await response.json();
console.log(result.data.id); // New comment UUID
```

**Create an Anonymous Reply:**
```javascript
const response = await fetch('http://localhost:8030/api/blog/blog-id-123/comments', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    author_name: 'Anonymous User',
    content: 'Thanks for sharing this!',
    reply_to_comment_id: 'parent-comment-uuid'
  })
});
```

**Get Comments with Pagination:**
```javascript
const response = await fetch('http://localhost:8030/api/blog/blog-id-123/comments?page=1&limit=20');
const result = await response.json();
console.log(result.data.comments);
console.log(`Showing ${result.data.page} of ${result.data.total_pages} pages`);
```

**Get Replies to a Comment:**
```javascript
const response = await fetch('http://localhost:8030/api/comments/comment-uuid/replies?page=1&limit=10');
const result = await response.json();
console.log(result.data.comments); // Array of replies
```

**Update Comment (with Auth):**
```javascript
const token = localStorage.getItem('jwt_token');
const response = await fetch('http://localhost:8030/api/comments/comment-uuid', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    is_spam: true // Mark as spam
  })
});
```

**Delete Comment (with Auth):**
```javascript
const token = localStorage.getItem('jwt_token');
const response = await fetch('http://localhost:8030/api/comments/comment-uuid', {
  method: 'DELETE',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

---

## Data Model

### Comment Entity

```typescript
interface Comment {
  id: string;                      // UUID, unique identifier
  blog_post_id: string;            // UUID, associated blog post
  reply_to_comment_id: string | null; // UUID or null (null = top-level comment)
  author_name: string;             // Required, commenter's name
  author_email: string | null;     // Optional, commenter's email
  content: string;                 // Required, comment text
  rating: number | null;           // Optional, 1-5 star rating
  is_spam: boolean;                // Default: false, admin-marked spam
  created_at: string;              // Timestamp "YYYY-MM-DD HH:mm:ss"
  updated_at: string;              // Timestamp "YYYY-MM-DD HH:mm:ss"
}
```

### Comment Query Response

```typescript
interface CommentListResponse {
  comments: Comment[];
  total: number;              // Total comments matching filter
  page: number;               // Current page number
  limit: number;              // Items per page
  total_pages: number;        // Total pages available
}
```

---

## Pagination

All list endpoints support pagination:

- `page`: Current page (minimum 1)
- `limit`: Items per page (1-100, default 10)

**Example with defaults (page 1, limit 10):**
```
GET /blog/blog-id/comments
```

**Example with custom pagination:**
```
GET /blog/blog-id/comments?page=3&limit=25
```

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "Error message describing the issue",
  "status": 400
}
```

**Common HTTP Status Codes:**

| Status | Meaning |
|--------|---------|
| 200 | Success (OK) |
| 201 | Success (Created) |
| 400 | Bad Request (invalid data or reference) |
| 401 | Unauthorized (missing/invalid JWT token) |
| 404 | Not Found (comment/post doesn't exist) |
| 500 | Server Error |

---

## Features

✅ **Public Access**: All users can view and post comments/replies without authentication
✅ **Anonymous Posting**: Comments can be posted without email (username-only)
✅ **Nested Replies**: Reply to specific comments, no depth limit
✅ **Star Ratings**: Optional 1-5 star ratings for comments
✅ **Admin Moderation**: Admins can edit content, ratings, and mark spam
✅ **Cascade Delete**: Deleting a comment automatically removes all replies
✅ **Spam Filtering**: Spam-marked comments hidden from public views
✅ **Pagination**: Efficient large-scale comment retrieval
✅ **Timestamps**: All comments include creation and update times

---

## Integration Checklist

- [ ] Create comment form component
- [ ] Display top-level comments with pagination
- [ ] Implement nested reply display (threaded view)
- [ ] Add star rating widget (1-5 selector)
- [ ] Create reply form (with parent comment context)
- [ ] Handle anonymous posting (optional email field)
- [ ] Display timestamps (format: "May 21, 2026 2:30 PM")
- [ ] Implement spam filter UI (hide spam by default)
- [ ] Add admin moderation panel (edit/delete/mark spam)
- [ ] Error handling and validation feedback
- [ ] Loading states and optimistic updates
- [ ] Empty state messaging

---

## Testing with Postman

Use the included `postman.json` collection which includes:
- **Get Blog Post Comments**
- **Get Comment Replies**
- **Create Comment**
- **Reply to Comment**
- **Update Comment**
- **Delete Comment**

All endpoints are pre-configured with variables for easy testing.

---

## Support

For backend API issues, contact the development team or check the server logs at:
```
docker logs quoteyouros-backend
```
