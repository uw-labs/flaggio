import React, { useEffect, useState } from 'react';
import { makeStyles } from '@material-ui/styles';
import { Redirect, useParams } from 'react-router-dom';
import { Grid } from '@material-ui/core';
import { useMutation, useQuery } from '@apollo/react-hooks';
import { UserDetails } from './components';
import { DELETE_EVALUATION_QUERY, DELETE_USER_QUERY, USER_QUERY } from './queries';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
    [theme.breakpoints.down('sm')]: {
      padding: theme.spacing(2),
    },
    [theme.breakpoints.down('xs')]: {
      padding: theme.spacing(0),
    },
  },
}));

const rowsPerPageKey = 'rowsPerPage';
const rowsPerPageOptions = [10, 25, 50, { value: -1, label: 'All' }];

const UserForm = () => {
  const { id } = useParams();
  const [toUsersPage, setToUsersPage] = React.useState(false);
  const [search, setSearch] = useState();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(Number(localStorage.getItem(rowsPerPageKey) || 25));
  const { loading, error, data = {}, refetch } = useQuery(USER_QUERY, {
    variables: {
      id,
      search: search || undefined,
      offset: page * rowsPerPage,
      limit: rowsPerPage > 0 ? rowsPerPage : undefined,
    },
    fetchPolicy: 'cache-and-network',
    skip: id === undefined,
  });
  const [deleteUser] = useMutation(DELETE_USER_QUERY);
  const [deleteEvaluation] = useMutation(DELETE_EVALUATION_QUERY);
  const handleDeleteUser = id => {
    deleteUser({ variables: { id } })
      .then(() => setToUsersPage(true));
  };
  const handleDeleteEvaluation = id => {
    deleteEvaluation({ variables: { id } })
      .then(refetch);
  };
  useEffect(() => {
    const handleEsc = (event) => {
      if (event.key === 'Escape') setToUsersPage(true);
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, []);
  const classes = useStyles();
  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>"Error while loading user details :("</div>;

  return (
    <div className={classes.root}>
      {toUsersPage && <Redirect to='/users'/>}
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <UserDetails
            user={data.user}
            onDeleteUser={handleDeleteUser}
            onDeleteEvaluation={handleDeleteEvaluation}
            search={search}
            onSearch={v => {
              setSearch(v);
              setPage(0);
            }}
            rowsPerPage={rowsPerPage}
            rowsPerPageOptions={rowsPerPageOptions}
            onRowsPerPageChange={v => {
              localStorage.setItem(rowsPerPageKey, v);
              setRowsPerPage(v);
            }}
            page={page}
            onPageChange={setPage}
          />
        </Grid>
      </Grid>
    </div>
  );
};

export default UserForm;
