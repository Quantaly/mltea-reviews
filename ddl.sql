CREATE TABLE tea
(
    id                  serial PRIMARY KEY,
    name                text UNIQUE NOT NULL,
    caffeinated         boolean     NOT NULL,
    currently_available boolean     NOT NULL DEFAULT true
);

CREATE TABLE review
(
    id       serial PRIMARY KEY,
    tea_id   integer  NOT NULL REFERENCES tea (id),
    reviewer text     NOT NULL,
    rating   smallint NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment  text     NOT NULL
);

CREATE TABLE faq
(
    ordinal  integer PRIMARY KEY,
    question text NOT NULL,
    answer   text NOT NULL
);
