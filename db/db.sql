CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE admin (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_link TEXT,
    admin_id INT REFERENCES admin(id) ON DELETE CASCADE
);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE
);



