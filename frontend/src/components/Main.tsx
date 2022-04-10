
import { Toolbar } from '@mui/material';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import CssBaseline from '@mui/material/CssBaseline';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import ModalProvider from 'mui-modal-provider';
import React, { ReactElement, useEffect } from 'react';
import { Route, Routes } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { APIStatusType } from '../reducer/common';
import { getNotesAsync, getNotesTreeAsync } from '../reducer/noteSlice';
import { getUserProfileAsync, selectUser, selectUserStatus, userLoading, userLogout } from '../reducer/userSlice';
import AppBar from './AppBar';
import AppDrawer from './AppDrawer';
import Editor from './Editor';
import Finder from './Finder';
import RepoSetupDialog from './RepoSetupDialog';
import Settings from './Settings';
import Viewer from './Viewer';

const Main: React.FC = (): ReactElement => {
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
    (async () => {
      if (userStatus == APIStatusType.IDLE && user != null) {
        await dispatch(getNotesTreeAsync())
        dispatch(getNotesAsync(""))
      }
    })()
  }, [userStatus, user])

  return (
    <ThemeProvider theme={createTheme()}>
      <ModalProvider>
        <Box sx={{ display: 'flex' }}>
          <CssBaseline />
          <AppBar userStatus={userStatus} setUserStatus={setUserStatus} handleLogout={handleLogout} user={user} />
          <AppDrawer user={user} />
          <Box component="main" sx={{
            backgroundColor: (theme) => theme.palette.mode === 'light'
              ? theme.palette.grey[100] : theme.palette.grey[900], flexGrow: 1, height: '100vh', overflow: 'auto',
          }}>
            <Toolbar variant="dense" />
            {user != null && !user?.default_repo?.name && <RepoSetupDialog open={true}></RepoSetupDialog>}
            <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
              <Routes>
                <Route path="/" element={<Finder />} ></Route>
                <Route path="/new" element={<Editor key={'new'} />} ></Route>
                <Route path="/edit" element={<Editor key="edit" />} ></Route>
                <Route path="/view" element={<Viewer key={'view'} />} ></Route>
                <Route path="/settings" element={<Settings user={user} />} ></Route>
              </Routes>
            </Container>
          </Box>
        </Box>
      </ModalProvider>
    </ThemeProvider>
  );
}

export default Main;
