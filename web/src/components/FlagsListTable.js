import React from "react";
import AppBar from "@material-ui/core/AppBar/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import SearchIcon from "@material-ui/icons/Search";
import {
  Chip,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  withStyles
} from "@material-ui/core";
import {useMutation, useQuery} from "@apollo/react-hooks";
import {FLAGS_QUERY, TOGGLE_FLAG_QUERY} from "../Queries";
import NewFlagModal from "./NewFlagModal";
import {Link} from "react-router-dom";
import moment from 'moment';

const styles = theme => ({
  searchBar: {
    borderBottom: '1px solid rgba(0, 0, 0, 0.12)',
  },
  searchInput: {
    fontSize: theme.typography.fontSize,
  },
  block: {
    display: 'block',
  },
  addFlag: {
    marginRight: theme.spacing(1),
  },
});

function FlagsTable(props) {
  const {classes, toggleFlag, flags} = props;
  return (
    <Table
      className={classes.table}
      aria-labelledby="tableTitle"
      aria-label="enhanced table"
    >
      <TableHead>
        <TableRow>
          <TableCell/>
          <TableCell>Name</TableCell>
          <TableCell>Key</TableCell>
          <TableCell align="right">Added</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {flags.map(({id, key, name, enabled, createdAt}) => (
          <TableRow
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
            <TableCell>
              <Link to={`/flags/${id}`}>
                {name}
              </Link>
            </TableCell>
            <TableCell><Chip label={key}/></TableCell>
            <TableCell align="right">
              {moment(createdAt).fromNow()}
            </TableCell>
          </TableRow>
        ))
        }
      </TableBody>
    </Table>
  );
}

function EmptyMessage({message}) {
  return (
    <Typography color="textSecondary" align="center" style={{margin: '40px 16px'}}>
      {message}
    </Typography>
  )
}

function FlagsListTable(props) {
  const {classes} = props;
  const [open, setOpen] = React.useState(false);
  const handleClickOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);
  const {loading, error, data} = useQuery(FLAGS_QUERY);
  const [toggleFlag] = useMutation(TOGGLE_FLAG_QUERY);
  let content;
  switch (true) {
    case loading:
      content = <EmptyMessage message="No flags for this project yet"/>;
      break;
    case  error:
      content = <EmptyMessage message="There were an error while loading the flag list :("/>;
      break;
    default:
      content = <FlagsTable classes={classes} flags={data.flags} toggleFlag={toggleFlag}/>;
  }
  return (
    <div>
      <AppBar className={classes.searchBar} position="static" color="default" elevation={0}>
        <Toolbar>
          <Grid container spacing={2} alignItems="center">
            <Grid item>
              <SearchIcon className={classes.block} color="inherit"/>
            </Grid>
            <Grid item xs>
              <TextField
                fullWidth
                placeholder="Search by flag name, key, or ID"
                InputProps={{
                  disableUnderline: true,
                  className: classes.searchInput,
                }}
              />
            </Grid>
            <NewFlagModal open={open} handleClose={handleClose}/>
            <Grid item>
              <Button variant="contained" color="primary" className={classes.addFlag} onClick={handleClickOpen}>
                Add flag
              </Button>
            </Grid>
          </Grid>
        </Toolbar>
      </AppBar>
      {content}
    </div>
  )
}

export default withStyles(styles)(FlagsListTable);
