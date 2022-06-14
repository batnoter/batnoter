import { ThemeProvider } from '@emotion/react';
import { Box, createTheme, CssBaseline, useMediaQuery } from '@mui/material';
import ModalProvider from 'mui-modal-provider';
import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import './App.scss';
import Main from './components/Main';
import { useAppDispatch, useAppSelector } from './app/hooks';
import { selectThemeMode } from './reducer/preferenceSlice';
import { setThemeMode } from './reducer/preferenceSlice';

const App: React.FC = () => {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

  const themeMode = useAppSelector(selectThemeMode);
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    dispatch(setThemeMode(prefersDarkMode ? 'dark' : 'light'));
  }, [prefersDarkMode]);

  const theme = createTheme({
    palette: { mode: themeMode === 'dark' ? 'dark' : 'light' }
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
