import { Box, Typography, Alert, Paper } from '@mui/material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism';

const Query = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Query Command
      </Typography>
      
      <Typography variant="body1" paragraph>
        After connecting to a database, you can execute SQL queries and manage transactions. The Query functionality 
        supports both single-query execution and transaction-based operations.
      </Typography>

      <Alert severity="info" sx={{ mb: 3 }}>
        <Typography variant="body1">
          You must be connected to a database using the Connect command before executing queries.
        </Typography>
      </Alert>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Single Query Execution
      </Typography>

      <Typography variant="body1" paragraph>
        Execute individual SQL statements directly:
      </Typography>

      <Paper sx={{ p: 2, mb: 4 }}>
        <SyntaxHighlighter language="sql" style={atomDark}>
          {"SELECT * FROM table_name;\nINSERT INTO table_name VALUES (1, 'value');\nUPDATE table_name SET column = 'new_value' WHERE id = 1;\nDELETE FROM table_name WHERE condition;"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Transaction Management
      </Typography>

      <Typography variant="body1" paragraph>
        For operations that require atomicity and consistency, use transactions:
      </Typography>

      <Box sx={{ mb: 4 }}>
        <Typography variant="subtitle1" gutterBottom>
          1. Start a Transaction
        </Typography>
        <Paper sx={{ p: 2, mb: 3 }}>
          <SyntaxHighlighter language="sql" style={atomDark}>
            {"start transaction"}
          </SyntaxHighlighter>
        </Paper>

        <Typography variant="subtitle1" gutterBottom>
          2. Execute Queries
        </Typography>
        <Paper sx={{ p: 2, mb: 3 }}>
          <SyntaxHighlighter language="sql" style={atomDark}>
            {"INSERT INTO users (id, name) VALUES (1, 'John');\nUPDATE accounts SET balance = balance - 100 WHERE user_id = 1;\nINSERT INTO transactions (user_id, amount) VALUES (1, -100);"}
          </SyntaxHighlighter>
        </Paper>

        <Typography variant="subtitle1" gutterBottom>
          3. Commit or Rollback
        </Typography>
        <Paper sx={{ p: 2, mb: 3 }}>
          <SyntaxHighlighter language="sql" style={atomDark}>
            {"-- To save changes:\ncommit\n\n-- To discard changes:\nrollback"}
          </SyntaxHighlighter>
        </Paper>
      </Box>

      <Typography variant="h5" gutterBottom>
        Complete Transaction Example
      </Typography>

      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="sql" style={atomDark}>
          {"start transaction\n-- Create a new user\nINSERT INTO users VALUES (1, 'John');\n\n-- Update user's name\nUPDATE users SET name = 'Johnny' WHERE id = 1;\n\n-- Save all changes\ncommit"}
        </SyntaxHighlighter>
      </Paper>

      <Alert severity="warning" sx={{ mt: 4 }}>
        <Typography variant="body1">
          <strong>Important:</strong> Always ensure you have the necessary permissions to execute queries on the tables. 
          Use transactions when performing multiple related operations to maintain data consistency.
        </Typography>
      </Alert>

      <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
        Best Practices
      </Typography>
      
      <Box sx={{ pl: 2 }}>
        <Typography variant="body1" paragraph>
          • Use transactions for operations that must be executed together
        </Typography>
        <Typography variant="body1" paragraph>
          • Always commit or rollback your transactions
        </Typography>
        <Typography variant="body1" paragraph>
          • Keep transactions as short as possible
        </Typography>
        <Typography variant="body1" paragraph>
          • Test queries with SELECT before performing updates or deletes
        </Typography>
        <Typography variant="body1" paragraph>
          • Use appropriate WHERE clauses to avoid unintended modifications
        </Typography>
      </Box>
    </Box>
  );
};

export default Query;