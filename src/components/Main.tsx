import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import CssBaseline from '@mui/material/CssBaseline';
import Link from '@mui/material/Link';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import Typography from '@mui/material/Typography';
import * as React from 'react';
import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { getUserProfileAsync, logout, selectUser } from '../reducer/user/userSlice';
import AppBar from './AppBar';
import AppDrawer from './AppDrawer';

function MainContent() {
  const dispatch = useAppDispatch();
  useEffect(() => {
    dispatch(getUserProfileAsync())
  }, [])

  const [open, setOpen] = React.useState(true);
  const toggleDrawer = (open: boolean) => {
    setOpen(!open);
  };
  const logoutHandler = () => {
    dispatch(logout())
  }
  const user = useAppSelector(selectUser);

  return (
    <ThemeProvider theme={createTheme()}>
      <Box sx={{ display: 'flex' }}>
        <CssBaseline />
        <AppBar logoutHandler={logoutHandler} user={user} isOpen={open} toggleDrawer={toggleDrawer} />
        <AppDrawer user={user} isOpen={open} toggleDrawer={toggleDrawer} />
        <Box
          component="main" sx={{
            backgroundColor: (theme) => theme.palette.mode === 'light'
              ? theme.palette.grey[100] : theme.palette.grey[900],
            flexGrow: 1, height: '100vh', overflow: 'auto',
          }}
        >
          <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Typography variant="body2" color="text.secondary" align="center" sx={{ pt: 4 }}>
              {'Copyright © '} <Link color="inherit" href="https://batnoter.com/"> batnoter.com </Link>
              {' '}{new Date().getFullYear()} {'.'}
            </Typography>
          </Container>
        </Box>
      </Box>
    </ThemeProvider>
  );
}

export default function Main() {
  return <MainContent />;
}
