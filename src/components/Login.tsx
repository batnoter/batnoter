import GitHubIcon from '@mui/icons-material/GitHub';
import { LoadingButton } from '@mui/lab';
import { Box, Container, Toolbar, Typography } from '@mui/material';
import React, { ReactElement, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { APIStatusType } from '../reducer/common';
import { User } from '../reducer/userSlice';

interface Props {
  user: User | null
  userAPIStatus: APIStatusType
  handleLogin: () => void
}

const isLoading = (apiStatus: APIStatusType, user: User | null): boolean => {
  return apiStatus === APIStatusType.LOADING || user != null;
}

const Login: React.FC<Props> = ({ user, handleLogin, userAPIStatus }): ReactElement => {
  const navigate = useNavigate();


  useEffect(() => {
    user != null && navigate("/", { replace: true });
  }, [user]);

  return (
    <Container maxWidth="xl">
      <Toolbar variant="dense" />

      <Box display="flex" sx={{ my: 2 }} alignItems="center" justifyContent={'space-around'}>
        <Box flexGrow={1} display={{ xs: "none", md: "block" }}>
          {/* TODO: Show animated GIF to showcase app features*/}
        </Box>
        <Box sx={{ my: 6, p: 2, width: '400px', height: '100%', border: '1px solid grey', borderRadius: 2 }}>
          <Typography variant="h5" align="center">GET STARTED</Typography>
          <p>Welcome to BatNoter &#127881;. Please login with your github account to start using the application</p>
          <LoadingButton onClick={() => handleLogin()}
            loading={isLoading(userAPIStatus, user)} fullWidth sx={{ my: 2 }}
            variant="contained" startIcon={<GitHubIcon />}>Login with Github</LoadingButton>
        </Box>
      </Box>
    </Container>
  )
}

export default Login;
