CREATE TABLE "flag" (
	id SERIAL PRIMARY KEY,
	userid INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	name TEXT,
	description TEXT,
	createat TIMESTAMPTZ,
	updateat TIMESTAMPTZ,
	summary TEXT
);