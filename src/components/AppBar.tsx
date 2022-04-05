import { Login as LoginIcon } from '@mui/icons-material';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import { Avatar, Box, Button, CircularProgress, Menu, MenuItem, Toolbar, Typography } from '@mui/material';
import AppBarComponent from '@mui/material/AppBar';
import React, { ReactElement } from 'react';
import { Link } from 'react-router-dom';
import { User, UserStatus } from '../reducer/userSlice';

interface Props {
  user: User | null
  userStatus: UserStatus
  setUserStatus: (userStatus: UserStatus) => void
  handleLogout: () => void
}

const AppBar: React.FC<Props> = ({ user, userStatus, setUserStatus, handleLogout }): ReactElement => {
  const isLoading = userStatus === UserStatus.LOADING
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <AppBarComponent position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
      <Toolbar variant="dense">
        <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1, display: "flex" }}>
          GIT NOTER
        </Typography>
        <Box sx={{ flexGrow: 1 }}></Box> {user == null ?
          (
            !isLoading ?
              <Button color="inherit" href="/api/v1/oauth2/login/github" endIcon={<LoginIcon />}
                onClick={() => setUserStatus(UserStatus.LOADING)}>Login</Button>
              : <CircularProgress color="inherit" />
          )
          :
          <>
            <Button color="inherit" component={Link} to={"/new"} startIcon={<AddCircleIcon />}
              onClick={() => setUserStatus(UserStatus.LOADING)}>Create Note</Button>
            <Avatar onClick={handleMenu} alt={user.name} src={user.avatar_url} sx={{ "cursor": "pointer" }}></Avatar>
            <Menu autoFocus={false} sx={{ mt: '5px' }} id="menu-appbar" anchorEl={anchorEl} anchorOrigin={{
              vertical: 'bottom', horizontal: 'right'
            }} transformOrigin={{ vertical: 'top', horizontal: 'right', }} open={Boolean(anchorEl)} onClose={handleClose}>
              <MenuItem component={Link} to={"/settings"} onClick={handleClose}>Setting</MenuItem>
              <MenuItem onClick={handleLogout}>Logout</MenuItem>
            </Menu>
          </>
        }
      </Toolbar>
    </AppBarComponent>
  )
}

export default AppBar;
