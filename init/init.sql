CREATE TABLE customer (
    customer_number SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    phone           VARCHAR(20),
    birth_date      DATE,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customer_note (
    id              SERIAL PRIMARY KEY,
    customer_number INTEGER NOT NULL REFERENCES customer(customer_number) ON DELETE CASCADE,
    note            TEXT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO customer (name, email, phone, birth_date)
VALUES ('John Doe', 'john.doe@example.com', '08123456789', '1990-05-15');

INSERT INTO customer_note (customer_number, note)
VALUES (1, 'Customer called about billing issue');
