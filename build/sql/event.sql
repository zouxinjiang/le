CREATE TABLE "event"(
	id SERIAL PRIMARY KEY,
	userid INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	happentime TIMETZ,
	place TEXT,
	contenttype TEXT,
	content TEXT,
	summary TEXT,
	createat TIMETZ,
	updateat TIMETZ
);