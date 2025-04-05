# Railway Deployment Instructions

## Setup

1. Install the Railway CLI: 
   ```
   npm i -g @railway/cli
   ```

2. Login to Railway:
   ```
   railway login
   ```

3. Link your project:
   ```
   railway link
   ```

4. Deploy:
   ```
   railway up
   ```

## Environment Variables

Make sure these environment variables are set in your Railway project:

- `ServerPort`: 4000
- `ServerAddr`: 0.0.0.0
- `DBdir`: /app/storge/
- `ServerDbFile`: server/db.sqlite3
- `ENCRYPTION_KEY`: A15pG0m3hwf0tfpVW6m92eZ6vRmAQA3C

## Database

The SQLite database will be automatically created at `/app/storge/server/db.sqlite3` during the container build and initialized with the schema defined in `scheme.sql`.

The initialization creates the following tables:
- `user` - User accounts and authentication
- `DB` - Database definitions
- `tables` - Table definitions
- `dbprivilege` - Database access privileges
- `tableprivilege` - Table-specific access privileges
- `postgres` - PostgreSQL connection configurations

**Note**: Since Railway containers are ephemeral, any data stored in the SQLite database will be lost when the container restarts. For production use, consider using a persistent database service. 