CREATE TABLE chains (
	id BIGSERIAL PRIMARY KEY,
	client_id BYTEA NOT NULL,
	last_usage_time TIMESTAMP NOT NULL,
	UNIQUE (client_id)
);

CREATE TABLE dynamo_migration_statuses (
	chain_id BIGINT REFERENCES chains(id),
	data_type INTEGER,
	earliest_mtime BIGINT NOT NULL,
	PRIMARY KEY (chain_id, data_type)
);

CREATE TABLE entities (
	id UUID,
	chain_id BIGINT NOT NULL REFERENCES chains(id),
	data_type INTEGER NOT NULL,
	ctime BIGINT NOT NULL,
	mtime BIGINT NOT NULL,
	specifics BYTEA STORAGE EXTERNAL NOT NULL,
	deleted BOOL NOT NULL,
	client_defined_unique_tag TEXT STORAGE PLAIN,
	server_defined_unique_tag TEXT STORAGE PLAIN,
	folder BOOLEAN,
	version BIGINT NOT NULL,
	name TEXT STORAGE PLAIN,
	originator_cache_guid TEXT STORAGE PLAIN,
	originator_client_item_id TEXT STORAGE PLAIN,
	parent_id TEXT STORAGE PLAIN,
	non_unique_name TEXT STORAGE PLAIN,
	unique_position BYTEA STORAGE PLAIN,
	PRIMARY KEY (id, chain_id),
	UNIQUE (chain_id, client_defined_unique_tag)
);
CREATE INDEX entities_chain_id_idx ON entities (chain_id);
CREATE INDEX entities_data_type_mtime_idx ON entities (data_type, mtime);
-- or maybe make a partial index for history entities and mtime, while keeping the chainid datattype and mtime index
-- CREATE INDEX entities_chain_id_data_type_mtime_idx ON entities (chain_id, data_type, mtime);
