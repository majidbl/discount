CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS gift_charges CASCADE;
DROP TABLE IF EXISTS reports CASCADE;

CREATE TABLE gift_charges (
    id                    SERIAL PRIMARY KEY,
    code                  VARCHAR(255) UNIQUE,
    validity_period_start TIMESTAMPTZ,
    validity_period_end   TIMESTAMPTZ,
    amount                BIGINT,
    max_usage_count       INT,
    created_at            TIMESTAMPTZ NOT NULL,
    updated_at            TIMESTAMPTZ
);

CREATE TABLE reports (
    id            SERIAL PRIMARY KEY,
    mobile        VARCHAR(50) NOT NULL,
    gift_code     VARCHAR(255) NOT NULL,
    charge_amount BIGINT NOT NULL,
    report_time   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL
);
