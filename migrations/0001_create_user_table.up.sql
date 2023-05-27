CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    email    VARCHAR(200) NOT NULL,
    password VARCHAR(100) NOT NULL,
    salt     VARCHAR(100) NOT NULL

);
