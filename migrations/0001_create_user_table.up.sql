CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    email    VARCHAR(200) NOT NULL,
    password VARCHAR(100) NOT NULL,
    salt     VARCHAR(100) NOT NULL

);

INSERT INTO users (email, password, salt)
VALUES ('user1@example.com', 'password1', 'salt1'),
       ('user2@example.com', 'password2', 'salt2'),
       ('admin@example.com', 'adminpassword', 'adminsalt');
