import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import { makeStyles } from '@material-ui/styles';
import AddFlagButton from './AddFlagButton';
import { SearchInput } from '../../../../components';

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
    marginRight: theme.spacing(5),
  },
}));

const FlagsToolbar = props => {
  const { className, onSearch, ...rest } = props;

  const classes = useStyles();

  return (
    <div
      {...rest}
      className={clsx(classes.root, className)}
    >
      <div className={classes.row}>
        <SearchInput
          className={classes.searchInput}
          onChange={onSearch}
          placeholder="Search flag"
        />
        <AddFlagButton/>
      </div>
    </div>
  );
};

FlagsToolbar.propTypes = {
  className: PropTypes.string,
  onSearch: PropTypes.func.isRequired,
};

export default FlagsToolbar;
