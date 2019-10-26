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
  Avatar,
  Checkbox,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  TablePagination
} from '@material-ui/core';

import { getInitials } from 'helpers';

const useStyles = makeStyles(theme => ({
  root: {},
  content: {
    padding: 0
  },
  inner: {
    minWidth: 1050
  },
  nameContainer: {
    display: 'flex',
    alignItems: 'center'
  },
  avatar: {
    marginRight: theme.spacing(2)
  },
  actions: {
    justifyContent: 'flex-end'
  }
}));

const FlagsTable = props => {
  const { className, flags, ...rest } = props;

  const classes = useStyles();

  const [selectedFlags, setSelectedFlags] = useState([]);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [page, setPage] = useState(0);

  const handleSelectAll = event => {
    const { flags } = props;

    let selectedFlags;

    if (event.target.checked) {
      selectedFlags = flags.map(flag => flag.id);
    } else {
      selectedFlags = [];
    }

    setSelectedFlags(selectedFlags);
  };

  const handleSelectOne = (event, id) => {
    const selectedIndex = selectedFlags.indexOf(id);
    let newSelectedFlags = [];

    if (selectedIndex === -1) {
      newSelectedFlags = newSelectedFlags.concat(selectedFlags, id);
    } else if (selectedIndex === 0) {
      newSelectedFlags = newSelectedFlags.concat(selectedFlags.slice(1));
    } else if (selectedIndex === selectedFlags.length - 1) {
      newSelectedFlags = newSelectedFlags.concat(selectedFlags.slice(0, -1));
    } else if (selectedIndex > 0) {
      newSelectedFlags = newSelectedFlags.concat(
        selectedFlags.slice(0, selectedIndex),
        selectedFlags.slice(selectedIndex + 1)
      );
    }

    setSelectedFlags(newSelectedFlags);
  };

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
                    <Checkbox
                      checked={selectedFlags.length === flags.length}
                      color="primary"
                      indeterminate={
                        selectedFlags.length > 0 &&
                        selectedFlags.length < flags.length
                      }
                      onChange={handleSelectAll}
                    />
                  </TableCell>
                  <TableCell>Name</TableCell>
                  <TableCell>Email</TableCell>
                  <TableCell>Location</TableCell>
                  <TableCell>Phone</TableCell>
                  <TableCell>Registration date</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {flags.slice(0, rowsPerPage).map(flag => (
                  <TableRow
                    className={classes.tableRow}
                    hover
                    key={flag.id}
                    selected={selectedFlags.indexOf(flag.id) !== -1}
                  >
                    <TableCell padding="checkbox">
                      <Checkbox
                        checked={selectedFlags.indexOf(flag.id) !== -1}
                        color="primary"
                        onChange={event => handleSelectOne(event, flag.id)}
                        value="true"
                      />
                    </TableCell>
                    <TableCell>
                      <div className={classes.nameContainer}>
                        <Avatar
                          className={classes.avatar}
                          src={flag.avatarUrl}
                        >
                          {getInitials(flag.name)}
                        </Avatar>
                        <Typography variant="body1">{flag.name}</Typography>
                      </div>
                    </TableCell>
                    <TableCell>{flag.email}</TableCell>
                    <TableCell>
                      {flag.address.city}, {flag.address.state},{' '}
                      {flag.address.country}
                    </TableCell>
                    <TableCell>{flag.phone}</TableCell>
                    <TableCell>
                      {moment(flag.createdAt).format('DD/MM/YYYY')}
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
  flags: PropTypes.array.isRequired
};

export default FlagsTable;
