BEGIN;

INSERT INTO accounts
(first_name, last_name, dob, gender, email, confirmed_email, phone_number, confirmed_phone, passhash)
VALUES
('Daniel', 'Martinez Olivas', '1985-05-08', 'M', 'example1@gmail.com', true, '6666665', true, crypt('asdf', gen_salt('bf', 8))),
('Jose', 'Saldivia Suarez', '1985-05-07', 'M', 'example2@gmail.com', true, '6666666', true, crypt('asdf', gen_salt('bf', 8)));

-- INSERT INTO roles
--     ()
-- VALUES
--     ();

COMMIT;