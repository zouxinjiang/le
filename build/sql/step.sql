CREATE TABLE "step" (
	id SERIAL PRIMARY KEY,
	parentid INTEGER REFERENCES step(id) ON DELETE CASCADE
	title TEXT,
	stepnumber Integer,
	description TEXT,
	planstart TIMETZ,
	planend TIMETZ,
	summary TEXT
);