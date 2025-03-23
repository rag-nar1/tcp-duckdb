import { Box, Typography, Alert, Paper } from '@mui/material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism';

const Connect = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Connect Command
      </Typography>
      
      <Typography variant="body1" paragraph>
        The Connect command is used to establish a connection to a specific database for executing queries. 
        Once connected, you can perform various database operations including running queries and managing transactions.
      </Typography>

      <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
        Syntax
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"connect [database_name]"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h6" gutterBottom>
        Parameters
      </Typography>
      
      <Box sx={{ pl: 2, mb: 3 }}>
        <Typography variant="body1" component="div">
          <strong>[database_name]</strong>: The name of the database you want to connect to
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom>
        Example
      </Typography>
      
      <Paper sx={{ p: 2, mb: 4 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"connect mydb"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        After Connecting
      </Typography>

      <Typography variant="body1" paragraph>
        Once connected to a database, you can:
      </Typography>

      <Box sx={{ pl: 2, mb: 4 }}>
        <Typography variant="body1" paragraph>
          • Execute single queries
        </Typography>
        <Typography variant="body1" paragraph>
          • Start transactions
        </Typography>
        <Typography variant="body1" paragraph>
          • Commit or rollback changes
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom>
        Transaction Commands
      </Typography>

      <Box sx={{ mb: 4 }}>
        <Typography variant="body2" sx={{ mb: 2 }}>
          Start a transaction:
        </Typography>
        <Paper sx={{ p: 2, mb: 3 }}>
          <SyntaxHighlighter language="sql" style={atomDark}>
            {"start transaction"}
          </SyntaxHighlighter>
        </Paper>

        <Typography variant="body2" sx={{ mb: 2 }}>
          Commit changes:
        </Typography>
        <Paper sx={{ p: 2, mb: 3 }}>
          <SyntaxHighlighter language="sql" style={atomDark}>
            {"commit"}
          </SyntaxHighlighter>
        </Paper>

        <Typography variant="body2" sx={{ mb: 2 }}>
          Rollback changes:
        </Typography>
        <Paper sx={{ p: 2, mb: 3 }}>
          <SyntaxHighlighter language="sql" style={atomDark}>
            {"rollback"}
          </SyntaxHighlighter>
        </Paper>
      </Box>

      <Typography variant="h6" gutterBottom>
        Transaction Example
      </Typography>

      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="sql" style={atomDark}>
          {"start transaction\nINSERT INTO users VALUES (1, 'John');\nUPDATE users SET name = 'Johnny' WHERE id = 1;\ncommit"}
        </SyntaxHighlighter>
      </Paper>

      <Alert severity="info" sx={{ mt: 4 }}>
        <Typography variant="body1">
          <strong>Note:</strong> You must have appropriate permissions to connect to a database and execute queries. 
          These permissions are managed through the Grant command.
        </Typography>
      </Alert>
    </Box>
  );
};

export default Connect;