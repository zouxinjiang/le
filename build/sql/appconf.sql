CREATE TABLE "appconf" (
    id SERIAL PRIMARY KEY,
    name text UNIQUE ,
    value text,
    state int
);
