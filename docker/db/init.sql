CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    birth_date DATE,
    gender VARCHAR(10),
    interests TEXT,
    city VARCHAR(50),
    username VARCHAR(50) UNIQUE,
    password VARCHAR(100)
  );
  