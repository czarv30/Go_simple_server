drop database if exists school_records;
create database school_records;

\c school_records;

create table students (

    student_id SERIAL PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    birth_date DATE NOT NULL

);

INSERT INTO students (first_name, last_name, birth_date) VALUES
    ('Amy', 'Winehouse', '2015-01-16'),
    ('Bob', 'Marley', '2013-06-25'),
    ('Charles', 'Ray', '2017-08-30');


CREATE TYPE grade_type AS ENUM ('A+', 'A', 'A-', 'B+', 'B', 'B-', 'C+', 'C', 'C-', 'D+', 'D', 'D-', 'F');

CREATE TYPE semester_type AS ENUM ('Fall','Winter','Spring','Summer');

create table class_grades (
    class_id SERIAL PRIMARY KEY,
    class_name VARCHAR NOT NULL,
    student_id INTEGER REFERENCES students(student_id),
    semester semester_type NOT NULL,
    year INTEGER NOT NULL CHECK (year >= 2000 AND year <= 2100),
    grade grade_type NOT NULL
);