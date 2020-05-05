import React from 'react';
import { Redirect, Switch } from 'react-router-dom';
import { RouteWithLayout } from './components';
import { Main as MainLayout, Minimal as MinimalLayout } from './layouts';
import {
  FlagForm as FlagFormView,
  FlagList as FlagListView,
  NotFound as NotFoundView,
  SegmentList as SegmentListView,
  SegmentForm as SegmentFormView,
  UserList as UserListView,
  UserForm as UserFormView,
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
      <RouteWithLayout
        component={FlagFormView}
        exact
        layout={MainLayout}
        path="/flags/new"
      />
      <RouteWithLayout
        component={FlagFormView}
        exact
        layout={MainLayout}
        path="/flags/:id"
      />
      <RouteWithLayout
        component={SegmentListView}
        exact
        layout={MainLayout}
        path="/segments"
      />
      <RouteWithLayout
        component={SegmentFormView}
        exact
        layout={MainLayout}
        path="/segments/new"
      />
      <RouteWithLayout
        component={SegmentFormView}
        exact
        layout={MainLayout}
        path="/segments/:id"
      />
      <RouteWithLayout
        component={UserListView}
        exact
        layout={MainLayout}
        path="/users"
      />
      <RouteWithLayout
        component={UserFormView}
        exact
        layout={MainLayout}
        path="/users/:id"
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
