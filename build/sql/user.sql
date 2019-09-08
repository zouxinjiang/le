CREATE TABLE "user"(
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE,
	icon TEXT,
	name TEXT,
	pwd BYTEA,
	mobile TEXT UNIQUE,
	email TEXT UNIQUE,
	uuid TEXT UNIQUE,
	state INTEGER,
	locktime TIMESTAMPTZ,
	lockreason TEXT,
	createat TIMESTAMPTZ,
	updateat TIMESTAMPTZ
);