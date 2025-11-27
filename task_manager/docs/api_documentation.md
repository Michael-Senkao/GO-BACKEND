# Task Manager API Documentation

### GET `/tasks`

Retrieve all tasks.

### GET `/tasks/:id`

Retrieve a task by ID.

### POST `/tasks`

Create a new task.

**Payload example:**

```json
{
  "id": "4",
  "title": "New Task",
  "description": "Some description",
  "due_date": "2025-01-01T00:00:00Z",
  "status": "Pending"
}
```
