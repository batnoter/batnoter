import { Login as LoginIcon } from '@mui/icons-material';
import ThemeToggleIconDark from '@mui/icons-material/DarkMode';
import FavoriteIcon from '@mui/icons-material/Favorite';
import GitHubIcon from '@mui/icons-material/GitHub';
import ThemeToggleIconLight from '@mui/icons-material/LightMode';
import TwitterIcon from '@mui/icons-material/Twitter';
import { Avatar, Box, Button, CircularProgress, Link, LinkProps, LinkTypeMap, Menu, MenuItem, SvgIconTypeMap, Toolbar } from '@mui/material';
import AppBarComponent from '@mui/material/AppBar';
import { OverridableComponent } from '@mui/material/OverridableComponent';
import { Ladybug, MessageQuestion, PlusBox } from 'mdi-material-ui';
import React, { ReactElement } from 'react';
import { NavLink } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { RootState } from '../app/store';
import { APIStatusType } from '../reducer/common';
import { setThemeMode } from '../reducer/preferenceSlice';
import { User } from '../reducer/userSlice';
import { URL_FAQ, URL_GITHUB, URL_ISSUES, URL_SPONSOR, URL_TWITTER_HANDLE } from '../util/util';

interface Props {
  user: User | null
  userAPIStatus: APIStatusType
  handleLogin: () => void
  handleLogout: () => void
}

export type AppBarLinkProps = {
  label: string
  icon: OverridableComponent<SvgIconTypeMap>
  iconColor?: string
}

const AppBarLink = <D extends React.ElementType = LinkTypeMap["defaultComponent"], P = AppBarLinkProps>
  ({ label, icon: Icon, iconColor, children, ...rest }: LinkProps<D, P> & AppBarLinkProps) =>
  <Link {...rest} sx={{
    p: 0.2, mx: 0.5, borderRadius: '50%', color: 'inherit', display: 'flex',
    bgcolor: { xs: 'action.disabled', lg: 'unset' }
  }} {...(rest.href ? { target: "_blank", rel: "noopener" } : {})}>
    <Icon sx={{ m: 0.5, verticalAlign: 'middle', color: iconColor }} fontSize="inherit" />
    <Box sx={{ display: { xs: 'none', lg: 'block' } }}>{label}</Box>
    {children}
  </Link>

const isLoading = (apiStatus: APIStatusType): boolean => {
  return apiStatus === APIStatusType.LOADING;
}

const AppBar: React.FC<Props> = ({ user, userAPIStatus, handleLogin, handleLogout }): ReactElement => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const dispatch = useAppDispatch();
  const themeMode = useAppSelector((state: RootState) => state.preference.themeMode);


  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleThemeModeToggle = () => {
    if (themeMode === 'light') { dispatch(setThemeMode('dark')) }
    else if (themeMode === 'dark') { dispatch(setThemeMode('light')) }
  }

  return (
    <AppBarComponent position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
      <Toolbar variant="dense" sx={{ justifyContent: "space-between" }}>
        <Link variant="h6" noWrap component={NavLink} to={"/"} sx={{ flexGrow: 1, display: "flex", color: 'inherit' }}>BATNOTER</Link>
        <Button sx={{ mx: 1, color: 'inherit' }} onClick={handleThemeModeToggle}>
          {themeMode === 'dark' ? <ThemeToggleIconLight /> : <ThemeToggleIconDark />}
        </Button>
        <AppBarLink href={URL_SPONSOR} label="sponsor" icon={FavoriteIcon} iconColor="#d489b5" />
        <AppBarLink href={URL_TWITTER_HANDLE} label="@batnoter" icon={TwitterIcon} iconColor="#b1d5ff" />
        <AppBarLink href={URL_FAQ} label="faq" icon={MessageQuestion} iconColor="#c7d097" />
        <AppBarLink href={URL_ISSUES} label="bug report" icon={Ladybug} iconColor="#eeb082" />
        <AppBarLink href={URL_GITHUB} label="github" icon={GitHubIcon} iconColor="#dadada" />
        {user && <AppBarLink component={NavLink} to="/new" label="create note" icon={PlusBox} iconColor="#c1f497" />}

        <Box sx={{ ml: 1 }}>
          {user == null ?
            (isLoading(userAPIStatus) ? <CircularProgress color="inherit" /> :
              <Button color="inherit" endIcon={<LoginIcon />} onClick={() => handleLogin()}>Login</Button>)
            :
            <>
              <Avatar onClick={handleMenu} alt={user.name} src={user.avatar_url} sx={{ cursor: "pointer" }} />
              <Menu autoFocus={false} sx={{ mt: '5px' }} id="menu-appbar" anchorEl={anchorEl} anchorOrigin={{
                vertical: 'bottom', horizontal: 'right'
              }} transformOrigin={{ vertical: 'top', horizontal: 'right', }} open={Boolean(anchorEl)} onClose={handleClose}>
                <MenuItem component={NavLink} to="/settings" onClick={handleClose}>Setting</MenuItem>
                <MenuItem onClick={handleLogout}>Logout</MenuItem>
              </Menu>
            </>
          }
        </Box>
      </Toolbar>
    </AppBarComponent>
  )
}

export default AppBar;
