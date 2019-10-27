import React, { useState } from 'react';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import moment from 'moment';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { makeStyles } from '@material-ui/styles';
import {
  Card,
  CardActions,
  CardContent,
  Chip,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TablePagination,
  TableRow,
} from '@material-ui/core';
import { Link } from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  root: {},
  content: {
    padding: 0,
  },
  inner: {
    minWidth: 1050,
  },
  nameContainer: {
    display: 'flex',
    alignItems: 'center',
  },
  avatar: {
    marginRight: theme.spacing(2),
  },
  actions: {
    justifyContent: 'flex-end',
  },
}));

const FlagsTable = props => {
  const {className, flags, toggleFlag, ...rest} = props;

  const classes = useStyles();

  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [page, setPage] = useState(0);

  const handlePageChange = (event, page) => {
    setPage(page);
  };

  const handleRowsPerPageChange = event => {
    setRowsPerPage(event.target.value);
  };

  return (
    <Card
      {...rest}
      className={clsx(classes.root, className)}
    >
      <CardContent className={classes.content}>
        <PerfectScrollbar>
          <div className={classes.inner}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell padding="checkbox">
                    &nbsp;
                  </TableCell>
                  <TableCell>Name</TableCell>
                  <TableCell>Key</TableCell>
                  <TableCell>Created</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {flags.map(flag => (
                  <TableRow
                    className={classes.tableRow}
                    hover
                    key={flag.id}
                  >
                    <TableCell padding="checkbox">
                      <Switch
                        checked={flag.enabled}
                        onChange={e => toggleFlag({variables: {id: flag.id, input: {enabled: e.target.checked}}})}
                        color="primary"
                        inputProps={{'aria-label': 'primary checkbox'}}
                      />
                    </TableCell>
                    <TableCell>
                      <Link to={`/flags/${flag.id}`}>{flag.name}</Link>
                    </TableCell>
                    <TableCell>
                      <Chip size="small" variant="outlined" label={flag.key}/>
                    </TableCell>
                    <TableCell>
                      {moment(flag.createdAt).fromNow()}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </PerfectScrollbar>
      </CardContent>
      <CardActions className={classes.actions}>
        <TablePagination
          component="div"
          count={flags.length}
          onChangePage={handlePageChange}
          onChangeRowsPerPage={handleRowsPerPageChange}
          page={page}
          rowsPerPage={rowsPerPage}
          rowsPerPageOptions={[5, 10, 25]}
        />
      </CardActions>
    </Card>
  );
};

FlagsTable.propTypes = {
  className: PropTypes.string,
  flags: PropTypes.array.isRequired,
};

export default FlagsTable;