BEGIN;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.updated_at = now(); 
      RETURN NEW;
   ELSE
      RETURN OLD;
   END IF;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_accounts_updated_at_time BEFORE UPDATE ON accounts FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_roles_updated_at_time BEFORE UPDATE ON roles FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_account_events_updated_at_time BEFORE UPDATE ON account_events FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_addresses_updated_at_time BEFORE UPDATE ON addresses FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_sessions_updated_at_time BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_phone_numbers_updated_at_time BEFORE UPDATE ON phone_numbers FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_equipment_updated_at_time BEFORE UPDATE ON equipment FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_authorizations_updated_at_time BEFORE UPDATE ON authorizations FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER update_account_authorizations_updated_at_time BEFORE UPDATE ON account_authorizations FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

COMMIT;