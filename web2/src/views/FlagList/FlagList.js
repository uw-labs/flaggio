import React, { useState } from 'react';
import { makeStyles } from '@material-ui/styles';

import { FlagsToolbar, FlagsTable } from './components';
import mockData from './data';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(3)
  },
  content: {
    marginTop: theme.spacing(2)
  }
}));

const FlagList = () => {
  const classes = useStyles();

  const [flags] = useState(mockData);

  return (
    <div className={classes.root}>
      <FlagsToolbar />
      <div className={classes.content}>
        <FlagsTable flags={flags} />
      </div>
    </div>
  );
};

export default FlagList;
