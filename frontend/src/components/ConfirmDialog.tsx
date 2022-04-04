
import { Button, Dialog, DialogActions, DialogContent, DialogProps, DialogTitle } from '@mui/material';
import React from 'react';

type Props = DialogProps & {
  desc: string
  onConfirm: () => void
}

const ConfirmDialog = (props: Props) => {
  const { desc, onConfirm, ...otherProps } = props;
  return (
    <Dialog {...otherProps}>
      <DialogTitle id="confirm-dialog">Please Confirm</DialogTitle>
      <DialogContent>{desc}</DialogContent>
      <DialogActions>
        <Button variant="contained" onClick={(e) => otherProps.onClose?.(e, "backdropClick")} color="inherit">CANCEL</Button>
        <Button variant="contained" onClick={(e) => { onConfirm(); otherProps.onClose?.(e, "backdropClick") }}>YES</Button>
      </DialogActions>
    </Dialog>
  );
};

export default ConfirmDialog;
