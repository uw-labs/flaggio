import React from 'react';
import PropTypes from 'prop-types';
import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle } from '@material-ui/core';

const DeleteEvaluationDialog = ({ open, evaluation = {}, onConfirm, onClose }) => {
  return (
    <Dialog
      open={open}
      onClose={onClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">Delete evaluation?</DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          Are you sure you want to delete evaluation for flag "{evaluation.flagKey}"?
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

DeleteEvaluationDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  evaluation: PropTypes.object,
  onConfirm: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};

export default DeleteEvaluationDialog;
