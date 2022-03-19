import AddCircleIcon from '@mui/icons-material/AddCircle';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import HelpCenterIcon from '@mui/icons-material/HelpCenter';
import NotesIcon from '@mui/icons-material/Notes';
import SettingsIcon from '@mui/icons-material/Settings';
import StarIcon from '@mui/icons-material/Star';
import { Divider, IconButton, List, ListItemButton, ListItemIcon, ListItemText, styled, Toolbar, Typography } from '@mui/material';
import MuiDrawer from '@mui/material/Drawer';
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { User } from '../reducer/user/userSlice';


interface Props {
  user: User | null
  isOpen: boolean,
  toggleDrawer: (isOpen: boolean) => void
}

export const DRAWER_WIDTH = 240;
const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
  ({ theme, open }) => ({
    '& .MuiDrawer-paper': {
      position: 'relative',
      whiteSpace: 'nowrap',
      width: DRAWER_WIDTH,
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen,
      }),
      boxSizing: 'border-box',
      ...(!open && {
        overflowX: 'hidden',
        transition: theme.transitions.create('width', {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.leavingScreen,
        }),
        width: theme.spacing(7),
        [theme.breakpoints.up('sm')]: {
          width: theme.spacing(9),
        },
      }),
    },
  }),
);

const AppDrawer: React.FC<Props> = ({ isOpen, toggleDrawer }) => {
  const { pathname } = useLocation();
  return (
    <Drawer variant="permanent" open={isOpen}>
      <Toolbar sx={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-end', px: [1], }} >
        <Typography component="h1" variant="h5" color="inherit" noWrap sx={{ flexGrow: 1 }}>
          Git Noter
        </Typography>
        <IconButton onClick={() => toggleDrawer(isOpen)}> <ChevronLeftIcon /> </IconButton>
      </Toolbar>
      <Divider />
      <List component="nav">
        <ListItemButton component={Link} to={"/"} selected={pathname === '/'}> <ListItemIcon> <NotesIcon /> </ListItemIcon><ListItemText primary="My Notes" /></ListItemButton>
        <ListItemButton component={Link} to={"/new"} selected={pathname === '/new'}> <ListItemIcon> <AddCircleIcon /> </ListItemIcon><ListItemText primary="Create Note" /></ListItemButton>
        <ListItemButton component={Link} to={"/favorites"} selected={pathname === '/favorites'}> <ListItemIcon> <StarIcon /> </ListItemIcon> <ListItemText primary="Favorites" /> </ListItemButton>
        <ListItemButton component={Link} to={"/settings"} selected={pathname === '/settings'}> <ListItemIcon> <SettingsIcon /> </ListItemIcon><ListItemText primary="Settings" /></ListItemButton>
        <ListItemButton component={Link} to={"/help"} selected={pathname === '/help'}> <ListItemIcon> <HelpCenterIcon /> </ListItemIcon> <ListItemText primary="Help" /> </ListItemButton>
        <Divider sx={{ my: 1 }} />
      </List>
    </Drawer>
  )
}

export default AppDrawer
