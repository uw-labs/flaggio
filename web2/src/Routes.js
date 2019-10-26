import React from 'react';
import { Switch, Redirect } from 'react-router-dom';

import { RouteWithLayout } from './components';
import { Main as MainLayout, Minimal as MinimalLayout } from './layouts';

import {
  FlagList as FlagListView,
} from './views';

const Routes = () => {
  return (
    <Switch>
      <Redirect
        exact
        from="/"
        to="/flags"
      />
      <RouteWithLayout
        component={FlagListView}
        exact
        layout={MainLayout}
        path="/flags"
      />
      {/*<RouteWithLayout*/}
      {/*  component={SignUpView}*/}
      {/*  exact*/}
      {/*  layout={MinimalLayout}*/}
      {/*  path="/sign-up"*/}
      {/*/>*/}
      {/*<RouteWithLayout*/}
      {/*  component={SignInView}*/}
      {/*  exact*/}
      {/*  layout={MinimalLayout}*/}
      {/*  path="/sign-in"*/}
      {/*/>*/}
      <Redirect to="/not-found" />
    </Switch>
  );
};

export default Routes;
