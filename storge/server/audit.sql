-- Complete PostgreSQL Audit System with Automatic Triggers
-- This script creates a comprehensive audit system that:
-- 1. Records all data changes (inserts, updates, deletes) across all tables
-- 2. Automatically adds audit triggers to newly created tables
-- 3. Provides utility functions for managing the audit system

-- Step 1: Create the audit schema
CREATE SCHEMA IF NOT EXISTS audit;

-- Step 2: Create the audit table to store all change records
CREATE TABLE IF NOT EXISTS audit.logged_actions (
    event_id BIGSERIAL PRIMARY KEY,
    schema_name TEXT NOT NULL,
    table_name TEXT NOT NULL,
    table_pk TEXT,              -- Primary key column(s) value
    table_pk_column TEXT,       -- Primary key column(s) name
    action_tstamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    action TEXT NOT NULL CHECK (action IN ('I','D','U','T')), -- Insert, Delete, Update, Truncate
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
COMMENT ON COLUMN audit.logged_actions.action IS 'Action type; I = insert, D = delete, U = update, T = truncate';
COMMENT ON COLUMN audit.logged_actions.original_data IS 'Record value before modification (for updates and deletes)';
COMMENT ON COLUMN audit.logged_actions.new_data IS 'New record value (for inserts and updates)';
COMMENT ON COLUMN audit.logged_actions.changed_fields IS 'Updated fields (for updates only)';
COMMENT ON COLUMN audit.logged_actions.client_addr IS 'IP address of client that issued query';
COMMENT ON COLUMN audit.logged_actions.client_port IS 'Port address of client that issued query';
COMMENT ON COLUMN audit.logged_actions.session_user_name IS 'Login / session user whose actions are being audited';
COMMENT ON COLUMN audit.logged_actions.current_user_name IS 'User who actually performed the action';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS logged_actions_schema_table_idx 
ON audit.logged_actions(schema_name, table_name);

CREATE INDEX IF NOT EXISTS logged_actions_action_tstamp_idx 
ON audit.logged_actions(action_tstamp);

CREATE INDEX IF NOT EXISTS logged_actions_action_idx 
ON audit.logged_actions(action);

-- Step 3: Create the audit trigger function that handles data changes
CREATE OR REPLACE FUNCTION audit.if_modified_func() 
RETURNS TRIGGER AS $body$
DECLARE
    v_old_data JSONB;
    v_new_data JSONB;
    v_changed_fields JSONB;
    v_primary_key_column TEXT;
    v_primary_key_value TEXT;
    v_action TEXT;
BEGIN
    -- Set the action type
    IF TG_OP = 'INSERT' THEN
        v_action := 'I';
    ELSIF TG_OP = 'UPDATE' THEN
        v_action := 'U';
    ELSIF TG_OP = 'DELETE' THEN
        v_action := 'D';
    ELSIF TG_OP = 'TRUNCATE' THEN
        v_action := 'T';
    ELSE
        RAISE EXCEPTION 'Unsupported trigger operation: %', TG_OP;
    END IF;

    -- Validate trigger type
    IF TG_WHEN <> 'AFTER' THEN
        RAISE EXCEPTION 'audit.if_modified_func() may only run as an AFTER trigger';
    END IF;

    -- Determine primary key column if it exists
    -- This improved version will not fail if no primary key exists
    SELECT a.attname INTO v_primary_key_column
    FROM pg_index i
    JOIN pg_attribute a ON a.attrelid = i.indrelid
        AND a.attnum = ANY(i.indkey)
    WHERE i.indrelid = TG_RELID
        AND i.indisprimary
    LIMIT 1;

    -- Get the primary key value if primary key exists
    IF v_primary_key_column IS NOT NULL THEN
        IF TG_OP = 'UPDATE' OR TG_OP = 'DELETE' THEN
            EXECUTE 'SELECT ($1).' || quote_ident(v_primary_key_column)
            INTO v_primary_key_value
            USING OLD;
        ELSIF TG_OP = 'INSERT' THEN
            EXECUTE 'SELECT ($1).' || quote_ident(v_primary_key_column)
            INTO v_primary_key_value
            USING NEW;
        END IF;
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
        WHERE NOT jsonb_contains(v_old_data, jsonb_build_object(key, value))
           OR NOT v_old_data ? key;
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
        v_action,                                       -- action
        v_old_data,                                     -- original_data
        v_new_data,                                     -- new_data
        v_changed_fields,                               -- changed_fields
        txid_current(),                                 -- transaction_id
        current_setting('application_name', true),      -- application_name
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

-- Step 4: Function to automatically create audit triggers for all tables in a schema
CREATE OR REPLACE FUNCTION audit.create_audit_triggers_for_schema(target_schema TEXT) 
RETURNS VOID AS $$
DECLARE
    table_record RECORD;
BEGIN
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
            DROP TRIGGER IF EXISTS audit_trigger_row ON %I.%I;
            CREATE TRIGGER audit_trigger_row
            AFTER INSERT OR UPDATE OR DELETE ON %I.%I
            FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();
            
            DROP TRIGGER IF EXISTS audit_trigger_stmt ON %I.%I;
            CREATE TRIGGER audit_trigger_stmt
            AFTER TRUNCATE ON %I.%I
            FOR EACH STATEMENT EXECUTE FUNCTION audit.if_modified_func();
            ', 
            target_schema, table_record.tablename,
            target_schema, table_record.tablename,
            target_schema, table_record.tablename,
            target_schema, table_record.tablename
        );
        
        RAISE NOTICE 'Added audit triggers to table: %.%', target_schema, table_record.tablename;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.create_audit_triggers_for_schema(TEXT) IS 'Function to add audit triggers to all tables in a schema';

