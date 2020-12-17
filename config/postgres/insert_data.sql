BEGIN;

INSERT INTO emails
(email, confirmed)
VALUES
('email001@example.com', false),
('email002@example.com', false),
('email006@example.com', false);

INSERT INTO phone_numbers
("number", confirmed)
VALUES
('+1 (320)-414-7788', false),
('+1 (310)-422-1128', false);

INSERT INTO addresses
(country, city, state_code, street, zip_code, "type")
VALUES
('US', 'Chicago', 'IL', '1100 W Armitage Ave', '60646', 0),
('US', 'Chicago', 'IL', '1100 W Armitage Ave', '60646', 0),
('US', 'Chicago', 'IL', '1100 W Armitage Ave', '60646', 0);

INSERT INTO accounts
    (
        email_id,
        address_id,
        phone_number_id,
        first_name,
        last_name,
        dob,
        gender,
        active,
        passhash,
        failed_logins_count,
        review_time
    )
VALUES
    (1, 1, NULL, 'Daniel', 'Martinez Olivas', '1999-01-01', 'm', false, crypt('asdf', gen_salt('bf', 8)), 0, null),
    (1, 1, NULL, 'Jose', 'Saldivia Suárez', '1985-11-09', 'M', false, crypt('asdf', gen_salt('bf', 8)), 0, null);

COMMIT;