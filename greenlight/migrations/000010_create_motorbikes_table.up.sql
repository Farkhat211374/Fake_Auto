
CREATE TABLE IF NOT EXISTS motorbikes (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null default NOW(),
    name text not null,
    horsepower numeric not null,
    type text not null,
    weight numeric not null,
    third_place boolean default false,
    cylinders integer not null,
    acceleration numeric not null,
    displacement numeric not null,
    origin text not null,
    version integer NOT NULL DEFAULT 1
);