-- Step 5: Function to create audit triggers for all tables in the database (except system schemas)
CREATE OR REPLACE FUNCTION audit.create_audit_triggers_for_all_tables() 
RETURNS VOID AS $$
DECLARE
    schema_record RECORD;
BEGIN
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

    RAISE NOTICE 'Added audit triggers to all tables in all non-system schemas';
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.create_audit_triggers_for_all_tables() IS 'Function to add audit triggers to all tables in all non-system schemas';

-- Step 6: Function to automatically add audit triggers to newly created tables
CREATE OR REPLACE FUNCTION audit.add_audit_trigger_for_new_table()
RETURNS event_trigger AS $$
DECLARE
    obj RECORD;
    schema_name TEXT;
    table_name TEXT;
BEGIN
    -- Loop through objects created in this transaction
    FOR obj IN SELECT * FROM pg_event_trigger_ddl_commands() WHERE command_tag = 'CREATE TABLE'
    LOOP
        -- Extract schema and table name
        schema_name := obj.schema_name;
        table_name := obj.object_identity;
        
        -- Skip audit tables to avoid recursion
        IF schema_name = 'audit' THEN
            CONTINUE;
        END IF;
        
        -- Skip tables with table name that contains the full schema path
        IF position('.' in table_name) > 0 THEN
            table_name := substring(table_name from position('.' in table_name) + 1);
        END IF;
        
        -- Add audit triggers to the new table
        EXECUTE format('
            CREATE TRIGGER audit_trigger_row
            AFTER INSERT OR UPDATE OR DELETE ON %I.%I
            FOR EACH ROW EXECUTE FUNCTION audit.if_modified_func();
            
            CREATE TRIGGER audit_trigger_stmt
            AFTER TRUNCATE ON %I.%I
            FOR EACH STATEMENT EXECUTE FUNCTION audit.if_modified_func();',
            schema_name, table_name,
            schema_name, table_name
        );
        
        RAISE NOTICE 'Added audit triggers to new table: %.%', schema_name, table_name;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.add_audit_trigger_for_new_table() IS 'Function to automatically add audit triggers to newly created tables';

-- Step 7: Create the event trigger that fires when a table is created
DROP EVENT TRIGGER IF EXISTS audit_table_creation_trigger;
CREATE EVENT TRIGGER audit_table_creation_trigger
ON ddl_command_end
WHEN TAG IN ('CREATE TABLE')
EXECUTE FUNCTION audit.add_audit_trigger_for_new_table();

