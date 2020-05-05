import React from 'react';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import moment from 'moment';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { makeStyles } from '@material-ui/styles';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  Hidden,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TablePagination,
  TableRow,
  Tooltip,
} from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';

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

const displayContextValue = (value) => {
  if (['string', 'number'].includes(typeof value)) {
    return value;
  }
  return JSON.stringify(value);
};

const EvaluationsTable = props => {
  const {
    className,
    evaluations,
    page,
    rowsPerPage,
    rowsPerPageOptions,
    onPageChange,
    onRowsPerPageChange,
    onDeleteEvaluation,
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
                <col style={{ width: '35%' }}/>
                <col style={{ width: '30%' }}/>
                <col style={{ width: '5%' }}/>
              </colgroup>
              <TableHead>
                <TableRow>
                  <TableCell>Flag Key</TableCell>
                  <Hidden xsDown>
                    <TableCell>Value</TableCell>
                  </Hidden>
                  <Hidden smDown>
                    <TableCell>Evaluated</TableCell>
                  </Hidden>
                  <TableCell/>
                </TableRow>
              </TableHead>
              <TableBody>
                {evaluations.evaluations.map(evaluation => (
                  <TableRow
                    className={classes.tableRow}
                    hover
                    key={evaluation.id}
                  >
                    <TableCell>
                      {evaluation.flagKey}
                    </TableCell>
                    <TableCell>
                      {displayContextValue(evaluation.value)}
                    </TableCell>
                    <Hidden smDown>
                      <TableCell>
                        {moment(evaluation.createdAt).fromNow()}
                      </TableCell>
                    </Hidden>
                    <TableCell>
                      <Tooltip title="Delete evaluation" placement="top">
                        <Button
                          color="secondary"
                          onClick={() => onDeleteEvaluation(evaluation)}
                        >
                          <DeleteIcon/>
                        </Button>
                      </Tooltip>
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
          count={evaluations.total}
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

EvaluationsTable.propTypes = {
  className: PropTypes.string,
  evaluations: PropTypes.object.isRequired,
  page: PropTypes.number.isRequired,
  rowsPerPage: PropTypes.number.isRequired,
  rowsPerPageOptions: PropTypes.array.isRequired,
  onPageChange: PropTypes.func.isRequired,
  onRowsPerPageChange: PropTypes.func.isRequired,
  onDeleteEvaluation: PropTypes.func.isRequired,
};

export default EvaluationsTable;
