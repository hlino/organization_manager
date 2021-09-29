CREATE TABLE organizations
(
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    creation_date TIMESTAMP,
    employee_count INTEGER,
    is_public BOOLEAN
);
