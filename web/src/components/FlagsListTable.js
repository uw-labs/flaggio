import React from "react";
import AppBar from "@material-ui/core/AppBar/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Grid from "@material-ui/core/Grid";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import Tooltip from "@material-ui/core/Tooltip";
import IconButton from "@material-ui/core/IconButton";
import RefreshIcon from '@material-ui/icons/Refresh';
import SearchIcon from "@material-ui/icons/Search";
import {Chip, Switch, Table, TableBody, TableCell, TableHead, TableRow, withStyles} from "@material-ui/core";
import {useMutation, useQuery} from "@apollo/react-hooks";
import {FLAGS_QUERY, TOGGLE_FLAG_QUERY} from "../Queries";
import NewFlagModal from "./NewFlagModal";

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
  contentWrapper: {
    // margin: '40px 16px',
  },
});

function FlagsListTable(props) {
  const [open, setOpen] = React.useState(false);
  const handleClickOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);
  const {loading, error, data} = useQuery(FLAGS_QUERY);
  const [toggleFlag] = useMutation(TOGGLE_FLAG_QUERY);
  const {classes} = props;
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error :(</p>;
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
      <div className={classes.contentWrapper}>
        {/*<Typography color="textSecondary" align="center">*/}
        {/*  No users for this project yet*/}
        {/*</Typography>*/}
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
    </div>
  )
}

export default withStyles(styles)(FlagsListTable);
