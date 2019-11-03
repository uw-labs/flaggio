import React from 'react';
import { makeStyles } from '@material-ui/styles';
import { useQuery } from '@apollo/react-hooks';
import { SegmentsTable, SegmentsToolbar } from './components';
import { SEGMENTS_QUERY } from './queries';
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

const SegmentList = () => {
  const classes = useStyles();
  const { loading, error, data } = useQuery(SEGMENTS_QUERY);
  let content;
  switch (true) {
    case loading:
      content = <CircularProgress className={classes.progress}/>;
      break;
    case  error:
      content = <EmptyMessage message="There were an error while loading the segment list :("/>;
      break;
    case !data || !data.segments || data.segments.length === 0:
      content = <EmptyMessage message="No segments for this project yet"/>;
      break;
    default:
      content = <SegmentsTable segments={data.segments}/>;
  }

  return (
    <div className={classes.root}>
      <SegmentsToolbar/>
      <div className={classes.content}>
        {content}
      </div>
    </div>
  );
};

export default SegmentList;
