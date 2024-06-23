CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name VARCHAR(100) NOT NULL
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    date DATE NOT NULL,
    transaction DECIMAL(10, 2) NOT NULL
);

INSERT INTO users (id, name, email) VALUES
(80208352, 'Eduard Reyes', 'reydward@gmail.com');

INSERT INTO accounts (id, user_id, name) values
(111, 80208352, 'Cuenta de ahorros'),
(112, 80208352, 'Cuenta corriente'),
(113, 80208352, 'Tarjeta de crédito'),
(114, 80208352, 'Cuenta de inversión');

INSERT INTO transactions (account_id, date, amount) values
(111, '2024-07-15', 60.50),
(111, '2024-07-28', -10.30),
(111, '2024-08-02', -20.46),
(111, '2024-08-13', 10),
(211, '2024-06-25', 5000.00);
