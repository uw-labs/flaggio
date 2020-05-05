import React  from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import { debounce } from 'lodash';
import { makeStyles } from '@material-ui/styles';
import { Input, Paper } from '@material-ui/core';
import SearchIcon from '@material-ui/icons/Search';

const useStyles = makeStyles(theme => ({
  root: {
    borderRadius: '4px',
    alignItems: 'center',
    padding: theme.spacing(1),
    display: 'flex',
    flexGrow: 1,
  },
  icon: {
    marginRight: theme.spacing(1),
    color: theme.palette.text.secondary,
  },
  input: {
    flexGrow: 1,
    fontSize: '14px',
    lineHeight: '16px',
    letterSpacing: '-0.05px',
  },
}));

const SearchInput = props => {
  const { className, onChange, style, ...rest } = props;

  const classes = useStyles();
  const handler = debounce(onChange, 300);

  return (
    <Paper
      {...rest}
      className={clsx(classes.root, className)}
      style={style}
    >
      <SearchIcon className={classes.icon}/>
      <Input
        {...rest}
        autoFocus
        className={classes.input}
        disableUnderline
        onChange={e => handler(e.target.value)}
      />
    </Paper>
  );
};

SearchInput.propTypes = {
  className: PropTypes.string,
  onChange: PropTypes.func.isRequired,
  style: PropTypes.object,
};

export default SearchInput;
