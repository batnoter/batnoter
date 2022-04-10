import { Alert, Typography } from '@mui/material';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import { unwrapResult } from '@reduxjs/toolkit';
import React, { ReactElement } from 'react';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { APIStatus, APIStatusType } from '../reducer/common';
import { autoSetupRepoAsync, selectPreferenceStatus } from '../reducer/preferenceSlice';
import { getSanitizedErrorMessage } from '../util/util';
import RepoSelectDialog from './RepoSelectDialog';

interface Props {
  open: boolean
  setOpen?: (isOpen: boolean) => void
}

const autoSetupRepoName = "notes";

const isLoading = (apiStatus: APIStatus): boolean => {
  const { autoSetupRepoAsync } = apiStatus;
  return autoSetupRepoAsync === APIStatusType.LOADING;
}

const isFailed = (apiStatus: APIStatus): boolean => {
  const { autoSetupRepoAsync } = apiStatus;
  return autoSetupRepoAsync === APIStatusType.FAIL;
}

const RepoSetupDialog: React.FC<Props> = ({ open, setOpen }): ReactElement => {
  const [openRepoSelectDialog, setOpenRepoSelectDialog] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState("");

  const dispatch = useAppDispatch();
  const apiStatus = useAppSelector(selectPreferenceStatus);

  const handleRepoSelect = () => {
    setOpenRepoSelectDialog(true);
  }

  const handleAutoSetupRepo = async () => {
    await dispatch(autoSetupRepoAsync(autoSetupRepoName)).then(unwrapResult)
      .catch(err => setErrorMessage(getSanitizedErrorMessage(err)));
    setOpen && setOpen(false);
  }

  return (
    <Dialog disableEscapeKeyDown open={open} fullWidth>
      <DialogTitle>Setup Notes Repository</DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexWrap: 'wrap' }}>
          {isFailed(apiStatus) && <Alert severity="error" sx={{ width: "100%" }}>{errorMessage}</Alert>}
          <Typography gutterBottom paragraph>
            You may choose to automatically setup your notes repository or manually select an existing repository for storing notes.
            The automatic setup will create a new private repository &quot;{autoSetupRepoName}&quot; and set it as your notes repository.
          </Typography>
          <Typography gutterBottom paragraph>
            Do you want to setup the notes repository automatically?
          </Typography>

          {isFailed(apiStatus) && <Alert severity="warning" sx={{ width: "100%" }}>
            If you already have repository with name: &quot;{autoSetupRepoName}&quot; Then please use SELECT EXISTING REPO option.
          </Alert>}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button disabled={isLoading(apiStatus)} onClick={() => handleRepoSelect()}>SELECT EXISTING REPO</Button>
        <Button disabled={isLoading(apiStatus)} onClick={() => handleAutoSetupRepo()}>YES, SETUP AUTOMATICALLY</Button>
      </DialogActions>
      <RepoSelectDialog open={openRepoSelectDialog} setOpen={setOpenRepoSelectDialog} />
    </Dialog>
  );
}

export default RepoSetupDialog;
