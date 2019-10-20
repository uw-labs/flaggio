import React from 'react';
import {Redirect, useParams} from "react-router-dom";
import {reject} from "lodash";
import Content from "../theme/Content";
import {
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  FormControlLabel,
  MenuItem,
  Paper,
  Select,
  TextField,
  withStyles
} from "@material-ui/core";
import {Delete as DeleteIcon} from "@material-ui/icons";
import Grid from "@material-ui/core/Grid";
import {useMutation, useQuery} from "@apollo/react-hooks";
import {DELETE_FLAG_QUERY, FLAG_QUERY, FLAGS_QUERY} from "../Queries";

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
    margin: theme.spacing(2),
  },
  textField: {
    // marginLeft: theme.spacing(1),
    // marginRight: theme.spacing(1),
  },
  section1: {
    margin: theme.spacing(0, 0, 2, 0),
  },
  section2: {
    margin: theme.spacing(0),
  },
  footer: {
    marginTop: theme.spacing(1),
    paddingTop: theme.spacing(1),
    borderTop: 'solid #E4E4E4 1px',
  },
  footerDiv: {
    flexGrow: 1,
  }
});

function FlagForm({classes, flag: flg, operations, handleDeleteFlag}) {
  const [flag, setFlag] = React.useState(flg);
  const [deleteFlagDlgOpen, setDeleteFlagDlgOpen] = React.useState(false);
  const setFlagField = (field, value) => setFlag({...flag, [field]: value});
  const handleClickOpen = () => setDeleteFlagDlgOpen(true);
  const handleClose = () => setDeleteFlagDlgOpen(false);
  const confirmDeleteFlag = () => {
    handleDeleteFlag(flag.id);
    setDeleteFlagDlgOpen(false);
  };
  return (
    <form className={classes.container} noValidate autoComplete="off">
      <Grid container spacing={2} className={classes.section1}>
        <Grid item xs={12} sm={6}>
          <TextField
            label="Name"
            className={classes.textField}
            value={flag.name}
            onChange={e => setFlagField('name', e.target.value)}
            fullWidth
          />
        </Grid>
        <Grid item xs={12} sm={6}>
          <TextField
            label="Key"
            className={classes.textField}
            value={flag.key}
            onChange={e => setFlagField('key', e.target.value)}
            fullWidth
          />
        </Grid>
        <Grid item xs={12} sm={12}>
          <TextField
            label="Description"
            className={classes.textField}
            value={flag.description || ''}
            onChange={e => setFlagField('description', e.target.value)}
            fullWidth
          />
        </Grid>
      </Grid>

      {
        flag.variants.map(variant => (
          <Grid container spacing={2} className={classes.section2} key={variant.id}>
            <Grid item xs={6} sm={6}>
              <TextField
                label="Value"
                className={classes.textField}
                value={variant.value}
                // onChange={e => setFlag('name', e.target.value)}
                fullWidth
              />
            </Grid>
            <Grid item xs={6} sm={6}>
              <TextField
                label="Name"
                className={classes.textField}
                value={variant.name}
                // onChange={e => setFlag('name', e.target.value)}
                fullWidth
              />
            </Grid>
            <Grid item xs={4} sm={4}>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={variant.defaultWhenOn}
                    // onChange={handleChange('checkedB')}
                    color="primary"
                  />
                }
                label="Use when flag is on"
              />
            </Grid>
            <Grid item xs={4} sm={4}>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={variant.defaultWhenOff}
                    // onChange={handleChange('checkedB')}
                    color="primary"
                  />
                }
                label="Use when flag is off"
              />
            </Grid>
          </Grid>
        ))
      }

      {
        flag.rules.map(rule => (
          <Paper key={rule.id}>
            {rule.constraints.map(constraint => (
              <Grid container spacing={2} className={classes.section2} key={constraint.id}>
                <Grid item xs={4} sm={4}>
                  <TextField
                    label="Property"
                    className={classes.textField}
                    value={constraint.property}
                    // onChange={e => setFlag('name', e.target.value)}
                    fullWidth
                  />
                </Grid>
                <Grid item xs={4} sm={4}>
                  <Select
                    value={constraint.operation}
                    // onChange={handleChange}
                  >
                    {operations.map(operation => (
                      <MenuItem key={operation} value={operation.name}>{operation.name}</MenuItem>
                    ))}
                  </Select>
                </Grid>
                <Grid item xs={4} sm={4}>
                  <TextField
                    label="Value"
                    className={classes.textField}
                    value={constraint.values.join(',')}
                    // onChange={e => setFlag('name', e.target.value)}
                    fullWidth
                  />
                </Grid>
              </Grid>
            ))}
          </Paper>
        ))
      }
      <Grid container alignContent="space-around" direction="row-reverse" className={classes.footer}>
        <Grid item>
          <Button color="primary">
            Save
          </Button>
        </Grid>
        <Grid item>
          <Button color="secondary">
            Cancel
          </Button>
        </Grid>
        <Grid item>
          <DeleteFlagDialog open={deleteFlagDlgOpen} onConfirm={confirmDeleteFlag} handleClose={handleClose}
                            flag={flag}/>
          <Button color="secondary" onClick={handleClickOpen}>
            <DeleteIcon/>
          </Button>
        </Grid>
      </Grid>
    </form>
  )
}

function DeleteFlagDialog({open, flag, onConfirm, handleClose}) {
  return (
    <Dialog
      open={open}
      onClose={handleClose}
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
        <Button onClick={handleClose} color="primary">
          No, keep it
        </Button>
        <Button onClick={onConfirm} color="secondary" autoFocus>
          Yes, delete it
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function EditFlagPage({classes}) {
  const [toFlagsPage, setToFlagsPage] = React.useState(false);
  let {id} = useParams();
  const {loading, error, data} = useQuery(FLAG_QUERY, {variables: {id}});
  const [deleteFlag] = useMutation(DELETE_FLAG_QUERY, {
    update(cache, {data: {deleteFlag: id}}) {
      const {flags} = cache.readQuery({query: FLAGS_QUERY});
      cache.writeQuery({
        query: FLAGS_QUERY,
        data: {flags: reject(flags, {id})},
      });
    }
  });
  if (toFlagsPage === true) {
    return <Redirect to='/flags'/>
  }
  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>"Error while loading flag details :("</div>;
  const handleDeleteFlag = (id) => {
    deleteFlag({variables: {id}}).then(() => setToFlagsPage(true));
  };

  return (
    <Content>
      <FlagForm classes={classes}
                flag={data.flag}
                operations={data.operations.enumValues}
                handleDeleteFlag={handleDeleteFlag}/>
    </Content>
  )
}

export default withStyles(styles)(EditFlagPage);