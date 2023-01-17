CREATE TABLE players (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     score INTEGER NOT NULL
);
create table other_scores(
    locked integer not null,
    unknown_locked integer not null,
    unknown integer not null
);
CREATE TABLE place_scores (
      id SERIAL PRIMARY KEY,
      place INTEGER NOT NULL,
      score INTEGER NOT NULL
);
CREATE TABLE unknown_names (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL
);