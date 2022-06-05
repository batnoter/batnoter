
import { CircularProgress, Toolbar } from '@mui/material';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import React, { ReactElement, useEffect, useState } from 'react';
import { Outlet, Route, Routes } from 'react-router-dom';
import { API_URL } from '../api/api';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { APIStatusType } from '../reducer/common';
import { getNotesAsync, getNotesTreeAsync } from '../reducer/noteSlice';
import { getUserProfileAsync, selectUser, selectUserAPIStatus, User, userLoading, userLogout } from '../reducer/userSlice';
import ErrorPage from "./404";
import AppBar from './AppBar';
import AppDrawer from './AppDrawer';
import Editor from './Editor';
import Finder from './Finder';
import RequireAuth from './lib/RequireAuth';
import Login from './Login';
import RepoSetupDialog from './RepoSetupDialog';
import Settings from './Settings';
import Viewer from './Viewer';
const DrawerLayout: React.FC<{ user: User | null }> = ({ user }): ReactElement => {
  return (
    <Box sx={{ display: 'flex', flexGrow: 1 }}>
      <AppDrawer user={user} />
      <Box component="main" sx={{
        flexGrow: 1, height: '100vh', overflow: 'auto'
      }}>
        <Toolbar variant="dense" />
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
          <Outlet />
        </Container>
      </Box>
    </Box>
  );
}

const isUserAPILoading = (userAPIStatus: APIStatusType): boolean => {
  return userAPIStatus === APIStatusType.LOADING;
}

const Main: React.FC = (): ReactElement => {
  const dispatch = useAppDispatch();
  const user = useAppSelector(selectUser);
  const userAPIStatus = useAppSelector(selectUserAPIStatus);
  const [apiTriggered, setAPITriggered] = useState(false);

  const handleLogin = () => {
    dispatch(userLoading());
    window.location.href = API_URL + "/oauth2/login/github";
  }

  const handleLogout = () => {
    dispatch(userLogout());
  }

  useEffect(() => {
    dispatch(getUserProfileAsync());
    setAPITriggered(true);
  }, [])

  useEffect(() => {
    (async () => {
      if (userAPIStatus == APIStatusType.IDLE && user != null) {
        await dispatch(getNotesTreeAsync());
        dispatch(getNotesAsync(""));
      }
    })()
  }, [userAPIStatus, user]);

  return (
    <>
      <AppBar userAPIStatus={userAPIStatus} handleLogin={handleLogin} handleLogout={handleLogout} user={user} />
      <Container maxWidth="xl">
        {user != null && !user?.default_repo?.name && <RepoSetupDialog open={true}></RepoSetupDialog>}
        {
          !apiTriggered || isUserAPILoading(userAPIStatus) ? <CircularProgress color="inherit" sx={{ ml: '50%', mt: 10 }} /> :
            <Routes>
              <Route path="/login" element={<Login userAPIStatus={userAPIStatus} handleLogin={handleLogin} user={user} />} />
              <Route path="/" element={<DrawerLayout user={user} />} >
                <Route index element={<RequireAuth user={user}><Finder /></RequireAuth>} />
                <Route path="/new" element={<RequireAuth user={user}><Editor key={'new'} /></RequireAuth>} />
                <Route path="/edit" element={<RequireAuth user={user}><Editor key={'edit'} /></RequireAuth>} />
                <Route path="/view" element={<RequireAuth user={user}><Viewer key={'view'} /></RequireAuth>} />
                <Route path="/settings" element={<RequireAuth user={user}><Settings user={user} /></RequireAuth>} />
              </Route>
              <Route path="*" element={<ErrorPage />} />
            </Routes>
        }
      </Container>
    </>
  );
}

export default Main;
