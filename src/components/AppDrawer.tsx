import AddCircleIcon from '@mui/icons-material/AddCircle';
import HelpCenterIcon from '@mui/icons-material/HelpCenter';
import NotesIcon from '@mui/icons-material/Notes';
import SettingsIcon from '@mui/icons-material/Settings';
import StarIcon from '@mui/icons-material/Star';
import { Divider, List, ListItemButton, ListItemIcon, ListItemText, styled, Toolbar } from '@mui/material';
import MuiDrawer from '@mui/material/Drawer';
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { User } from '../reducer/userSlice';

interface Props {
  user: User | null
}

export const DRAWER_WIDTH = 240;
const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
  ({ theme }) => ({
    '& .MuiDrawer-paper': {
      boxSizing: 'border-box',
      position: 'relative',
      whiteSpace: 'nowrap',
      width: DRAWER_WIDTH,
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen,
      }),
      [theme.breakpoints.down('sm')]: {
        width: theme.spacing(9),
      },
    },
  }),
);

const AppDrawer: React.FC<Props> = ({ user }) => {
  const { pathname } = useLocation();
  return (
    <Drawer variant="permanent">
      <Toolbar />
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
