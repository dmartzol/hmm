# Intallation 

## For first time dabase setup run these:

`psql -tc "SELECT 1 FROM pg_database WHERE datname = 'hmm'" | grep -q 1 || psql -c "CREATE DATABASE hmm"`

`psql -f create.sql`


## Next times you just need to run:

`psql -f rebuild.sql --quiet`
