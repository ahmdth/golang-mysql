-- delete database if it's created before
DROP DATABASE IF EXISTS blog;

-- create new database
CREATE DATABASE blog;

USE blog;

-- create table
CREATE TABLE posts(
    id INT(11) NOT NULL PRIMARY KEY ,
    title VARCHAR(255) UNIQUE NOT NULL,
    body TEXT NOT NULL,
    is_published DATE
)