
import { Button, Dialog, DialogActions, DialogContent, DialogProps, DialogTitle } from '@mui/material';
import React from 'react';

type Props = DialogProps & {
  desc: string
  onConfirm: () => void
}

const ConfirmDialog: React.FC<Props> = (props: Props) => {
  const { desc, onConfirm, ...otherProps } = props;
  return (
    <Dialog {...otherProps}>
      <DialogTitle id="confirm-dialog">Please Confirm</DialogTitle>
      <DialogContent>{desc}</DialogContent>
      <DialogActions>
        <Button variant="outlined" onClick={(e) => otherProps.onClose?.(e, "backdropClick")}>CANCEL</Button>
        <Button variant="contained" onClick={(e) => { onConfirm(); otherProps.onClose?.(e, "backdropClick") }}>YES</Button>
      </DialogActions>
    </Dialog>
  );
};

export default ConfirmDialog;
