# Industrix Todo App

Aplikasi todo list full-stack dengan Go (Gin + GORM) backend dan React (TypeScript + Ant Design) frontend.

## Fitur

- CRUD todos (create, read, update, delete)
- Mark complete/incomplete
- Kategori dengan warna custom
- Pagination, search, dan filter
- Responsive design (desktop, tablet, mobile)
- Backend unit tests
- TypeScript frontend

## Tech Stack

**Backend:** Go 1.21, Gin, GORM, PostgreSQL
**Frontend:** React 18, TypeScript, Ant Design 5, Vite, Context API

## Cara Menjalankan

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### 1. Clone & Setup Database
```bash
git clone https://github.com/nathwiels26/industrix-todo-app
cd industrix-todo-app

# Buat database
psql -U postgres
CREATE DATABASE industrix_todo;
\q
```

### 2. Jalankan Backend
```bash
cd backend

# Buat file .env dengan isi:
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=your_password
# DB_NAME=industrix_todo
# SERVER_PORT=8080

go mod download
go run cmd/server/main.go
```
Backend jalan di `http://localhost:8080`

### 3. Jalankan Frontend
```bash
cd frontend
npm install
npm run dev
```
Frontend jalan di `http://localhost:5173`

### 4. Jalankan Tests
```bash
cd backend
go test ./tests/... -v
```

## API Endpoints

### Todos
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | /api/todos | List todos (dengan pagination, search, filter) |
| POST | /api/todos | Buat todo baru |
| GET | /api/todos/:id | Get todo by ID |
| PUT | /api/todos/:id | Update todo |
| DELETE | /api/todos/:id | Hapus todo |
| PATCH | /api/todos/:id/complete | Toggle status complete |

**Query params untuk GET /api/todos:**
- `page`, `limit` - pagination
- `search` - cari by title
- `category_id`, `completed`, `priority` - filter
- `sort_by`, `sort_order` - sorting

### Categories
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | /api/categories | List semua kategori |
| POST | /api/categories | Buat kategori baru |
| PUT | /api/categories/:id | Update kategori |
| DELETE | /api/categories/:id | Hapus kategori |

---

## Technical Questions

### 1. Database tables apa yang dibuat dan kenapa?

Saya membuat 2 tabel:

**Categories:** id, name, color, created_at, updated_at
**Todos:** id, title, description, completed, priority, due_date, category_id (FK), created_at, updated_at

**Relasi:** One-to-Many (1 kategori punya banyak todos)

**Alasan:** Struktur sederhana, normalized, efisien untuk query yang dibutuhkan.

### 2. Bagaimana handle pagination dan filtering?

**Pagination:** LIMIT dan OFFSET (`OFFSET = (page-1) * limit`)

**Filtering:** WHERE clause untuk search (ILIKE), category_id, completed, priority

**Index yang ditambahkan:**
- idx_todos_completed, idx_todos_priority, idx_todos_category_id, idx_todos_created_at
- GIN index untuk full-text search

### 3. Bagaimana implementasi responsive design?

**Breakpoints:** Mobile (<576px), Tablet (576-992px), Desktop (>992px)

**Adaptasi:** Kolom tabel hide di layar kecil, filter stack vertikal di mobile, horizontal scroll untuk tabel

**Komponen Ant Design:** Row/Col dengan breakpoints, Table responsive property, scroll prop

### 4. Bagaimana struktur React components?

```
App → TodoProvider (Context) → Layout → TodoList → TodoForm + CategoryManager
```

**State management:** Global state di TodoContext (todos, categories, pagination, filters), local state untuk UI

### 5. Arsitektur backend apa yang dipilih?

**Layered Architecture:**
```
Handlers → Services → Repository → Database
```

**Struktur:** handlers (HTTP), services (business logic), repository (database), models (data structures)

**Error handling:** Custom error types, service return meaningful errors, handler translate ke HTTP status

### 6. Bagaimana handle validasi data?

**Frontend:** Ant Design Form rules, TypeScript type checking
**Backend:** Gin binding tags, service layer validation

**Rules:** Title required (max 255), priority harus high/medium/low, color valid hex

**Alasan validasi di keduanya:** Frontend untuk UX, backend untuk security

### 7. Apa yang di-unit test dan kenapa?

**Yang ditest:** TodoService dan CategoryService (semua CRUD operations)

**Edge cases:** Empty fields, invalid values, non-existent IDs, partial updates

**Struktur:** Table-driven tests, mock repository dengan testify/mock

**Alasan service layer:** Contains business logic, mudah mock dependencies, high value tests

### 8. Jika punya waktu lebih, apa yang akan ditambah?

**Technical debt:** Integration tests, E2E tests, proper logging, request tracing

**Fitur baru:** Authentication (JWT), notifications, recurring todos, bulk operations, dark mode

**Refactoring:** Rate limiting, caching (Redis), API versioning

