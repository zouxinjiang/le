CREATE TABLE "event"(
	id SERIAL PRIMARY KEY,
	happentime TIMETZ,
	place TEXT,
	contenttype TEXT,
	content TEXT,
	summary TEXT
);