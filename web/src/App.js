import React, { Component } from 'react';
import { Router } from 'react-router-dom';
import ApolloClient from 'apollo-boost';
import { ApolloProvider } from '@apollo/react-hooks';
import { createBrowserHistory } from 'history';
import { ThemeProvider } from '@material-ui/styles';
import validate from 'validate.js';
import theme from './theme';
import 'react-perfect-scrollbar/dist/css/styles.css';
import './assets/scss/index.scss';
import validators from './common/validators';
import Routes from './Routes';

const browserHistory = createBrowserHistory();

const client = new ApolloClient({
  // uri: 'https://flags.vkt.sh/query',
  uri: 'http://localhost:8081/query',
});


validate.validators = {
  ...validate.validators,
  ...validators,
};

export default class App extends Component {
  render() {
    return (
      <ApolloProvider client={client}>
        <ThemeProvider theme={theme}>
          <Router history={browserHistory}>
            <Routes/>
          </Router>
        </ThemeProvider>
      </ApolloProvider>
    );
  }
}
