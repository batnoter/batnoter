
import { Toolbar } from '@mui/material';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import CssBaseline from '@mui/material/CssBaseline';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import * as React from 'react';
import { useEffect } from 'react';
import { Route, Routes } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { searchNotesAsync } from '../reducer/noteSlice';
import { getUserProfileAsync, selectUser, selectUserStatus, userLoading, userLogout, UserStatus } from '../reducer/userSlice';
import AppBar from './AppBar';
import AppDrawer from './AppDrawer';
import Editor from './Editor';
import { Favorites } from './Favorites';
import Finder from './Finder';
import RepoSelectDialog from './RepoSelectDialog';
import Settings from './Settings';

const Main = () => {
  const dispatch = useAppDispatch();
  useEffect(() => {
    dispatch(getUserProfileAsync())
  }, [])
  const user = useAppSelector(selectUser);
  const userStatus = useAppSelector(selectUserStatus);

  const setUserStatus = () => {
    dispatch(userLoading())
  }
  const handleLogout = () => {
    dispatch(userLogout())
  }

  useEffect(() => {
    if (userStatus == UserStatus.IDLE && user != null) {
      dispatch(searchNotesAsync())
    }
  }, [userStatus, user])

  return (
    <ThemeProvider theme={createTheme()}>
      <Box sx={{ display: 'flex' }}>
        <CssBaseline />
        <AppBar userStatus={userStatus} setUserStatus={setUserStatus} handleLogout={handleLogout} user={user} />
        <AppDrawer user={user} />
        <Box component="main" sx={{
          backgroundColor: (theme) => theme.palette.mode === 'light'
            ? theme.palette.grey[100] : theme.palette.grey[900], flexGrow: 1, height: '100vh', overflow: 'auto',
        }}>
          <Toolbar />
          <RepoSelectDialog open={user != null && !user?.default_repo?.name}></RepoSelectDialog>
          <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Routes>
              <Route path="/" element={<Finder />} ></Route>
              <Route path="/new" element={<Editor />} ></Route>
              <Route path="/favorites" element={<Favorites />} ></Route>
              <Route path="/settings" element={<Settings user={user} />} ></Route>
            </Routes>
          </Container>
        </Box>
      </Box>
    </ThemeProvider>
  );
}

export default Main
