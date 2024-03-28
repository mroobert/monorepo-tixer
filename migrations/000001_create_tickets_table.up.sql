CREATE TABLE IF NOT EXISTS tickets (
    id bigserial PRIMARY KEY NOT NULL, 
    public_id char(12) NOT NULL, 
    title text NOT NULL,
    price integer NOT NULL,
    version integer NOT NULL DEFAULT 1,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
