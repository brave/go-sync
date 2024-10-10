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
	specifics BYTEA STORAGE EXTERNAL NOT NULL,
	client_defined_unique_tag TEXT STORAGE PLAIN,
	server_defined_unique_tag TEXT STORAGE PLAIN,
	name TEXT STORAGE PLAIN,
	originator_cache_guid TEXT STORAGE PLAIN,
	originator_client_item_id TEXT STORAGE PLAIN,
	parent_id TEXT STORAGE PLAIN,
	non_unique_name TEXT STORAGE PLAIN,
	unique_position BYTEA STORAGE PLAIN,
	folder BOOLEAN,
	deleted BOOLEAN NOT NULL,
	PRIMARY KEY (id, chain_id),
	UNIQUE (chain_id, client_defined_unique_tag)
)
PARTITION BY RANGE (chain_id);
CREATE INDEX entities_chain_id_data_type_mtime_idx ON entities (chain_id, data_type, mtime);

SELECT partman.create_parent(
	p_parent_table := 'public.entities',
	p_control := 'chain_id',
	p_interval := '3500'
);

CREATE EXTENSION IF NOT EXISTS pg_cron;

SELECT cron.schedule('@hourly', $$CALL partman.run_maintenance_proc()$$);
