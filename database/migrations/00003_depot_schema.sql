-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA depot;

CREATE TABLE depot.shopping_lists
(
  id              text NOT NULL,
  order_id        text NOT NULL,
  stops           bytea NOT NULL,
  assigned_bot_id text NOT NULL,
  status          text NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT NOW(),
  updated_at      timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE INDEX shopping_lists_order_id_idx ON depot.shopping_lists (order_id);

CREATE INDEX shopping_lists_availability_idx ON depot.shopping_lists (status, created_at)
WHERE status = 'available';

CREATE TRIGGER created_at_shopping_lists_trgr BEFORE UPDATE ON depot.shopping_lists
FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_shopping_lists_trgr BEFORE UPDATE ON depot.shopping_lists
FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS depot CASCADE;
-- +goose StatementEnd
