import { Login, Logout, Menu } from '@mui/icons-material'
import { Toolbar, IconButton, Typography, Button, styled, Avatar } from '@mui/material'
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';

import React from 'react'
import { User } from '../reducer/user/userSlice';
import { DRAWER_WIDTH } from './AppDrawer';

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

interface Props {
  user: User | null
  logoutHandler: () => void
  isOpen: boolean
  toggleDrawer: (isOpen: boolean) => void
}

const AppBarComponent = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    marginLeft: DRAWER_WIDTH,
    width: `calc(100% - ${DRAWER_WIDTH}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const AppBar: React.FC<Props> = ({ user, logoutHandler, isOpen, toggleDrawer }) => {
  return (
    <AppBarComponent position="absolute" open={isOpen}>
      <Toolbar
        sx={{
          pr: '24px', // keep right padding when drawer closed
        }}>
        <IconButton edge="start" color="inherit" aria-label="open drawer" onClick={() => toggleDrawer(isOpen)}
          sx={{
            marginRight: '36px',
            ...(isOpen && { display: 'none' }),
          }}>
          <Menu />
        </IconButton>
        <Typography align="left" component="h1" variant="h5" color="inherit" noWrap sx={{ flexGrow: 1 }}>
          Git Noter
        </Typography>
        {user == null ?
          <Button variant="contained" color="secondary" href="/api/v1/oauth2/login/github" endIcon={<Login />}>
            Login
          </Button>
          :
          <>
            <Avatar sx={{ mx: 2 }} alt={user.name} src={user.avatar_url}></Avatar>
            <Button variant="contained" color="secondary" onClick={() => logoutHandler()} endIcon={<Logout />}>
              Logout
            </Button>
          </>
        }
      </Toolbar>
    </AppBarComponent>
  )
}

export default AppBar
