import { Avatar, Button, Container, Grid, Typography } from '@mui/material'
import { SourceBranch } from 'mdi-material-ui'
import React, { ReactElement } from 'react'
import { User } from '../reducer/userSlice'
import RepoSelectDialog from './RepoSelectDialog'

interface Props {
  user: User | null
}

const Settings: React.FC<Props> = ({ user }): ReactElement => {
  const [openRepoSelectDialog, setOpenRepoSelectDialog] = React.useState(false);

  return (
    <Container maxWidth="sm">
      <Grid container direction="column" textAlign="center">
        <Grid container direction="column">
          <Grid flexGrow={1} sx={{ backgroundImage: `url('${user?.avatar_url}')`, backgroundPosition: "center", filter: "blur(30px)", height: "150px" }}  ></Grid>
          <Avatar alt={user?.name} src={user?.avatar_url} sx={{ width: 100, height: 100, alignSelf: "center", marginTop: "-100px" }} />
        </Grid>
        <Grid container direction="column" marginY={2}>
          <Typography m={0} variant="h5" gutterBottom component="div"> {user?.name || user?.email} </Typography>
          <Typography color="text.secondary" variant="body1" gutterBottom component="div"> {user?.location} </Typography>
          {user?.default_repo?.default_branch && <Typography color="text.secondary" m={0} variant="h6" gutterBottom component="div">Notes Repository: {user?.default_repo?.name} (<SourceBranch sx={{ verticalAlign: 'middle' }} fontSize='inherit' /> {user?.default_repo?.default_branch})</Typography>}
          <Button onClick={() => setOpenRepoSelectDialog(true)}>Change Notes Repository</Button>
          <RepoSelectDialog open={openRepoSelectDialog} setOpen={setOpenRepoSelectDialog} defaultRepo={user?.default_repo?.name} />
        </Grid>
      </Grid>
    </Container>
  )
}

export default Settings;
