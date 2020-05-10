# Intallation 

## First time dabase setup. Run these:

psql -tc "SELECT 1 FROM pg_database WHERE datname = 'hackerspace'" | grep -q 1 || psql -c "CREATE DATABASE hackerspace"
psql -f create.sql


## Next times you just need to run:

psql -f rebuild.sql --quiet