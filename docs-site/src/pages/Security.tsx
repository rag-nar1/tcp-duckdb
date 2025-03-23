import { Box, Typography, Alert, Paper, Grid } from '@mui/material';

const Security = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Security and Access Control
      </Typography>
      
      <Typography variant="body1" paragraph>
        The TCP Server for DuckDB enforces strict security measures and access control mechanisms to ensure data protection and proper authorization.
      </Typography>

      <Alert severity="warning" sx={{ mb: 4 }}>
        <Typography variant="body1">
          <strong>IMPORTANT:</strong> The super user is 'duck', which has privileges to create databases and users. 
          The default password is 'duck' - it is crucial to change this password immediately after setting up your project for security purposes.
        </Typography>
      </Alert>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Authentication System
      </Typography>

      <Paper sx={{ p: 3, mb: 4, border: '1px solid', borderColor: 'divider' }}>
        <Typography variant="body1" paragraph>
          The server implements a username and password-based authentication system. All operations require proper authentication, and the server maintains session information to track authenticated users.
        </Typography>

        <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold', mt: 2 }}>
          Key Authentication Features:
        </Typography>

        <Box sx={{ pl: 2 }}>
          <Typography variant="body1" paragraph>
            • Mandatory login before performing any operations
          </Typography>
          <Typography variant="body1" paragraph>
            • Password verification for all authentication attempts
          </Typography>
          <Typography variant="body1" paragraph>
            • Session tracking for authenticated users
          </Typography>
        </Box>
      </Paper>

      <Typography variant="h5" gutterBottom>
        Access Control Mechanisms
      </Typography>

      <Paper sx={{ p: 3, mb: 4, border: '1px solid', borderColor: 'divider' }}>
        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>
              User Privileges
            </Typography>
            <Box sx={{ pl: 2 }}>
              <Typography variant="body1" paragraph>
                • <strong>Super User:</strong> Only the super user 'duck' can create databases and users or grant permissions
              </Typography>
              <Typography variant="body1" paragraph>
                • <strong>Regular Users:</strong> Can only perform operations they have been explicitly granted access to
              </Typography>
            </Box>
          </Grid>

          <Grid item xs={12} md={6}>
            <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>
              Permission Levels
            </Typography>
            <Box sx={{ pl: 2 }}>
              <Typography variant="body1" paragraph>
                • <strong>Database-level:</strong> Read or write access to entire databases
              </Typography>
              <Typography variant="body1" paragraph>
                • <strong>Table-level:</strong> Granular permissions for select, insert, update, and delete operations
              </Typography>
            </Box>
          </Grid>
        </Grid>

        <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold', mt: 2 }}>
          Permission Enforcement:
        </Typography>
        <Box sx={{ pl: 2 }}>
          <Typography variant="body1" paragraph>
            • Permissions are checked for every connection and query operation
          </Typography>
          <Typography variant="body1" paragraph>
            • Users can only connect to databases they have been granted access to
          </Typography>
          <Typography variant="body1" paragraph>
            • Table operations are restricted based on granted permissions
          </Typography>
        </Box>
      </Paper>

      <Typography variant="h5" gutterBottom>
        Security Best Practices
      </Typography>

      <Paper sx={{ p: 3, mb: 4, border: '1px solid', borderColor: 'divider' }}>
        <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>
          Recommended Security Measures:
        </Typography>
        <Box sx={{ pl: 2 }}>
          <Typography variant="body1" paragraph>
            • Change the default super user password immediately after installation
          </Typography>
          <Typography variant="body1" paragraph>
            • Use strong, complex passwords for all user accounts
          </Typography>
          <Typography variant="body1" paragraph>
            • Grant minimal necessary permissions to users (principle of least privilege)
          </Typography>
          <Typography variant="body1" paragraph>
            • Regularly review and audit user permissions
          </Typography>
          <Typography variant="body1" paragraph>
            • Secure the network environment where the TCP server is deployed
          </Typography>
          <Typography variant="body1" paragraph>
            • Consider using encrypted connections for sensitive data
          </Typography>
        </Box>
      </Paper>

      <Alert severity="info" sx={{ mt: 4 }}>
        <Typography variant="body1">
          For more detailed information about specific commands related to security, please refer to the 
          <strong> Login</strong>, <strong>Create</strong>, and <strong>Grant</strong> command documentation.
        </Typography>
      </Alert>
    </Box>
  );
};

export default Security;