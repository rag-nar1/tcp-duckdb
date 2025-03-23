import { Box, Typography, Paper, Grid, Button, Card, CardContent, CardMedia } from '@mui/material';
import { Link } from 'react-router-dom';

const Home = () => {
  return (
    <Box>
      {/* Hero Section */}
      <Paper 
        elevation={0}
        sx={{
          p: 6,
          borderRadius: 2,
          backgroundImage: 'linear-gradient(to right, #1A1A1A, #212121)',
          mb: 6,
          position: 'relative',
          overflow: 'hidden',
          border: '1px solid',
          borderColor: 'rgba(255, 193, 7, 0.2)'
        }}
      >
        <Box sx={{ position: 'relative', zIndex: 1 }}>
          <Typography variant="h3" component="h1" gutterBottom sx={{ fontWeight: 'bold', color: 'primary.main' }}>
            TCP Server for DuckDB
          </Typography>
          <Typography variant="h6" paragraph sx={{ maxWidth: '800px', mb: 4 }}>
            A TCP server that enables management and interaction with DuckDB databases over a network.
            The server provides functionality for user authentication, database operations, and access control.
          </Typography>
          <Box sx={{ mt: 4 }}>
            <Button 
              component={Link} 
              to="/features" 
              variant="contained" 
              color="primary"
              size="large"
              sx={{ mr: 2, fontWeight: 'bold' }}
            >
              Explore Features
            </Button>
            <Button 
              component={Link} 
              to="/commands/login" 
              variant="outlined" 
              color="primary"
              size="large"
              sx={{ fontWeight: 'bold' }}
            >
              View Commands
            </Button>
          </Box>
        </Box>
      </Paper>

      {/* Features Section */}
      <Typography variant="h4" component="h2" gutterBottom sx={{ mb: 4, fontWeight: 'bold' }}>
        Key Features
      </Typography>
      
      <Grid container spacing={4} sx={{ mb: 6 }}>
        {[
          {
            title: 'User Authentication',
            description: 'Secure login system with user-based access control',
            icon: '🔐'
          },
          {
            title: 'Database Management',
            description: 'Create and manage DuckDB databases with ease',
            icon: '💾'
          },
          {
            title: 'Access Control',
            description: 'Fine-grained table-level permissions system',
            icon: '🛡️'
          },
          {
            title: 'PostgreSQL Linking',
            description: 'Connect and sync with PostgreSQL databases',
            icon: '🔄'
          },
          {
            title: 'Transaction Support',
            description: 'Full support for database transactions',
            icon: '📊'
          },
          {
            title: 'Query Execution',
            description: 'Execute SQL queries with proper access control',
            icon: '⚡'
          }
        ].map((feature, index) => (
          <Grid item xs={12} sm={6} md={4} key={index}>
            <Card 
              sx={{ 
                height: '100%', 
                display: 'flex', 
                flexDirection: 'column',
                transition: 'transform 0.3s ease-in-out, box-shadow 0.3s ease-in-out, border 0.3s ease-in-out',
                border: '1px solid transparent',
                '&:hover': {
                  transform: 'translateY(-5px)',
                  boxShadow: '0 8px 16px rgba(255, 193, 7, 0.15)',
                  border: '1px solid',
                  borderColor: 'primary.main'
                }
              }}
            >
              <CardContent sx={{ flexGrow: 1 }}>
                <Typography variant="h1" component="div" sx={{ fontSize: '3rem', mb: 2 }}>
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

      {/* Getting Started Section */}
      <Paper 
        elevation={0} 
        sx={{ 
          p: 4, 
          borderRadius: 2, 
          bgcolor: 'background.paper',
          border: '1px solid',
          borderColor: 'rgba(255, 193, 7, 0.3)',
          backgroundImage: 'linear-gradient(to bottom right, rgba(255, 193, 7, 0.05), transparent)'
        }}
      >
        <Typography variant="h5" component="h2" gutterBottom sx={{ fontWeight: 'bold' }}>
          Ready to Get Started?
        </Typography>
        <Typography paragraph>
          Check out our comprehensive documentation to learn how to use the TCP Server for DuckDB.
        </Typography>
        <Button 
          component={Link} 
          to="/commands/login" 
          variant="contained" 
          color="primary"
          sx={{ fontWeight: 'bold' }}
        >
          View Command Reference
        </Button>
      </Paper>
    </Box>
  );
};

export default Home;