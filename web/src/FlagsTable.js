import React from 'react';
import PropTypes from 'prop-types';
import clsx from 'clsx';
import {
  Chip,
  Paper,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TablePagination,
  TableRow,
  Toolbar,
  Typography
} from "@material-ui/core";
import {lighten, makeStyles} from "@material-ui/core/styles";
import './App.css';
import {useMutation, useQuery} from '@apollo/react-hooks';
import {gql} from 'apollo-boost';

export const FLAGS_QUERY = gql`
    {
        flags {
            id
            key
            name
            enabled
            createdAt
        }
    }
`;

const TOGGLE_FLAG_QUERY = gql`
    mutation($id: ID!, $input: UpdateFlag!) {
        updateFlag(id: $id, input: $input) {
            id
            enabled
        }
    }
`;

const useToolbarStyles = makeStyles(theme => ({
  root: {
    paddingLeft: theme.spacing(2),
    paddingRight: theme.spacing(1),
  },
  highlight:
    theme.palette.type === 'light'
      ? {
        color: theme.palette.secondary.main,
        backgroundColor: lighten(theme.palette.secondary.light, 0.85),
      }
      : {
        color: theme.palette.text.primary,
        backgroundColor: theme.palette.secondary.dark,
      },
  spacer: {
    flex: '1 1 100%',
  },
  actions: {
    color: theme.palette.text.secondary,
  },
  title: {
    flex: '0 0 auto',
  },
}));

const EnhancedTableToolbar = props => {
  const classes = useToolbarStyles();
  const {numSelected} = props;

  return (
    <Toolbar
      className={clsx(classes.root, {
        [classes.highlight]: numSelected > 0,
      })}
    >
      <div className={classes.title}>
        {numSelected > 0 ? (
          <Typography color="inherit" variant="subtitle1">
            {numSelected} selected
          </Typography>
        ) : (
          <Typography variant="h6" id="tableTitle">
            Flags
          </Typography>
        )}
      </div>
      <div className={classes.spacer}/>
    </Toolbar>
  );
};

EnhancedTableToolbar.propTypes = {
  numSelected: PropTypes.number.isRequired,
};


const useStyles = makeStyles(theme => ({
  root: {
    width: '100%',
    marginTop: theme.spacing(3),
  },
  paper: {
    width: '100%',
    marginBottom: theme.spacing(2),
  },
  table: {
    minWidth: 750,
  },
  tableWrapper: {
    overflowX: 'auto',
  },
  visuallyHidden: {
    border: 0,
    clip: 'rect(0 0 0 0)',
    height: 1,
    margin: -1,
    overflow: 'hidden',
    padding: 0,
    position: 'absolute',
    top: 20,
    width: 1,
  },
}));

function FlagsTable() {
  const classes = useStyles();
  const {loading, error, data} = useQuery(FLAGS_QUERY);
  const [toggleFlag] = useMutation(TOGGLE_FLAG_QUERY);
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error :(</p>;

  return (
    <div className={classes.root}>
      <Paper className={classes.paper}>
        <EnhancedTableToolbar numSelected={0}/>
        <div className={classes.tableWrapper}>
          <Table
            className={classes.table}
            aria-labelledby="tableTitle"
            aria-label="enhanced table"
          >
            <TableHead>
              <TableRow>
                <TableCell/>
                <TableCell component="th" scope="row" padding="none">Name</TableCell>
                <TableCell align="right">Key</TableCell>
                <TableCell align="right">Created At</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {data.flags.map(({id, key, name, enabled, createdAt}) => (
                <TableRow
                  hover
                  aria-checked={false}
                  tabIndex={-1}
                  key={key}
                >
                  <TableCell padding="checkbox">
                    <Switch
                      checked={enabled}
                      onChange={e => toggleFlag({variables: {id, input: {enabled: e.target.checked}}})}
                      color="primary"
                      inputProps={{'aria-label': 'primary checkbox'}}
                    />
                  </TableCell>
                  <TableCell component="th" id={1} scope="row" padding="none">
                    {name}
                  </TableCell>
                  <TableCell align="right"><Chip label={key}/></TableCell>
                  <TableCell align="right">{createdAt}</TableCell>
                </TableRow>
              ))
              }
            </TableBody>
          </Table>
        </div>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25]}
          component="div"
          count={data.flags.length}
          rowsPerPage={25}
          page={0}
          backIconButtonProps={{
            'aria-label': 'previous page',
          }}
          nextIconButtonProps={{
            'aria-label': 'next page',
          }}
          // onChangePage={handleChangePage}
          // onChangeRowsPerPage={handleChangeRowsPerPage}
        />
      </Paper>
    </div>
  );
}

export default FlagsTable;
