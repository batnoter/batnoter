import { Avatar, Button, Container, Grid, Typography } from '@mui/material'
import React from 'react'
import { User } from '../reducer/userSlice'
import RepoSelectDialog from './RepoSelectDialog'

interface Props {
  user: User | null
}

const Settings: React.FC<Props> = ({ user }) => {
  const [openRepoSelectDialog, setOpenRepoSelectDialog] = React.useState(false);

  return (

    <Container maxWidth="sm">
      <Grid
        container
        direction="column">
        <Grid flexDirection={'column'} justifyContent={'center'} display="flex" >
          <Grid flexGrow={1} sx={{ backgroundImage: `url('${user?.avatar_url}')`, backgroundPosition: "center", filter: "blur(30px)", height: "150px" }}  ></Grid>
          <Avatar alt={user?.name} src={user?.avatar_url} sx={{ width: 100, height: 100, alignSelf: "center", marginTop: "-100px" }} />

        </Grid>
        <Grid flexDirection={'column'} justifyContent={'center'} display="flex" marginY={2}>
          <Typography m={0} variant="h5" gutterBottom component="div"> {user?.name} </Typography>
          <Typography color="textSecondary" variant="body1" gutterBottom component="div"> {user?.location} </Typography>
          {user?.default_repo?.default_branch && <Typography color="textSecondary" m={0} variant="h6" gutterBottom component="div"> Default Repository: {user?.default_repo?.name} (Branch: {user?.default_repo?.default_branch})</Typography>}
          <Button onClick={() => setOpenRepoSelectDialog(true)}> Change Default Repository </Button>
          <RepoSelectDialog open={openRepoSelectDialog} setOpen={setOpenRepoSelectDialog} defaultRepo={user?.default_repo?.name} />
        </Grid>
      </Grid>
    </Container>
  )
}

export default Settings;