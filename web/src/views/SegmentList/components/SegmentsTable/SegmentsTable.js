import React, { useState } from 'react';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import moment from 'moment';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { makeStyles } from '@material-ui/styles';
import {
  Card,
  CardActions,
  CardContent, Hidden,
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

const SegmentsTable = props => {
  const { className, segments, toggleSegment, ...rest } = props;

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
          <div>
            <Table>
              <colgroup>
                <col style={{ width: '30%' }}/>
                <col style={{ width: '40%' }}/>
                <col style={{ width: '30%' }}/>
              </colgroup>
              <TableHead>
                <TableRow>
                  <TableCell>Name</TableCell>
                  <Hidden xsDown>
                    <TableCell>Description</TableCell>
                  </Hidden>
                  <Hidden smDown>
                    <TableCell>Created</TableCell>
                  </Hidden>
                </TableRow>
              </TableHead>
              <TableBody>
                {segments.map(segment => (
                  <TableRow
                    className={classes.tableRow}
                    hover
                    key={segment.id}
                  >
                    <TableCell>
                      <Link to={`/segments/${segment.id}`}>{segment.name}</Link>
                    </TableCell>
                    <Hidden xsDown>
                      <TableCell>
                        {segment.description}
                      </TableCell>
                    </Hidden>
                    <Hidden smDown>
                      <TableCell>
                        {moment(segment.createdAt).fromNow()}
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
          count={segments.length}
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

SegmentsTable.propTypes = {
  className: PropTypes.string,
  segments: PropTypes.array.isRequired,
};

export default SegmentsTable;
