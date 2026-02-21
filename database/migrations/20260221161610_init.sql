-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id INT REFERENCES departments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_dept_name_parent_id ON departments (name, parent_id) WHERE parent_id IS NOT NULL;
CREATE UNIQUE INDEX idx_dept_name_root ON departments (name) WHERE parent_id IS NULL;

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    department_id INT NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;
-- +goose StatementEnd




