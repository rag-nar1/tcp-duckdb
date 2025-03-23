import { createTheme } from '@mui/material/styles';

// DuckDB color scheme
const duckdbYellow = '#FFC107'; // Primary color - DuckDB yellow
const duckdbDarkGray = '#212121'; // Secondary color - dark gray
const duckdbMediumGray = '#424242'; // Medium gray for accents

// Create a theme instance with DuckDB colors
const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: duckdbYellow,
      contrastText: '#000',
    },
    secondary: {
      main: duckdbMediumGray,
    },
    background: {
      default: duckdbDarkGray,
      paper: '#1A1A1A',
    },
    text: {
      primary: '#fff',
      secondary: 'rgba(255, 255, 255, 0.7)',
    },
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontWeight: 700,
    },
    h2: {
      fontWeight: 700,
    },
    h3: {
      fontWeight: 600,
    },
    h4: {
      fontWeight: 600,
    },
    h5: {
      fontWeight: 500,
    },
    h6: {
      fontWeight: 500,
    },
  },
  components: {
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: '#1A1A1A',
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          backgroundColor: '#1A1A1A',
          borderRight: '1px solid rgba(255, 255, 255, 0.12)',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          textTransform: 'none',
          fontWeight: 600,
        },
        containedPrimary: {
          '&:hover': {
            backgroundColor: '#e6ad00',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          backgroundColor: '#212121',
          borderRadius: 12,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
        },
      },
    },
  },
});

export default theme;