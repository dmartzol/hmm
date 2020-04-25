BEGIN;

INSERT INTO accounts
    (first_name, last_name, passhash, dob)
VALUES
    ('John', 'Dow', crypt('randomPassword', gen_salt('bf')), '1980-01-01');

COMMIT;