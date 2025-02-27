-- First, let's create the audit table to store all change records
CREATE TABLE audit.logged_actions (
    event_id BIGSERIAL PRIMARY KEY,
    schema_name TEXT NOT NULL,
    table_name TEXT NOT NULL,
    table_pk TEXT,              -- Primary key column(s) value
    table_pk_column TEXT,       -- Primary key column(s) name
    action_tstamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    action TEXT NOT NULL CHECK (action IN ('I','D','U')), -- Insert, Delete, Update
    original_data JSONB,        -- Previous data for updates and deletes
    new_data JSONB,             -- New data for inserts and updates
    changed_fields JSONB,       -- Changed fields for updates
    transaction_id BIGINT,      -- Transaction ID
    application_name TEXT,      -- Application name from current_setting
    client_addr INET NULL,      -- Client IP address
    client_port INTEGER NULL,   -- Client port
    session_user_name TEXT,     -- Session user
    current_user_name TEXT      -- Current user performing action
);

-- Create comments on the audit table for better documentation
COMMENT ON TABLE audit.logged_actions IS 'History of auditable actions on audited tables';
COMMENT ON COLUMN audit.logged_actions.event_id IS 'Unique identifier for each auditable event';
COMMENT ON COLUMN audit.logged_actions.schema_name IS 'Database schema audited table is in';
COMMENT ON COLUMN audit.logged_actions.table_name IS 'Non-schema-qualified table name of table event occured in';
COMMENT ON COLUMN audit.logged_actions.table_pk IS 'Primary key value of the affected row';
COMMENT ON COLUMN audit.logged_actions.table_pk_column IS 'Name of the primary key column';
COMMENT ON COLUMN audit.logged_actions.action_tstamp IS 'Transaction start timestamp for tx in which the audited event occurred';
COMMENT ON COLUMN audit.logged_actions.action IS 'Action type; I = insert, D = delete, U = update';
COMMENT ON COLUMN audit.logged_actions.original_data IS 'Record value before modification (for updates and deletes)';
COMMENT ON COLUMN audit.logged_actions.new_data IS 'New record value (for inserts and updates)';
COMMENT ON COLUMN audit.logged_actions.changed_fields IS 'Updated fields (for updates only)';
COMMENT ON COLUMN audit.logged_actions.client_addr IS 'IP address of client that issued query';
COMMENT ON COLUMN audit.logged_actions.client_port IS 'Port address of client that issued query';
COMMENT ON COLUMN audit.logged_actions.session_user_name IS 'Login / session user whose actions are being audited';
COMMENT ON COLUMN audit.logged_actions.current_user_name IS 'User who actually performed the action';

-- Create index for better query performance
CREATE INDEX logged_actions_schema_table_idx 
ON audit.logged_actions(schema_name, table_name);

CREATE INDEX logged_actions_action_tstamp_idx 
ON audit.logged_actions(action_tstamp);

CREATE INDEX logged_actions_action_idx 
ON audit.logged_actions(action);

-- Now, let's create the audit trigger function
CREATE OR REPLACE FUNCTION audit.if_modified_func() 
RETURNS TRIGGER AS $body$
DECLARE
    v_old_data JSONB;
    v_new_data JSONB;
    v_changed_fields JSONB;
    v_primary_key_column TEXT;
    v_primary_key_value TEXT;
