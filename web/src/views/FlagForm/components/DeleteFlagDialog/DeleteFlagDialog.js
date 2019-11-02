import React from 'react';
import PropTypes from 'prop-types';
import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from '@material-ui/core';

const DeleteFlagDialog = ({open, flag, onConfirm, onClose}) => {
  return (
    <Dialog
      open={open}
      onClose={onClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">Delete flag?</DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          Are you sure you want to delete flag "{flag.name}"?
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="primary">
          No, keep it
        </Button>
        <Button onClick={onConfirm} color="secondary" autoFocus>
          Yes, delete it
        </Button>
      </DialogActions>
    </Dialog>
  );
};

DeleteFlagDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  flag: PropTypes.object.isRequired,
  onConfirm: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};

export default DeleteFlagDialog;
