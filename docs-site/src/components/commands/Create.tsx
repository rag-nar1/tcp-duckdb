import { Box, Typography, Alert, Paper } from '@mui/material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism';

const Create = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Create Command
      </Typography>
      
      <Typography variant="body1" paragraph>
        The Create command is used to create databases and users in the TCP-DuckDB server. This command requires super user privileges.
      </Typography>

      <Alert severity="info" sx={{ mb: 3 }}>
        <Typography variant="body1">
          Only the super user (duck) has the privileges to create databases and users.
        </Typography>
      </Alert>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Create Database
      </Typography>

      <Typography variant="h6" gutterBottom>
        Syntax
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"create database [database_name]"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h6" gutterBottom>
        Parameters
      </Typography>
      
      <Box sx={{ pl: 2, mb: 3 }}>
        <Typography variant="body1" component="div">
          <strong>[database_name]</strong>: The name of the database to create
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom>
        Example
      </Typography>
      
      <Paper sx={{ p: 2, mb: 4 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"create database mydb"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Create User
      </Typography>

      <Typography variant="h6" gutterBottom>
        Syntax
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"create user [username] [password]"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h6" gutterBottom>
        Parameters
      </Typography>
      
      <Box sx={{ pl: 2, mb: 3 }}>
        <Typography variant="body1" component="div">
          <strong>[username]</strong>: The username for the new user
        </Typography>
        <Typography variant="body1" component="div">
          <strong>[password]</strong>: The password for the new user
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom>
        Example
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"create user john pass123"}
        </SyntaxHighlighter>
      </Paper>

      <Alert severity="warning" sx={{ mt: 4 }}>
        <Typography variant="body1">
          <strong>Security Note:</strong> Always use strong passwords when creating new users. Passwords should be complex 
          and not easily guessable.
        </Typography>
      </Alert>
    </Box>
  );
};

export default Create;