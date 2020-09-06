BEGIN;

insert into accounts
(first_name, last_name, dob, gender, email, confirmed_email, phone_number, confirmed_phone, passhash)
values
('Daniel', 'Martinez Olivas', '1985-05-08', 'M', 'example1@gmail.com', true, '6666665', true, crypt('randompass', gen_salt('bf', 8))),
('Jose', 'Saldivia Suarez', '1985-05-07', 'M', 'example2@gmail.com', true, '6666666', true, crypt('randompass', gen_salt('bf', 8)));

INSERT INTO accounts
    (first_name, last_name, passhash, dob)
VALUES
    ('John', 'Dow', crypt('randomPassword', gen_salt('bf')), '1980-01-01');

INSERT INTO roles
    ()
VALUES
    ();

COMMIT;