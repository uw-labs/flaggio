import React from 'react';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/styles';
import {
  Button,
  Card,
  CardContent,
  CardHeader,
  Divider,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Tooltip,
} from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';
import { DeleteEvaluationDialog, DeleteUserDialog, EvaluationsTable, EvaluationsToolbar } from '../';

const useStyles = makeStyles(theme => ({
  root: {},
  formControl: {
    fullWidth: true,
    display: 'flex',
    wrap: 'nowrap',
  },
  actionButton: {
    margin: theme.spacing(1),
  },
  userCard: {
    marginBottom: theme.spacing(3),
  },
  evaluationsTable: {
    marginTop: theme.spacing(1),
  },
  deleteUserButton: {
    height: theme.spacing(6),
  },
}));

const displayContextValue = (value) => {
  if (['string', 'number'].includes(typeof value)) {
    return value;
  }
  return JSON.stringify(value);
};

const UserDetails = props => {
  const {
    className,
    user,
    onDeleteUser,
    onDeleteEvaluation,
    search,
    onSearch,
    page,
    onPageChange,
    rowsPerPage,
    rowsPerPageOptions,
    onRowsPerPageChange, ...rest
  } = props;
  const [deleteUserDlgOpen, setDeleteUserDlgOpen] = React.useState(false);
  const [deleteEvalDlgOpen, setDeleteEvalDlgOpen] = React.useState(false);
  const [evalToDelete, setEvalToDelete] = React.useState();
  const classes = useStyles();

  return (
    <>
      <Card
        {...rest}
        className={clsx(classes.root, classes.userCard, className)}
      >
        <DeleteUserDialog
          user={user}
          open={deleteUserDlgOpen}
          onConfirm={() => onDeleteUser(user.id)}
          onClose={() => setDeleteUserDlgOpen(false)}
        />
        <CardHeader
          subheader="A user is identified by an ID and has a list of properties (context)"
          title="User"
          action={
            <Tooltip title="Delete user" placement="top">
              <Button
                className={classes.deleteUserButton}
                color="secondary"
                onClick={() => setDeleteUserDlgOpen(true)}
              >
                <DeleteIcon/>
              </Button>
            </Tooltip>
          }
        />
        <Divider/>
        <CardContent>
          <Grid container spacing={3} justify="center">
            <Grid item md={8} sm={10} xs={12}>
              <Table className={classes.table} size="small" aria-label="a dense table">
                <TableHead>
                  <TableRow>
                    <TableCell>Property</TableCell>
                    <TableCell>Value</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {Object.keys(user.context).map((property) => (
                    <TableRow key={property}>
                      <TableCell component="th" scope="row">
                        {property}
                      </TableCell>
                      <TableCell>{displayContextValue(user.context[property])}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      <div className={classes.root}>
        <DeleteEvaluationDialog
          evaluation={evalToDelete}
          open={deleteEvalDlgOpen}
          onConfirm={() => {
            setDeleteEvalDlgOpen(false);
            onDeleteEvaluation(evalToDelete.id);
          }}
          onClose={() => setDeleteEvalDlgOpen(false)}
        />

        <EvaluationsToolbar
          search={search}
          onSearch={onSearch}
        />
        <div className={classes.evaluationsTable}>
          <EvaluationsTable
            evaluations={user.evaluations}
            page={page}
            rowsPerPage={rowsPerPage}
            rowsPerPageOptions={rowsPerPageOptions}
            onPageChange={onPageChange}
            onRowsPerPageChange={onRowsPerPageChange}
            onDeleteEvaluation={evltn => {
              setEvalToDelete(evltn);
              setDeleteEvalDlgOpen(true);
            }}
          />
        </div>
      </div>
    </>
  );
};

UserDetails.propTypes = {
  className: PropTypes.string,
  user: PropTypes.object.isRequired,
  search: PropTypes.string,
  onSearch: PropTypes.func.isRequired,
  page: PropTypes.number.isRequired,
  onPageChange: PropTypes.func.isRequired,
  rowsPerPage: PropTypes.number.isRequired,
  rowsPerPageOptions: PropTypes.array.isRequired,
  onRowsPerPageChange: PropTypes.func.isRequired,
  onDeleteUser: PropTypes.func.isRequired,
  onDeleteEvaluation: PropTypes.func.isRequired,
};

export default UserDetails;
