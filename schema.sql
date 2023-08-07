CREATE TABLE IF NOT EXISTS places
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    budget NUMERIC NOT NULL,
	location TEXT NOT NULL,
	mood TEXT NOT NULL,
	cuisine TEXT NOT NULL,
	mealtime TEXT NOT NULL,
	rating NUMERIC NOT NULL,
	approved BOOL DEFAULT false
)