-- null values are not appreciated in GoLang 
-- So all columns either not null or have default vals
CREATE TABLE IF NOT EXISTS cars (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null default NOW(),
    name text not null,
    body text not null,
    brake_system text not null,
    aspiration text not null,
    horsepower numeric not null,
    mpg numeric not null,
    cylinders integer not null,
    acceleration numeric not null,
    displacement numeric not null,
    origin text not null,
    version integer NOT NULL DEFAULT 1
);

