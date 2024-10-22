CREATE SCHEMA IF NOT EXISTS partman;
CREATE EXTENSION IF NOT EXISTS pg_partman SCHEMA partman;

CREATE TABLE chains (
	id BIGSERIAL PRIMARY KEY,
	last_usage_time TIMESTAMP NOT NULL,
	client_id BYTEA NOT NULL,
	UNIQUE (client_id)
);

CREATE TABLE dynamo_migration_statuses (
	chain_id BIGINT REFERENCES chains(id) ON DELETE CASCADE,
	-- null earliest_mtime indicates that all entities have been migrated
	earliest_mtime BIGINT,
	data_type INTEGER,
	PRIMARY KEY (chain_id, data_type)
);

CREATE TABLE entities (
	id UUID,
	chain_id BIGINT NOT NULL REFERENCES chains(id) ON DELETE CASCADE,
	ctime BIGINT NOT NULL,
	mtime BIGINT NOT NULL,
	version BIGINT NOT NULL,
	data_type INTEGER NOT NULL,
	specifics BYTEA NOT NULL,
	client_defined_unique_tag TEXT,
	server_defined_unique_tag TEXT,
	name TEXT,
	originator_cache_guid TEXT,
	originator_client_item_id TEXT,
	parent_id TEXT,
	non_unique_name TEXT,
	unique_position BYTEA,
	folder BOOLEAN,
	deleted BOOLEAN NOT NULL,
	PRIMARY KEY (id, chain_id),
	UNIQUE (chain_id, client_defined_unique_tag)
)
PARTITION BY RANGE (chain_id);

ALTER TABLE entities ALTER specifics SET STORAGE EXTERNAL;
ALTER TABLE entities ALTER client_defined_unique_tag SET STORAGE PLAIN;
ALTER TABLE entities ALTER server_defined_unique_tag SET STORAGE PLAIN;
ALTER TABLE entities ALTER name SET STORAGE PLAIN;
ALTER TABLE entities ALTER originator_cache_guid SET STORAGE PLAIN;
ALTER TABLE entities ALTER originator_client_item_id SET STORAGE PLAIN;
ALTER TABLE entities ALTER parent_id SET STORAGE PLAIN;
ALTER TABLE entities ALTER non_unique_name SET STORAGE PLAIN;
ALTER TABLE entities ALTER unique_position SET STORAGE PLAIN;

CREATE INDEX entities_chain_id_data_type_mtime_idx ON entities (chain_id, data_type, mtime);

DO $$
BEGIN
	-- for vanilla postgres
	PERFORM partman.create_parent(
		p_parent_table := 'public.entities',
		p_control := 'chain_id',
		p_interval := '3500',
		p_type := 'range'
	);
EXCEPTION WHEN OTHERS THEN
	-- for Aurora
	PERFORM partman.create_parent(
		p_parent_table := 'public.entities',
		p_control := 'chain_id',
		p_interval := '3500',
		p_type := 'native'
	);
END $$;

CREATE EXTENSION IF NOT EXISTS pg_cron;

SELECT cron.schedule('@hourly', $$CALL partman.run_maintenance_proc()$$);
