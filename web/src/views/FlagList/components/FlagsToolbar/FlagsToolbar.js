import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import { makeStyles } from '@material-ui/styles';
import { Button } from '@material-ui/core';
import { Link } from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  root: {},
  row: {
    height: '42px',
    display: 'flex',
    alignItems: 'center',
    marginTop: theme.spacing(1),
  },
  spacer: {
    flexGrow: 1,
  },
  searchInput: {
    marginRight: theme.spacing(1),
  },
}));

const FlagsToolbar = props => {
  const { className, ...rest } = props;

  const classes = useStyles();

  return (
    <div
      {...rest}
      className={clsx(classes.root, className)}
    >
      <div className={classes.row}>
        <span className={classes.spacer}/>
        <Link to="/flags/new">
          <Button
            color="primary"
            variant="contained"
          >
            Add flag
          </Button>
        </Link>
      </div>
      {/*<div className={classes.row}>*/}
      {/*  <SearchInput*/}
      {/*    className={classes.searchInput}*/}
      {/*    placeholder="Search flag"*/}
      {/*  />*/}
      {/*</div>*/}
    </div>
  );
};

FlagsToolbar.propTypes = {
  className: PropTypes.string,
};

export default FlagsToolbar;
