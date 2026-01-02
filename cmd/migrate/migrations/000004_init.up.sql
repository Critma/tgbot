CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER reminder_update_trigger BEFORE UPDATE ON reminders FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();