BEGIN
    IF TG_WHEN <> 'AFTER' THEN
        RAISE EXCEPTION 'audit.if_modified_func() may only run as an AFTER trigger';
    END IF;

    -- Determine primary key column
    SELECT a.attname INTO v_primary_key_column
    FROM pg_index i
    JOIN pg_attribute a ON a.attrelid = i.indrelid
        AND a.attnum = ANY(i.indkey)
    WHERE i.indrelid = TG_RELID
        AND i.indisprimary;

    -- Get the primary key value
    IF TG_OP = 'UPDATE' OR TG_OP = 'DELETE' THEN
        EXECUTE 'SELECT ($1).' || quote_ident(v_primary_key_column)
        INTO v_primary_key_value
        USING OLD;
    ELSIF TG_OP = 'INSERT' THEN
        EXECUTE 'SELECT ($1).' || quote_ident(v_primary_key_column)
        INTO v_primary_key_value
        USING NEW;
    END IF;

    -- Convert to JSON data formats based on operation
    IF (TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
        v_old_data = to_jsonb(OLD);
    END IF;
    
    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
        v_new_data = to_jsonb(NEW);
    END IF;

    -- For UPDATE operations, calculate the changed fields
    IF (TG_OP = 'UPDATE') THEN
        -- Calculate changed fields by comparing old and new data
        SELECT jsonb_object_agg(key, value) INTO v_changed_fields
        FROM jsonb_each(v_new_data)
        WHERE NOT jsonb_contains(v_old_data, jsonb_build_object(key, value));
    END IF;

    -- Insert audit log entry
    INSERT INTO audit.logged_actions (
        schema_name,
        table_name,
        table_pk,
        table_pk_column,
        action,
        original_data,
        new_data,
        changed_fields,
        transaction_id,
        application_name,
        client_addr,
        client_port,
        session_user_name,
        current_user_name
    ) VALUES (
        TG_TABLE_SCHEMA::TEXT,                          -- schema_name
        TG_TABLE_NAME::TEXT,                            -- table_name
        v_primary_key_value,                            -- table_pk
        v_primary_key_column,                           -- table_pk_column
        substring(TG_OP,1,1),                           -- action
        v_old_data,                                     -- original_data
        v_new_data,                                     -- new_data
        v_changed_fields,                               -- changed_fields
        txid_current(),                                 -- transaction_id
        current_setting('application_name'),            -- application_name
        inet_client_addr(),                             -- client_addr
        inet_client_port(),                             -- client_port
        session_user::TEXT,                             -- session_user_name
        current_user::TEXT                              -- current_user_name
    );

    RETURN NULL; -- For AFTER triggers, return value is ignored
END;
$body$
LANGUAGE plpgsql
SECURITY DEFINER;

COMMENT ON FUNCTION audit.if_modified_func() IS 'Trigger function that logs changes to the audit.logged_actions table';

-- Function to automatically create audit triggers for all tables in a schema
CREATE OR REPLACE FUNCTION audit.create_audit_triggers_for_schema(target_schema TEXT) 
RETURNS VOID AS $$
DECLARE
    table_record RECORD;
BEGIN
    -- Ensure audit schema exists
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'audit') THEN
        CREATE SCHEMA audit;
    END IF;

    -- Loop through all tables in the schema and add audit triggers
    FOR table_record IN 
        SELECT tablename 
        FROM pg_tables 
        WHERE schemaname = target_schema 
        AND tablename NOT LIKE 'pg_%' 
        AND tablename <> 'logged_actions'
    LOOP
        -- Create the audit trigger for the current table
        EXECUTE format('
            CREATE TRIGGER audit_trigger_row
            AFTER INSERT OR UPDATE OR DELETE ON %I.%I
            FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();
            
            CREATE TRIGGER audit_trigger_stmt
            AFTER TRUNCATE ON %I.%I
            FOR EACH STATEMENT EXECUTE FUNCTION audit.if_modified_func();
            ', 
            target_schema, table_record.tablename,
            target_schema, table_record.tablename
        );
        
        RAISE NOTICE 'Added audit triggers to table: %.%', target_schema, table_record.tablename;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.create_audit_triggers_for_schema(TEXT) IS 'Function to add audit triggers to all tables in a schema';

-- Function to create audit triggers for all tables in the database (except system schemas)
CREATE OR REPLACE FUNCTION audit.create_audit_triggers_for_all_tables() 
RETURNS VOID AS $$
DECLARE
    schema_record RECORD;
BEGIN
    -- Ensure audit schema exists
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'audit') THEN
        CREATE SCHEMA audit;
    END IF;

    -- Loop through all non-system schemas and add audit triggers
    FOR schema_record IN 
        SELECT nspname 
        FROM pg_namespace 
        WHERE nspname NOT LIKE 'pg_%' 
        AND nspname <> 'information_schema'
        AND nspname <> 'audit'
    LOOP
        -- Call the function to create audit triggers for all tables in the current schema
        PERFORM audit.create_audit_triggers_for_schema(schema_record.nspname);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.create_audit_triggers_for_all_tables() IS 'Function to add audit triggers to all tables in all non-system schemas';

-- Example usage:
-- STEP 1: Ensure the audit schema exists
CREATE SCHEMA IF NOT EXISTS audit;

-- STEP 2: Create the audit table and functions defined above

-- STEP 3: Add audit triggers to all tables in the database
-- SELECT audit.create_audit_triggers_for_all_tables();

-- STEP 4: Or, to add audit triggers to tables in a specific schema only:
-- SELECT audit.create_audit_triggers_for_schema('public');

-- STEP 5: To add audit trigger to a single table manually:
/*
CREATE TRIGGER audit_trigger_row
AFTER INSERT OR UPDATE OR DELETE ON schema_name.table_name
FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();

CREATE TRIGGER audit_trigger_stmt
AFTER TRUNCATE ON schema_name.table_name
FOR EACH STATEMENT EXECUTE FUNCTION audit.if_modified_func();
*/