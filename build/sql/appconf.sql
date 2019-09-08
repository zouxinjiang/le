CREATE TABLE "appconf" (
    id SERIAL PRIMARY KEY,
    name text UNIQUE ,
    value text,
    state int
);

INSERT INTO appconf VALUES
(1,'auth.twofactor.state',0,1),
(2,'auth.twofactor.email',1,1),
(3,'auth.twofactor.imagecode',1,1);
