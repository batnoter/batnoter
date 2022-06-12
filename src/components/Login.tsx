import GitHubIcon from '@mui/icons-material/GitHub';
import { LoadingButton } from '@mui/lab';
import { Box, Container, Toolbar, Typography } from '@mui/material';
import React, { ReactElement, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { getToken } from '../api/api';
import { useAppDispatch } from '../app/hooks';
import { APIStatusType } from '../reducer/common';
import { getUserProfileAsync, User } from '../reducer/userSlice';

interface Props {
  user: User | null
  userAPIStatus: APIStatusType
  handleLogin: () => void
}

const isLoading = (apiStatus: APIStatusType, user: User | null): boolean => {
  return apiStatus === APIStatusType.LOADING || user != null;
}

const Login: React.FC<Props> = ({ user, handleLogin, userAPIStatus }): ReactElement => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const loginSuccess = searchParams.get('success') === "true";

  useEffect(() => {
    user != null && navigate("/", { replace: true });
    if (loginSuccess) {
      // user has just completed the oauth login
      // call getToken api to get the app token and store it in localStorage
      getToken().then(() => {
        dispatch(getUserProfileAsync());
        navigate("/", { replace: true });
      })
    }
  }, [user, loginSuccess]);

  return (
    <Container maxWidth="xl">
      <Toolbar variant="dense" />

      <Box display="flex" sx={{ my: 2 }} alignItems="center" justifyContent={'space-around'}>
        <Box flexGrow={1} sx={{ mx: 0, my: 2 }} display={{ xs: "none", md: "block" }}>
          <img style={{ width: '100%', border: "1px solid #80808080", borderRadius: "8px" }} src="/demo.gif" />
        </Box>
        <Box flexShrink={0} sx={{ my: 6, ml: 4, p: 2, width: '400px', height: '100%', border: '1px solid grey', borderRadius: 2 }}>
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
