import { ThemeProvider } from '@emotion/react';
import { Box, createTheme, CssBaseline, useMediaQuery } from '@mui/material';
import ModalProvider from 'mui-modal-provider';
import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import './App.scss';
import Main from './components/Main';

const App: React.FC = () => {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

  const theme = createTheme({
    palette: {
      mode: prefersDarkMode ? 'dark' : 'light',
    },
  });

  return (
    <div className="App">
      <BrowserRouter basename={process.env.REACT_APP_BASENAME}>
        <ThemeProvider theme={theme}>
          <ModalProvider>
            <Box sx={{ display: 'flex' }}>
              <CssBaseline />
              <Main />
            </Box>
          </ModalProvider>
        </ThemeProvider>
      </BrowserRouter>
    </div>
  );
}

export default App;
