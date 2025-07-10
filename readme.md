# Taskify API

Taskify adalah RESTful API backend yang dibangun menggunakan **Go**, **Gin**, dan **GORM**. Aplikasi ini menyediakan fitur manajemen **user**, **project**, dan **task**, dilengkapi autentikasi JWT dan dukungan Docker untuk deployment.

---

## 🧩 Fitur

* 🔐 Register & Login dengan hashing password (bcrypt)
* 🧾 Manajemen Proyek (CRUD) per user
* ✅ Manajemen Tugas dalam proyek (CRUD)
* 📦 JWT Middleware (autentikasi & otorisasi)
* 🐳 Docker support

---

## 🚀 Quick Start

### 1. Clone & Build

```bash
git clone https://github.com/username/taskify.git
cd taskify
go mod tidy
go run main.go
```

### 2. Siapkan Database MySQL / MariaDB

Pastikan service database berjalan di `localhost:3306` dan kamu punya database kosong bernama `taskify`.

```sql
CREATE DATABASE taskify;
```

### 3. Buat File `.env`

Contoh isi file `.env`:

```ini
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=taskify
JWT_SECRET=rahasia
```

Letakkan `.env` di root proyek.

### 4. Jalankan aplikasi

```bash
go mod tidy
go run main.go
```

Server akan berjalan di `localhost:8080`

---

## 📬 Endpoint List (Postman)

### 🔐 1. REGISTER

* **Method**: POST
* **URL**: `/register`
* **Body (JSON)**:

```json
{
  "name": "Hamas",
  "email": "hamas@mail.com",
  "password": "rahasia123"
}
```

### 🔐 2. LOGIN

* **Method**: POST
* **URL**: `/login`
* **Body (JSON)**:

```json
{
  "email": "hamas@mail.com",
  "password": "rahasia123"
}
```

* **Response**:

```json
{
  "message": "Login successful",
  "token": "<JWT_TOKEN>",
  "user_id": "<UUID>"
}
```

Gunakan token untuk request berikutnya:

```
Authorization: Bearer <JWT_TOKEN>
```

---

## 📁 PROJECTS (Harus Login)

### ✅ 3. CREATE PROJECT

* **Method**: POST
* **URL**: `/projects`
* **Headers**: Authorization
* **Body (JSON)**:

```json
{
  "name": "Project Alpha",
  "description": "Proyek penting",
  "created_by": "<UUID_USER_DARI_LOGIN>"
}
```

### ✅ 4. GET ALL PROJECTS

* **Method**: GET
* **URL**: `/projects`
* **Headers**: Authorization

### ✅ 5. GET PROJECT BY ID

* **Method**: GET
* **URL**: `/projects/{id}`
* **Headers**: Authorization

### ✅ 6. UPDATE PROJECT

* **Method**: PUT
* **URL**: `/projects/{id}`
* **Headers**: Authorization
* **Body (JSON)**:

```json
{
  "name": "Project Beta",
  "description": "Deskripsi baru"
}
```

### ✅ 7. DELETE PROJECT

* **Method**: DELETE
* **URL**: `/projects/{id}`
* **Headers**: Authorization

---

## ✅ TASKS (Dalam Project, Harus Login)

### ✅ 8. CREATE TASK

* **Method**: POST
* **URL**: `/projects/{project_id}/tasks`
* **Headers**: Authorization
* **Body (JSON)**:

```json
{
  "title": "Tugas 1",
  "description": "Kerjakan API",
  "status": "todo",
  "deadline": "2025-07-20"
}
```

### ✅ 9. GET TASKS BY PROJECT

* **Method**: GET
* **URL**: `/projects/{project_id}/tasks`
* **Headers**: Authorization

### ✅ 10. GET TASK BY ID

* **Method**: GET
* **URL**: `/projects/{project_id}/tasks/{task_id}`
* **Headers**: Authorization

### ✅ 11. UPDATE TASK

* **Method**: PUT
* **URL**: `/projects/{project_id}/tasks/{task_id}`
* **Headers**: Authorization
* **Body (JSON)**:

```json
{
  "title": "Update Tugas",
  "description": "Revisi backend",
  "status": "in_progress",
  "deadline": "2025-07-22"
}
```

### ✅ 12. DELETE TASK

* **Method**: DELETE
* **URL**: `/projects/{project_id}/tasks/{task_id}`
* **Headers**: Authorization

---

## 📌 Status Enum untuk Task

* `todo`
* `in_progress`
* `done`

Gunakan salah satu nilai di atas saat membuat atau mengupdate task.

---

## 🛠️ Teknologi

* Go
* Gin Web Framework
* GORM ORM
* MySQL / MariaDB
* JWT
* Docker (opsional)

---

## 🧑‍💻 Kontributor

Muhammad Hamas - 2025

---

## 📄 Lisensi

MIT License
