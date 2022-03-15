import { Bookmark, ChevronLeft, HelpCenter, Pageview, Settings } from '@mui/icons-material';
import { Divider, IconButton, List, ListItemButton, ListItemIcon, ListItemText, styled, Toolbar } from '@mui/material';
import MuiDrawer from '@mui/material/Drawer';
import React from 'react';


interface Props {
    isOpen: boolean,
    toggleDrawer : Function
}
const drawerWidth: number = 240;
const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
    ({ theme, open }) => ({
      '& .MuiDrawer-paper': {
        position: 'relative',
        whiteSpace: 'nowrap',
        width: drawerWidth,
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
  
const AppDrawer: React.FC<Props> = ({isOpen, toggleDrawer}) => {
    return (
        <Drawer variant="permanent" open={isOpen}>
        <Toolbar
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'flex-end',
            px: [1],
          }}
        >
          <IconButton onClick={()=>toggleDrawer(isOpen)}>
            <ChevronLeft />
          </IconButton>
        </Toolbar>
        <Divider />
        <List component="nav">
            <ListItemButton>
              <ListItemIcon>
                <Pageview />
              </ListItemIcon>
              <ListItemText primary="Search" />
            </ListItemButton>
            <ListItemButton>
              <ListItemIcon>
                <Bookmark />
              </ListItemIcon>
              <ListItemText primary="Bookmarks" />
            </ListItemButton>
            <ListItemButton>
              <ListItemIcon>
                <Settings />
              </ListItemIcon>
              <ListItemText primary="Settings" />
            </ListItemButton>
            <ListItemButton>
              <ListItemIcon>
                <HelpCenter />
              </ListItemIcon>
              <ListItemText primary="Help" />
            </ListItemButton>
          <Divider sx={{ my: 1 }} />
        </List>
      </Drawer>
    )
}

export default AppDrawer
