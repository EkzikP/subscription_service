CREATE TABLE IF NOT EXISTS service.subscriptions(
    service_name text NOT NULL,
    price integer NOT NULL,
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date
);

CREATE INDEX IF NOT EXISTS subscriptions_user_id_idx ON service.subscriptions(user_id);
CREATE INDEX IF NOT EXISTS subscriptions_service_name_idx ON service.subscriptions(service_name);