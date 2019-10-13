import React from 'react';
import {AppBar, Container, CssBaseline, IconButton, Toolbar, Typography} from "@material-ui/core";
import {Menu as MenuIcon} from '@material-ui/icons';
import {makeStyles} from "@material-ui/core/styles";
import './App.css';
import FlagsTable from "./FlagsTable";
import ApolloClient from 'apollo-boost';
import { ApolloProvider } from '@apollo/react-hooks';

const client = new ApolloClient({
  uri: 'http://localhost:8081/query',
});

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
}));

function App() {
  const classes = useStyles();
  return (
    <ApolloProvider client={client}>
      <div>
        <CssBaseline />
        <AppBar position="static">
          <Toolbar>
            <IconButton edge="start" className={classes.menuButton} color="inherit" aria-label="menu">
              <MenuIcon/>
            </IconButton>
            <Typography variant="h6" className={classes.title}>
              Flaggio
            </Typography>
          </Toolbar>
        </AppBar>
        <Container fixed>
          <FlagsTable/>
        </Container>
      </div>
    </ApolloProvider>
  );
}

export default App;
