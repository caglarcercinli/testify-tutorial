DROP TABLE IF EXISTS stockprices;

CREATE TABLE stockprices (
    id SERIAL PRIMARY KEY,
    prices DECIMAL
);