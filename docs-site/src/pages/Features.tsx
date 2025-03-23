import { Box, Typography, Paper, Grid, Card, CardContent } from '@mui/material';

const Features = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Features
      </Typography>
      
      <Typography variant="body1" paragraph>
        The TCP Server for DuckDB provides a comprehensive set of features for managing and interacting with DuckDB databases over a network.
        Here's an overview of the key capabilities:
      </Typography>

      <Grid container spacing={4} sx={{ mb: 6 }}>
        {[
          {
            title: 'User Authentication and Authorization',
            description: 'Secure login system with user-based access control. The server enforces authentication for all operations and provides a robust permission system.',
            icon: '🔐'
          },
          {
            title: 'Database Creation and Management',
            description: 'Create and manage DuckDB databases with ease. Super users can create databases and grant access to other users.',
            icon: '💾'
          },
          {
            title: 'Table-level Access Control',
            description: 'Fine-grained permissions system that allows controlling access at the table level with specific operation types (select, insert, update, delete).',
            icon: '🛡️'
          },
          {
            title: 'PostgreSQL Database Linking',
            description: 'Connect and synchronize with PostgreSQL databases. Copy schemas and data from PostgreSQL to DuckDB and keep them in sync.',
            icon: '🔄'
          },
          {
            title: 'Transaction Support',
            description: 'Full support for database transactions, including start transaction, commit, and rollback operations to ensure data consistency.',
            icon: '📊'
          },
          {
            title: 'Query Execution',
            description: 'Execute SQL queries with proper access control. The server checks permissions for every query operation.',
            icon: '⚡'
          }
        ].map((feature, index) => (
          <Grid item xs={12} sm={6} md={4} key={index}>
            <Card 
              sx={{ 
                height: '100%', 
                display: 'flex', 
                flexDirection: 'column',
                transition: 'transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out',
                border: '1px solid',
                borderColor: 'divider',
                '&:hover': {
                  transform: 'translateY(-5px)',
                  boxShadow: 3,
                  borderColor: 'primary.main',
                }
              }}
            >
              <CardContent sx={{ flexGrow: 1 }}>
                <Typography variant="h2" component="div" sx={{ fontSize: '2.5rem', mb: 2, color: 'text.secondary' }}>
                  {feature.icon}
                </Typography>
                <Typography variant="h6" component="h3" gutterBottom sx={{ fontWeight: 'bold' }}>
                  {feature.title}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {feature.description}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Technical Specifications
      </Typography>

      <Paper sx={{ p: 3, mb: 4, border: '1px solid', borderColor: 'divider' }}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>
              Server Architecture
            </Typography>
            <Typography variant="body2" paragraph>
              • TCP-based network protocol
            </Typography>
            <Typography variant="body2" paragraph>
              • Multi-user support with concurrent connections
            </Typography>
            <Typography variant="body2" paragraph>
              • Persistent storage of user credentials and permissions
            </Typography>
          </Grid>
          <Grid item xs={12} md={6}>
            <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>
              Database Features
            </Typography>
            <Typography variant="body2" paragraph>
              • Full SQL support via DuckDB
            </Typography>
            <Typography variant="body2" paragraph>
              • ACID-compliant transactions
            </Typography>
            <Typography variant="body2" paragraph>
              • Integration with PostgreSQL databases
            </Typography>
          </Grid>
        </Grid>
      </Paper>

      <Typography variant="h5" gutterBottom>
        Security Features
      </Typography>

      <Paper sx={{ p: 3, mb: 4, border: '1px solid', borderColor: 'divider' }}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>
              Authentication
            </Typography>
            <Typography variant="body2" paragraph>
              • Username and password-based authentication
            </Typography>
            <Typography variant="body2" paragraph>
              • Super user privileges for administrative tasks
            </Typography>
            <Typography variant="body2" paragraph>
              • Session-based authentication
            </Typography>
          </Grid>
          <Grid item xs={12} md={6}>
            <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>
              Authorization
            </Typography>
            <Typography variant="body2" paragraph>
              • Database-level read/write permissions
            </Typography>
            <Typography variant="body2" paragraph>
              • Table-level operation permissions (select, insert, update, delete)
            </Typography>
            <Typography variant="body2" paragraph>
              • Permission checking for every database operation
            </Typography>
          </Grid>
        </Grid>
      </Paper>
    </Box>
  );
};

export default Features;