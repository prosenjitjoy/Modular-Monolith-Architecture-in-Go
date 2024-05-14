-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA ordering;

CREATE TABLE ordering.orders
(
  id          text NOT NULL,
  customer_id text NOT NULL,
  payment_id  text NOT NULL,
  invoice_id  text NOT NULL,
  shopping_id text NOT NULL,
  items       bytea NOT NULL,
  status      text NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT NOW(),
  updated_at  timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_orders_trgr BEFORE UPDATE ON ordering.orders
FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_orders_trgr BEFORE UPDATE ON ordering.orders
FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS ordering CASCADE;
-- +goose StatementEnd
