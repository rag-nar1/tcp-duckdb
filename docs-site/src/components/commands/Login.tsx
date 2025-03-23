import { Box, Typography, Alert, Paper } from '@mui/material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomDark } from 'react-syntax-highlighter/dist/esm/styles/prism';

const Login = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Login Command
      </Typography>
      
      <Typography variant="body1" paragraph>
        The Login command is used to authenticate users to access the TCP-DuckDB server. Authentication is required before performing any operations on the server.
      </Typography>

      <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
        Syntax
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"login [username] [password]"}
        </SyntaxHighlighter>
      </Paper>

      <Typography variant="h6" gutterBottom>
        Parameters
      </Typography>
      
      <Box sx={{ pl: 2 }}>
        <Typography variant="body1" component="div">
          <strong>[username]</strong>: Your username for authentication
        </Typography>
        <Typography variant="body1" component="div">
          <strong>[password]</strong>: Your password for authentication
        </Typography>
      </Box>

      <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
        Example
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <SyntaxHighlighter language="bash" style={atomDark}>
          {"login duck superpassword"}
        </SyntaxHighlighter>
      </Paper>

      <Alert severity="warning" sx={{ mt: 4, mb: 3 }}>
        <Typography variant="body1">
          <strong>IMPORTANT:</strong> The super user is 'duck', which has privileges to create databases and users. 
          The default password is 'duck' - it is crucial to change this password immediately after setting up your project for security purposes.
        </Typography>
      </Alert>

      <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
        Security Considerations
      </Typography>
      
      <Box sx={{ pl: 2 }}>
        <Typography variant="body1" paragraph>
          • Always use strong passwords and change the default super user password immediately
        </Typography>
        <Typography variant="body1" paragraph>
          • Keep your credentials secure and never share them
        </Typography>
        <Typography variant="body1" paragraph>
          • The server enforces authentication for all operations
        </Typography>
      </Box>
    </Box>
  );
};

export default Login;