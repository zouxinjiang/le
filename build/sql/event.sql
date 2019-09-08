CREATE TABLE "event"(
	id SERIAL PRIMARY KEY,
	userid INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	happentime TIMESTAMPTZ,
	place TEXT,
	contenttype TEXT,
	content TEXT,
	summary TEXT,
	createat TIMESTAMPTZ,
	updateat TIMESTAMPTZ
);