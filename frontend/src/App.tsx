import { ThemeProvider } from '@emotion/react';
import { Box, createTheme, CssBaseline } from '@mui/material';
import ModalProvider from 'mui-modal-provider';
import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import './App.scss';
import Main from './components/Main';

const App: React.FC = () => {
  return (
    <div className="App">
      <BrowserRouter>
        <ThemeProvider theme={createTheme()}>
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
