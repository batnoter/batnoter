import { Login as LoginIcon, Menu as MenuIcon } from '@mui/icons-material';
import { Avatar, Box, Button, CircularProgress, IconButton, Menu, MenuItem, styled, Toolbar } from '@mui/material';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import React from 'react';
import { User, UserStatus } from '../reducer/user/userSlice';
import { DRAWER_WIDTH } from './AppDrawer';


interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

interface Props {
  user: User | null
  userStatus: UserStatus
  setUserStatus: (userStatus: UserStatus) => void
  handleLogout: () => void
  isOpen: boolean
  toggleDrawer: (isOpen: boolean) => void
}

const AppBarComponent = styled(MuiAppBar, { shouldForwardProp: (prop) => prop !== 'open' })<AppBarProps>(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer + 1, transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }), ...(open && {
    marginLeft: DRAWER_WIDTH, width: `calc(100% - ${DRAWER_WIDTH}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp, duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const AppBar: React.FC<Props> = ({ user, userStatus, setUserStatus, handleLogout, isOpen, toggleDrawer }) => {
  const isLoading = userStatus === UserStatus.LOADING
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (<AppBarComponent position="absolute" open={isOpen}>
    <Toolbar sx={{ pr: '24px', /* keep right padding when drawer closed */ }}>
      <IconButton edge="start" color="inherit" aria-label="open drawer" onClick={() => toggleDrawer(isOpen)}
        sx={{ marginRight: '36px', ...(isOpen && { display: 'none' }), }}>
        <MenuIcon />
      </IconButton>
      <Box sx={{ flexGrow: 1 }}></Box> {user == null ?
        (!isLoading ? <Button color="inherit" href="/api/v1/oauth2/login/github" endIcon={<LoginIcon />}
          onClick={() => setUserStatus(UserStatus.LOADING)}>Login</Button> :
          <CircularProgress color="inherit" />)
        :
        <>
          <Avatar onClick={handleMenu} alt={user.name} src={user.avatar_url} sx={{ "cursor": "pointer" }}></Avatar>
          <Menu autoFocus={false} sx={{ mt: '45px' }} id="menu-appbar" anchorEl={anchorEl} anchorOrigin={{
            vertical: 'top', horizontal: 'right'
          }} transformOrigin={{ vertical: 'top', horizontal: 'right', }} open={Boolean(anchorEl)} onClose={handleClose}>
            <MenuItem onClick={handleClose}>Profile</MenuItem>
            <MenuItem onClick={handleLogout}>Logout</MenuItem>
          </Menu>
        </>
      }
    </Toolbar>
  </AppBarComponent>
  )
}

export default AppBar
