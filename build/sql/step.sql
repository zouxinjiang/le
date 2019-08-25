CREATE TABLE "step" (
	id SERIAL PRIMARY KEY,
	parentid INTEGER REFERENCES step(id) ON DELETE CASCADE,
	title TEXT,
	description TEXT,
	planstart TIMETZ,
	planend TIMETZ,
	summary TEXT
);