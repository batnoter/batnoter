import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import { SourceBranch } from 'mdi-material-ui';
import React, { ReactElement } from 'react';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { getUserReposAsync, PreferenceStatus, saveDefaultRepoAsync, selectPreferenceStatus, selectUserRepos } from '../reducer/preferenceSlice';
import { getUserProfileAsync } from '../reducer/userSlice';

interface Props {
  defaultRepo?: string
  open: boolean
  setOpen?: (isOpen: boolean) => void
}

const RepoSelectDialog: React.FC<Props> = ({ open, setOpen, defaultRepo }): ReactElement => {
  const dispatch = useAppDispatch();
  React.useEffect(() => {
    dispatch(getUserReposAsync())
  }, [])
  const repos = useAppSelector(selectUserRepos);
  const prefStatus = useAppSelector(selectPreferenceStatus);

  const [repoName, setDefaultRepoName] = React.useState<string>();

  const handleChange = (event: SelectChangeEvent<typeof repoName>) => {
    setDefaultRepoName(String(event.target.value) || '');
  }

  const handleSave = async () => {
    const selectedRepo = repos.filter(r => r.name === repoName)[0]
    await dispatch(saveDefaultRepoAsync(selectedRepo))
    await dispatch(getUserProfileAsync())
    setOpen && setOpen(false)
  }

  const handleClose = (event: React.SyntheticEvent<unknown>, reason?: string) => {
    if (reason !== 'backdropClick') {
      setOpen && setOpen(false)
    }
  }

  return (
    <Dialog disableEscapeKeyDown open={open} onClose={handleClose} fullWidth>
      <DialogTitle>Select Notes Repository</DialogTitle>
      <DialogContent>
        <Box component="form" sx={{ display: 'flex', flexWrap: 'wrap' }}>
          <FormControl fullWidth sx={{ m: 1, minWidth: 120 }}>
            <InputLabel id="notes-repo-select-label">Notes Repository</InputLabel>
            <Select autoWidth labelId="notes-repo-select-label" value={repoName || defaultRepo} onChange={handleChange} disabled={prefStatus === PreferenceStatus.LOADING} label="Notes Repository">
              {repos.map(r => <MenuItem key={r.name} value={r.name}>{r.name} (<SourceBranch sx={{ verticalAlign: 'middle' }} fontSize='inherit' /> {r.default_branch || 'main'})</MenuItem>)}
            </Select>
          </FormControl>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button disabled={!repoName || prefStatus === PreferenceStatus.LOADING} onClick={() => handleSave()}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}

export default RepoSelectDialog;