import { Box, Typography, Alert, Paper } from '@mui/material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism';

const LinkMigrate = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Link and Migrate Commands
      </Typography>
      
      <Typography variant="body1" paragraph>
        The Link and Migrate commands enable integration between DuckDB and PostgreSQL databases. These commands 
        facilitate data synchronization and are particularly useful for maintaining consistency across different 
        database systems.
      </Typography>

      <Alert severity="info" sx={{ mb: 3 }}>
        <Typography variant="body1">
          These commands require super user privileges and are essential for PostgreSQL database integration.
        </Typography>
      </Alert>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Link Command
      </Typography>

      <Typography variant="body1" paragraph>
        The Link command establishes a connection between DuckDB and PostgreSQL by:
        <ul>
          <li>Reading PostgreSQL table schemas</li>
          <li>Recreating the schemas in DuckDB</li>
          <li>Copying all data from PostgreSQL tables into DuckDB</li>
        </ul>
      </Typography>

      <Typography variant="h6" gutterBottom>
        Syntax
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"link [database_name] [postgresql_connection_string]"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h6" gutterBottom>
        Parameters
      </Typography>
      
      <Box sx={{ pl: 2, mb: 3 }}>
        <Typography variant="body1" component="div">
          <strong>[database_name]</strong>: The name of the DuckDB database to link
        </Typography>
        <Typography variant="body1" component="div">
          <strong>[postgresql_connection_string]</strong>: The PostgreSQL connection string
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom>
        Example
      </Typography>
      
      <Paper sx={{ p: 2, mb: 4 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"link mydb \"postgresql://user:password@localhost:5432/pgdb\""}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Migrate Command
      </Typography>

      <Typography variant="body1" paragraph>
        The Migrate command maintains synchronization between DuckDB and PostgreSQL databases by reading the audit 
        table to keep the DuckDB database in sync with PostgreSQL changes.
      </Typography>

      <Typography variant="h6" gutterBottom>
        Syntax
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"migrate [database_name]"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h6" gutterBottom>
        Parameters
      </Typography>
      
      <Box sx={{ pl: 2, mb: 3 }}>
        <Typography variant="body1" component="div">
          <strong>[database_name]</strong>: The name of the DuckDB database to synchronize
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom>
        Example
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"migrate mydb"}
        </SyntaxHighlighter>
      </Paper>

      <Alert severity="warning" sx={{ mt: 4 }}>
        <Typography variant="body1">
          <strong>Important:</strong> Ensure that the PostgreSQL database has the necessary audit tables and triggers 
          set up for proper change tracking. The migration process relies on these audit mechanisms to maintain data 
          consistency.
        </Typography>
      </Alert>

      <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
        Best Practices
      </Typography>
      
      <Box sx={{ pl: 2 }}>
        <Typography variant="body1" paragraph>
          • Regularly run the migrate command to keep databases synchronized
        </Typography>
        <Typography variant="body1" paragraph>
          • Monitor the audit tables for any potential synchronization issues
        </Typography>
        <Typography variant="body1" paragraph>
          • Ensure proper network connectivity between DuckDB and PostgreSQL servers
        </Typography>
        <Typography variant="body1" paragraph>
          • Use secure connection strings and protect database credentials
        </Typography>
      </Box>
    </Box>
  );
};

export default LinkMigrate;