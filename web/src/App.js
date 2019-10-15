import React from 'react';
import ApolloClient from 'apollo-boost';
import {ApolloProvider} from '@apollo/react-hooks';
import {BrowserRouter as Router} from "react-router-dom";
import Paperbase from "./theme/Paperbase";

const client = new ApolloClient({
  uri: 'http://localhost:8081/query',
});

function App() {
  return (
    <ApolloProvider client={client}>
      <Router>
        <Paperbase/>
      </Router>
    </ApolloProvider>
  );
}

export default App;
