import React from 'react';
import { Redirect, Switch } from 'react-router-dom';
import { RouteWithLayout } from './components';
import { Main as MainLayout, Minimal as MinimalLayout } from './layouts';
import { FlagForm as FlagFormView, FlagList as FlagListView, NotFound as NotFoundView } from './views';

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
      <RouteWithLayout
        component={FlagFormView}
        exact
        layout={MainLayout}
        path="/flags/:id"
      />
      <RouteWithLayout
        component={NotFoundView}
        exact
        layout={MinimalLayout}
        path="/not-found"
      />
      <Redirect to="/not-found"/>
    </Switch>
  );
};

export default Routes;
