import { Login, Menu } from '@mui/icons-material'
import { Toolbar, IconButton, Typography, Button, styled } from '@mui/material'
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';

import React from 'react'

const drawerWidth: number = 240;

interface AppBarProps extends MuiAppBarProps {
    open?: boolean;
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
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
        }),
    }),
}));

interface Props {
    isOpen: boolean,
    toggleDrawer: Function
}

const AppBar: React.FC<Props> = ({ isOpen, toggleDrawer }) => {
    return (
        <AppBarComponent position="absolute" open={isOpen}>
            <Toolbar
                sx={{
                    pr: '24px', // keep right padding when drawer closed
                }}>
                <IconButton
                    edge="start"
                    color="inherit"
                    aria-label="open drawer"
                    onClick={() => toggleDrawer(isOpen)}
                    sx={{
                        marginRight: '36px',
                        ...(isOpen && { display: 'none' }),
                    }}>
                    <Menu />
                </IconButton>
                <Typography
                    component="h1"
                    variant="h6"
                    color="inherit"
                    noWrap
                    sx={{ flexGrow: 1 }}>
                    BatNoter
                </Typography>
                <Button variant="contained" color="secondary" endIcon={<Login />}>
                    Login
                </Button>
            </Toolbar>
        </AppBarComponent>
    )
}

export default AppBar
