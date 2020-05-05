import React from 'react';
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
  Hidden,
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

const previewContext = (context) => {
  const maxItemsToDisplay = 3;
  let items = [];
  if (!context || typeof context !== 'object') {
    return '';
  }
  for (const key of Object.keys(context)) {
    if (!key.startsWith('$')) {
      items.push(`${key}: ${context[key]}`);
    }
  }
  let preview = items.slice(0, maxItemsToDisplay).join(', ');
  if (items.length > maxItemsToDisplay) {
    preview = `${preview} ...`
  }
  return preview;
};

const UsersTable = props => {
  const {
    className,
    users,
    page,
    rowsPerPage,
    rowsPerPageOptions,
    onPageChange,
    onRowsPerPageChange,
    ...rest
  } = props;

  const classes = useStyles();
  const handlePageChange = (event, page) => onPageChange(page);
  const handleRowsPerPageChange = event => onRowsPerPageChange(event.target.value);

  return (
    <Card
      {...rest}
      className={clsx(classes.root, className)}
    >
      <CardContent className={classes.content}>
        <PerfectScrollbar>
          <div>
            <Table>
              <colgroup>
                <col style={{ width: '30%' }}/>
                <col style={{ width: '40%' }}/>
                <col style={{ width: '30%' }}/>
              </colgroup>
              <TableHead>
                <TableRow>
                  <TableCell>User ID</TableCell>
                  <Hidden xsDown>
                    <TableCell>Context</TableCell>
                  </Hidden>
                  <Hidden smDown>
                    <TableCell>Updated</TableCell>
                  </Hidden>
                </TableRow>
              </TableHead>
              <TableBody>
                {users.users.map(user => (
                  <TableRow
                    className={classes.tableRow}
                    hover
                    key={user.id}
                  >
                    <TableCell>
                      <Link to={`/users/${user.id}`}>
                        <Chip size="small" variant="outlined" label={user.id} clickable/>
                      </Link>
                    </TableCell>
                    <TableCell>
                      <Link to={`/users/${user.id}`}>{previewContext(user.context)}</Link>
                    </TableCell>
                    <Hidden smDown>
                      <TableCell>
                        {moment(user.updatedAt).fromNow()}
                      </TableCell>
                    </Hidden>
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
          count={users.total}
          onChangePage={handlePageChange}
          onChangeRowsPerPage={handleRowsPerPageChange}
          page={page}
          rowsPerPage={rowsPerPage}
          rowsPerPageOptions={rowsPerPageOptions}
        />
      </CardActions>
    </Card>
  );
};

UsersTable.propTypes = {
  className: PropTypes.string,
  users: PropTypes.object.isRequired,
  page: PropTypes.number.isRequired,
  rowsPerPage: PropTypes.number.isRequired,
  rowsPerPageOptions: PropTypes.array.isRequired,
  onPageChange: PropTypes.func.isRequired,
  onRowsPerPageChange: PropTypes.func.isRequired,
};

export default UsersTable;
