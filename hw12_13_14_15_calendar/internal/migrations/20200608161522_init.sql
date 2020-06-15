-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
	id serial NOT NULL,
	title varchar(255) NOT NULL,
	datetime timestamptz NOT NULL,
	duration int8 NOT NULL,
	description text NULL,
	owner_id int4 NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NULL,
	CONSTRAINT events_pk PRIMARY KEY (id)
);
CREATE INDEX events_datetime_idx ON events USING btree (datetime);
CREATE UNIQUE INDEX events_owner_id_idx ON events USING btree (owner_id, datetime);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
