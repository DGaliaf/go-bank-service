BEGIN;

CREATE TABLE IF NOT EXISTS public.user (
    id SERIAL PRIMARY KEY NOT NULL,
    balance INT NOT NULL DEFAULT 0,
    CONSTRAINT positive_balance CHECK ( balance >= 0 )
);

END;