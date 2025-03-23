import { Box, Typography, Alert, Paper } from '@mui/material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism';

const Grant = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Grant Command
      </Typography>
      
      <Typography variant="body1" paragraph>
        The Grant command is used to manage access permissions for users to databases and tables. This command requires super user privileges.
      </Typography>

      <Alert severity="info" sx={{ mb: 3 }}>
        <Typography variant="body1">
          Only the super user (duck) can grant permissions to other users.
        </Typography>
      </Alert>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Database Access
      </Typography>

      <Typography variant="h6" gutterBottom>
        Syntax
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"grant database [database_name] [username] [access_type]"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h6" gutterBottom>
        Parameters
      </Typography>
      
      <Box sx={{ pl: 2, mb: 3 }}>
        <Typography variant="body1" component="div">
          <strong>[database_name]</strong>: The name of the database to grant access to
        </Typography>
        <Typography variant="body1" component="div">
          <strong>[username]</strong>: The user to grant access to
        </Typography>
        <Typography variant="body1" component="div">
          <strong>[access_type]</strong>: Type of access to grant (read, write)
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom>
        Example
      </Typography>
      
      <Paper sx={{ p: 2, mb: 4 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"grant database mydb john read"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Table Access
      </Typography>

      <Typography variant="h6" gutterBottom>
        Syntax
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"grant table [database_name] [table_name] [username] [access_type...]"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h6" gutterBottom>
        Parameters
      </Typography>
      
      <Box sx={{ pl: 2, mb: 3 }}>
        <Typography variant="body1" component="div">
          <strong>[database_name]</strong>: The name of the database containing the table
        </Typography>
        <Typography variant="body1" component="div">
          <strong>[table_name]</strong>: The name of the table to grant access to
        </Typography>
        <Typography variant="body1" component="div">
          <strong>[username]</strong>: The user to grant access to
        </Typography>
        <Typography variant="body1" component="div">
          <strong>[access_type...]</strong>: Types of access to grant (select, update, insert, delete)
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom>
        Example
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"grant table mydb users john select insert"}
        </SyntaxHighlighter>
      </Paper>

      <Alert severity="info" sx={{ mt: 4 }}>
        <Typography variant="body1">
          <strong>Note:</strong> Table-level permissions provide fine-grained control over what operations users can perform. 
          Users need appropriate database access before table permissions can be granted.
        </Typography>
      </Alert>

      <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
        Access Types
      </Typography>

      <Typography variant="subtitle1" gutterBottom>
        Database Level:
      </Typography>
      <Box sx={{ pl: 2, mb: 2 }}>
        <Typography variant="body1" paragraph>
          • <strong>read</strong>: Allows connecting to the database and reading data
        </Typography>
        <Typography variant="body1" paragraph>
          • <strong>write</strong>: Allows all read operations plus data modification
        </Typography>
      </Box>

      <Typography variant="subtitle1" gutterBottom>
        Table Level:
      </Typography>
      <Box sx={{ pl: 2 }}>
        <Typography variant="body1" paragraph>
          • <strong>select</strong>: Allows reading data from the table
        </Typography>
        <Typography variant="body1" paragraph>
          • <strong>insert</strong>: Allows adding new records to the table
        </Typography>
        <Typography variant="body1" paragraph>
          • <strong>update</strong>: Allows modifying existing records in the table
        </Typography>
        <Typography variant="body1" paragraph>
          • <strong>delete</strong>: Allows removing records from the table
        </Typography>
      </Box>
    </Box>
  );
};

export default Grant;