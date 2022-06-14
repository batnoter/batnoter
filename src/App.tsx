import { ThemeProvider } from '@emotion/react';
import { Box, createTheme, CssBaseline, useMediaQuery } from '@mui/material';
import ModalProvider from 'mui-modal-provider';
import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import './App.scss';
import Main from './components/Main';
import { useAppDispatch, useAppSelector } from './app/hooks';
import { selectAppTheme } from './reducer/preferenceSlice';
import { setAppTheme } from './reducer/preferenceSlice';

const App: React.FC = () => {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

  const appTheme = useAppSelector(selectAppTheme);
  const dispatch = useAppDispatch();

  React.useEffect(() => {
    dispatch(setAppTheme(prefersDarkMode ? 'dark' : 'light'));
  }, [prefersDarkMode]);

  return (
    <div className="App">
      <BrowserRouter basename={process.env.REACT_APP_BASENAME}>
        <ThemeProvider theme={createTheme({
          palette: { mode: appTheme === 'dark' ? 'dark' : 'light' }
        })}>
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
