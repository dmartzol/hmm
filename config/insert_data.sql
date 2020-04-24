BEGIN;

INSERT INTO accounts
    (first_name, last_name, email, password, dob)
VALUES
    ('John', 'Dow', 'randomemail@email.com', crypt('randomPassword', gen_salt('bf')), '1980-01-01');

COMMIT;