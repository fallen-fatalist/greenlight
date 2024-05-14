CREATE TABLE IF NOT EXISTS movies (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE movies ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);

ALTER TABLE movies ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));

ALTER TABLE movies ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);

INSERT INTO movies (title, year, runtime, genres) VALUES (
    'Spider man',
    2002,
    102,
    -- adventure genre should be added
    '{"sci-fi", "action"}'
);

INSERT INTO movies (title, year, runtime, genres) VALUES (
    'Attack of the titans',
    2013,
    1,
    -- should be added "anime" genre
    '{"adventure", "action", "fantasy", "drama", "cartoon"}'
);

INSERT INTO movies (title, year, runtime, genres) VALUES (
    'Grimgar of the fantasy and ash',
    -- should be 2016 year
    2015,
    1,
    '{"anime", "adventure", "action", "fantasy", "cartoon"}'
);