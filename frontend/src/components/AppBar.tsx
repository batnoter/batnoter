import { Login as LoginIcon } from '@mui/icons-material';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import BugReportIcon from '@mui/icons-material/BugReport';
import TwitterIcon from '@mui/icons-material/Twitter';
import { Avatar, Button, CircularProgress, Link, Menu, MenuItem, SvgIconTypeMap, Toolbar, Typography } from '@mui/material';
import AppBarComponent from '@mui/material/AppBar';
import { OverridableComponent } from '@mui/material/OverridableComponent';
import { HelpCircle } from 'mdi-material-ui';
import React, { ReactElement } from 'react';
import { NavLink} from 'react-router-dom';
import { APIStatusType } from '../reducer/common';
import { User } from '../reducer/userSlice';
import { URL_FAQ, URL_ISSUES, URL_TWITTER_HANDLE } from '../util/util';

interface Props {
  user: User | null
  userAPIStatus: APIStatusType
  handleLogin: () => void
  handleLogout: () => void
}


const getExternalLink = (url: string, label: string, Icon: OverridableComponent<SvgIconTypeMap>): ReactElement => {
  return <Link href={url} sx={{ mx: 1, color: 'inherit' }} target="_blank" rel="noopener">
    <Icon sx={{ mx: 0.5, verticalAlign: 'middle' }} fontSize="inherit" />{label}
  </Link>
}

const isLoading = (apiStatus: APIStatusType): boolean => {
  return apiStatus === APIStatusType.LOADING;
}

const AppBar: React.FC<Props> = ({ user, userAPIStatus, handleLogin, handleLogout }): ReactElement => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <AppBarComponent position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
      <Toolbar variant="dense" sx={{ justifyContent: "space-between" }}>
        <Link variant="h6" noWrap component={NavLink} to={"/"} sx={{ flexGrow: 1, display: "flex", color: 'inherit' }}>GITNOTER</Link>
        {getExternalLink(URL_TWITTER_HANDLE, "@gitnoter", TwitterIcon)}
        {getExternalLink(URL_FAQ, "faq", HelpCircle)}
        {getExternalLink(URL_ISSUES, "bug report", BugReportIcon)}
        {user == null ?
          (
            !isLoading(userAPIStatus) ? <Button color="inherit" endIcon={<LoginIcon />} onClick={() => handleLogin()}>Login</Button>
              : <CircularProgress color="inherit" />
          )
          :
          <>
            <Link component={NavLink} to={"/new"} sx={{ mx: 1, color: 'inherit' }}>
              <AddCircleIcon sx={{ mx: 0.5, verticalAlign: 'middle' }} fontSize="inherit" />create note
            </Link>

            <Avatar onClick={handleMenu} alt={user.name} src={user.avatar_url} sx={{ "cursor": "pointer" }}></Avatar>
            <Menu autoFocus={false} sx={{ mt: '5px' }} id="menu-appbar" anchorEl={anchorEl} anchorOrigin={{
              vertical: 'bottom', horizontal: 'right'
            }} transformOrigin={{ vertical: 'top', horizontal: 'right', }} open={Boolean(anchorEl)} onClose={handleClose}>
              <MenuItem component={NavLink} to={"/settings"} onClick={handleClose}>Setting</MenuItem>
              <MenuItem onClick={handleLogout}>Logout</MenuItem>
            </Menu>
          </>
        }
      </Toolbar>
    </AppBarComponent>
  )
}

export default AppBar;
