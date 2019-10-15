import React from 'react';
import {Route, Switch, useRouteMatch} from "react-router-dom";
import ListFlagsPage from "./ListFlagsPage";
import EditFlagPage from "./EditFlagPage";
import {withStyles} from "@material-ui/core";
import Header from "../theme/Header";
import Typography from "@material-ui/core/Typography";
import Link from "@material-ui/core/Link";

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="https://material-ui.com/">
        Your Website
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

const styles = theme => ({
  app: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
  },
  main: {
    flex: 1,
    padding: theme.spacing(6, 4),
    background: '#eaeff1',
  },
  footer: {
    padding: theme.spacing(2),
    background: '#eaeff1',
  },
});

function FlagsPage(props) {
  const {classes, onDrawerToggle} = props;
  let {path} = useRouteMatch();
  return (
    <div className={classes.app}>
      <Header title="Flags" onDrawerToggle={onDrawerToggle}/>
      <main className={classes.main}>
        <Switch>
          <Route exact path={path}>
            <ListFlagsPage/>
          </Route>
          <Route path={`${path}/:id`}>
            <EditFlagPage/>
          </Route>
        </Switch>
      </main>
      <footer className={classes.footer}>
        <Copyright/>
      </footer>
    </div>
  )
}

export default withStyles(styles)(FlagsPage);