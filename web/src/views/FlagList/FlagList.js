import React, { useState } from 'react';
import { makeStyles } from '@material-ui/styles';
import { useMutation, useQuery } from '@apollo/react-hooks';
import { FlagsTable, FlagsToolbar } from './components';
import { FLAGS_QUERY, TOGGLE_FLAG_QUERY } from './queries';
import { CircularProgress, Typography } from '@material-ui/core';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(3),
  },
  content: {
    marginTop: theme.spacing(2),
  },
  progress: {
    margin: theme.spacing(2),
  },
}));

const rowsPerPageKey = 'rowsPerPage';
const searchKey = 'flagSearch';
const rowsPerPageOptions = [10, 25, 50, { value: -1, label: 'All' }];

function EmptyMessage({ message }) {
  return (
    <Typography color="textSecondary" align="center" style={{ margin: '40px 16px' }}>
      {message}
    </Typography>
  )
}

const FlagList = () => {
  const classes = useStyles();
  const [search, setSearch] = useState(sessionStorage.getItem(searchKey) || '');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(Number(localStorage.getItem(rowsPerPageKey) || 25));
  const { loading, error, data } = useQuery(FLAGS_QUERY, {
    fetchPolicy: 'cache-and-network',
    variables: { search, offset: page * rowsPerPage, limit: rowsPerPage > 0 ? rowsPerPage : undefined },
  });
  const [toggleFlag] = useMutation(TOGGLE_FLAG_QUERY);
  let content;
  switch (true) {
    case loading:
      content = <CircularProgress className={classes.progress}/>;
      break;
    case error:
      content = <EmptyMessage message="There were an error while loading the flag list :("/>;
      break;
    case page === 0 && (!data || !data.flags || !data.flags.total):
      content = <EmptyMessage message="No flags for this project yet"/>;
      break;
    default:
      content = (
        <FlagsTable
          flags={data.flags}
          onToggleFlag={toggleFlag}
          page={page}
          rowsPerPage={rowsPerPage}
          rowsPerPageOptions={rowsPerPageOptions}
          onPageChange={setPage}
          onRowsPerPageChange={v => {
            localStorage.setItem(rowsPerPageKey, v);
            setRowsPerPage(v);
          }}
        />
      );
  }

  return (
    <div className={classes.root}>
      <FlagsToolbar
        search={search}
        onSearch={v => {
          setSearch(v);
          sessionStorage.setItem(searchKey, v);
          setPage(0);
        }}
      />
      <div className={classes.content}>
        {content}
      </div>
    </div>
  );
};

export default FlagList;
