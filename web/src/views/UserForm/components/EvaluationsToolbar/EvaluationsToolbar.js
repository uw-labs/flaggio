import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import { makeStyles } from '@material-ui/styles';
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
}));

const EvaluationsToolbar = props => {
  const { className, search, onSearch, ...rest } = props;

  const classes = useStyles();

  return (
    <div
      {...rest}
      className={clsx(classes.root, className)}
    >
      <div className={classes.row}>
        <SearchInput
          defaultValue={search}
          onChange={onSearch}
          placeholder="Search evaluations"
        />
      </div>
    </div>
  );
};

EvaluationsToolbar.propTypes = {
  className: PropTypes.string,
  search: PropTypes.string,
  onSearch: PropTypes.func.isRequired,
};

export default EvaluationsToolbar;
