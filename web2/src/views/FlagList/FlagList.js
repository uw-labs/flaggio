import React from 'react';
import { makeStyles } from '@material-ui/styles';
import { useMutation, useQuery } from '@apollo/react-hooks';
import { FlagsTable, FlagsToolbar } from './components';
import { FLAGS_QUERY, TOGGLE_FLAG_QUERY } from './queries';
import Typography from '../Typography';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(3),
  },
  content: {
    marginTop: theme.spacing(2),
  },
}));

function EmptyMessage({message}) {
  return (
    <Typography color="textSecondary" align="center" style={{margin: '40px 16px'}}>
      {message}
    </Typography>
  )
}

const FlagList = () => {
  const classes = useStyles();
  const {loading, error, data} = useQuery(FLAGS_QUERY);
  const [toggleFlag] = useMutation(TOGGLE_FLAG_QUERY);
  let content;
  switch (true) {
    case loading:
      content = <EmptyMessage message="No flags for this project yet"/>;
      break;
    case  error:
      content = <EmptyMessage message="There were an error while loading the flag list :("/>;
      break;
    default:
      content = <FlagsTable flags={data.flags} toggleFlag={toggleFlag}/>;
  }

  return (
    <div className={classes.root}>
      <FlagsToolbar/>
      <div className={classes.content}>
        {content}
      </div>
    </div>
  );
};

export default FlagList;
