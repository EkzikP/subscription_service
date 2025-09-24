CREATE TABLE IF NOT EXISTS subscriptions(
    service_name text NOT NULL,
    price integer NOT NULL CHECK (price > 0),
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date NULL
);

CREATE INDEX IF NOT EXISTS subscriptions_user_id_idx ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS subscriptions_service_name_idx ON subscriptions(service_name);
CREATE INDEX IF NOT EXISTS subscriptions_dates_idx ON subscriptions(start_date, end_date);