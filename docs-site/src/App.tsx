import { useState } from 'react';
import Home from './pages/Home';
import Features from './pages/Features';
import Security from './pages/Security';
import { Routes, Route, Link, useLocation } from 'react-router-dom';
import Login from './components/commands/Login';
import Create from './components/commands/Create';
import Connect from './components/commands/Connect';
import Grant from './components/commands/Grant';
import Query from './components/commands/Query';
import LinkMigrate from './components/commands/LinkMigrate';
import { 
  AppBar, Toolbar, Typography, Drawer, List, ListItem, ListItemText, 
  ListItemIcon, Box, Container, IconButton, useMediaQuery, Divider,
  useTheme, CssBaseline
} from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import HomeIcon from '@mui/icons-material/Home';
import CodeIcon from '@mui/icons-material/Code';
import SecurityIcon from '@mui/icons-material/Security';
import StorageIcon from '@mui/icons-material/Storage';
import LoginIcon from '@mui/icons-material/Login';
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import LinkIcon from '@mui/icons-material/Link';
import VpnKeyIcon from '@mui/icons-material/VpnKey';
import './App.css'

const drawerWidth = 240;

function App() {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const [mobileOpen, setMobileOpen] = useState(false);
  const location = useLocation();

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const menuItems = [
    { text: 'Home', icon: <HomeIcon />, path: '/' },
    { text: 'Features', icon: <CodeIcon />, path: '/features' },
    { text: 'Commands', icon: null, path: '', divider: true, header: true },
    { text: 'Login', icon: <LoginIcon />, path: '/commands/login' },
    { text: 'Create', icon: <PersonAddIcon />, path: '/commands/create' },
    { text: 'Connect', icon: <StorageIcon />, path: '/commands/connect' },
    { text: 'Grant', icon: <VpnKeyIcon />, path: '/commands/grant' },
    { text: 'Query', icon: <CodeIcon />, path: '/commands/query' },
    { text: 'Link & Migrate', icon: <LinkIcon />, path: '/commands/link-migrate' },
    { text: 'Security', icon: <SecurityIcon />, path: '/security' },
  ];

  const drawer = (
    <div>
      <Toolbar>
        <Typography variant="h6" noWrap component="div" sx={{ fontWeight: 'bold' }}>
          TCP-DuckDB
        </Typography>
      </Toolbar>
      <Divider />
      <List>
        {menuItems.map((item, index) => (
          item.header ? (
            <Box key={`header-${index}`} sx={{ mt: 2, mb: 1 }}>
              {item.divider && <Divider />}
              <Typography 
                variant="overline" 
                sx={{ 
                  pl: 2, 
                  display: 'block', 
                  color: 'primary.main',
                  fontWeight: 'bold'
                }}
              >
                {item.text}
              </Typography>
            </Box>
          ) : (
            <ListItem 
              button 
              key={item.text} 
              component={Link} 
              to={item.path}
              selected={location.pathname === item.path}
              onClick={isMobile ? handleDrawerToggle : undefined}
              sx={{
                '&.Mui-selected': {
                  backgroundColor: 'rgba(255, 193, 7, 0.1)',
                  borderLeft: '4px solid',
                  borderColor: 'primary.main',
                  '& .MuiListItemIcon-root': {
                    color: 'primary.main',
                  },
                },
                '&:hover': {
                  backgroundColor: 'rgba(255, 193, 7, 0.05)',
                },
              }}
            >
              {item.icon && <ListItemIcon>{item.icon}</ListItemIcon>}
              <ListItemText primary={item.text} />
            </ListItem>
          )
        ))}
      </List>
    </div>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar 
        position="fixed" 
        sx={{
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
          boxShadow: 'none',
          borderBottom: '1px solid rgba(255, 255, 255, 0.12)'
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div">
            TCP Server for DuckDB Documentation
          </Typography>
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile.
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{ 
          flexGrow: 1, 
          p: 3, 
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          minHeight: '100vh',
          backgroundColor: 'background.default'
        }}
      >
        <Toolbar />
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/features" element={<Features />} />
            <Route path="/commands/login" element={<Login />} />
            <Route path="/commands/create" element={<Create />} />
            <Route path="/commands/connect" element={<Connect />} />
            <Route path="/commands/grant" element={<Grant />} />
            <Route path="/commands/query" element={<Query />} />
            <Route path="/commands/link-migrate" element={<LinkMigrate />} />
            <Route path="/security" element={<Security />} />
          </Routes>
        </Container>
      </Box>
    </Box>
  )
}

export default App