COMMENT ON EVENT TRIGGER audit_table_creation_trigger IS 'Event trigger that adds audit triggers to newly created tables';

-- Step 8: Function to enable or disable the event trigger (useful for maintenance)
CREATE OR REPLACE FUNCTION audit.toggle_audit_creation_trigger(enable BOOLEAN)
RETURNS VOID AS $$
BEGIN
    IF enable THEN
        ALTER EVENT TRIGGER audit_table_creation_trigger ENABLE;
        RAISE NOTICE 'Automatic audit triggers for new tables enabled';
    ELSE
        ALTER EVENT TRIGGER audit_table_creation_trigger DISABLE;
        RAISE NOTICE 'Automatic audit triggers for new tables disabled';
    END IF;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.toggle_audit_creation_trigger(BOOLEAN) IS 'Function to enable or disable automatic audit triggers for new tables';

-- Step 9: Function to set up the complete audit system
CREATE OR REPLACE FUNCTION audit.setup_complete_audit_system()
RETURNS VOID AS $$
BEGIN
    -- Add triggers to all existing tables
    PERFORM audit.create_audit_triggers_for_all_tables();
    
    -- Enable the event trigger for new tables
    PERFORM audit.toggle_audit_creation_trigger(true);
    
    RAISE NOTICE 'Audit system is now fully set up for existing and future tables';
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.setup_complete_audit_system() IS 'Function to set up the complete audit system for existing and future tables';

-- Step 10: Function to get basic audit statistics
CREATE OR REPLACE FUNCTION audit.get_audit_statistics()
RETURNS TABLE(
    schema_name TEXT,
    table_name TEXT,
    inserts BIGINT,
    updates BIGINT,
    deletes BIGINT,
    truncates BIGINT,
    total_actions BIGINT,
    last_action_time TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        la.schema_name,
        la.table_name,
        SUM(CASE WHEN la.action = 'I' THEN 1 ELSE 0 END)::BIGINT AS inserts,
        SUM(CASE WHEN la.action = 'U' THEN 1 ELSE 0 END)::BIGINT AS updates,
        SUM(CASE WHEN la.action = 'D' THEN 1 ELSE 0 END)::BIGINT AS deletes,
        SUM(CASE WHEN la.action = 'T' THEN 1 ELSE 0 END)::BIGINT AS truncates,
        COUNT(*)::BIGINT AS total_actions,
        MAX(la.action_tstamp) AS last_action_time
    FROM audit.logged_actions la
    GROUP BY la.schema_name, la.table_name
    ORDER BY la.schema_name, la.table_name;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.get_audit_statistics() IS 'Function to get basic statistics about audit records';

-- Step 11: Function to clean up old audit records based on retention period
CREATE OR REPLACE FUNCTION audit.cleanup_old_audit_records(retention_days INTEGER DEFAULT 90)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM audit.logged_actions
    WHERE action_tstamp < (CURRENT_TIMESTAMP - (retention_days || ' days')::INTERVAL);
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    RAISE NOTICE 'Deleted % audit records older than % days', deleted_count, retention_days;
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION audit.cleanup_old_audit_records(INTEGER) IS 'Function to clean up old audit records based on retention period';

-- Step 12: Simple view to see recent audit activity
CREATE OR REPLACE VIEW audit.recent_activity AS
SELECT
    event_id,
    action_tstamp,
    CASE
        WHEN action = 'I' THEN 'INSERT'
        WHEN action = 'U' THEN 'UPDATE'
        WHEN action = 'D' THEN 'DELETE'
        WHEN action = 'T' THEN 'TRUNCATE'
    END AS action_type,
    schema_name || '.' || table_name AS table_path,
    table_pk,
    current_user_name,
    client_addr
FROM
    audit.logged_actions
ORDER BY
    action_tstamp DESC
LIMIT 100;

COMMENT ON VIEW audit.recent_activity IS 'View to see recent activity in the audit log';
SELECT audit.setup_complete_audit_system();

-- Final Step: Initialize everything with a single command
-- Uncomment and run this line to set up the entire audit system