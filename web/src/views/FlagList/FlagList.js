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

function EmptyMessage({ message }) {
  return (
    <Typography color="textSecondary" align="center" style={{ margin: '40px 16px' }}>
      {message}
    </Typography>
  )
}

const FlagList = () => {
  const classes = useStyles();
  const [search, setSearch] = useState();
  const { loading, error, data } = useQuery(FLAGS_QUERY, { variables: { search } });
  const [toggleFlag] = useMutation(TOGGLE_FLAG_QUERY);
  let content;
  switch (true) {
    case loading:
      content = <CircularProgress className={classes.progress}/>;
      break;
    case  error:
      content = <EmptyMessage message="There were an error while loading the flag list :("/>;
      break;
    case !data || !data.flags || data.flags.length === 0:
      content = <EmptyMessage message="No flags for this project yet"/>;
      break;
    default:
      content = <FlagsTable flags={data.flags} toggleFlag={toggleFlag}/>;
  }

  return (
    <div className={classes.root}>
      <FlagsToolbar onSearch={setSearch}/>
      <div className={classes.content}>
        {content}
      </div>
    </div>
  );
};

export default FlagList;
