-- Tabel untuk pengguna
CREATE TABLE users (
    id UUID PRIMARY KEY, 
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- Tabel untuk proyek
CREATE TABLE projects (
    id UUID PRIMARY KEY, 
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL, 
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Tabel untuk tugas
CREATE TABLE tasks (
    id UUID PRIMARY KEY, 
    project_id UUID NOT NULL, 
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status ENUM('todo', 'in_progress', 'done') NOT NULL,
    deadline DATE,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);