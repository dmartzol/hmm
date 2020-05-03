BEGIN;

CREATE OR REPLACE FUNCTION update_update_time_column()
RETURNS TRIGGER AS $$
BEGIN
   IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.update_time = now(); 
      RETURN NEW;
   ELSE
      RETURN OLD;
   END IF;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_accounts_update_time BEFORE UPDATE ON accounts FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();
CREATE TRIGGER update_roles_update_time BEFORE UPDATE ON roles FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();
CREATE TRIGGER update_account_events_update_time BEFORE UPDATE ON account_events FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();
CREATE TRIGGER update_addresses_update_time BEFORE UPDATE ON addresses FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();
CREATE TRIGGER update_sessions_update_time BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();
CREATE TRIGGER update_equipment_update_time BEFORE UPDATE ON equipment FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();
CREATE TRIGGER update_authorizations_update_time BEFORE UPDATE ON authorizations FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();
CREATE TRIGGER update_account_authorizations_update_time BEFORE UPDATE ON account_authorizations FOR EACH ROW EXECUTE PROCEDURE update_update_time_column();

COMMIT;