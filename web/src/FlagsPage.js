import React from 'react';
import {Route, Switch, useRouteMatch} from "react-router-dom";
import ListFlagsPage from "./ListFlagsPage";
import EditFlagPage from "./EditFlagPage";

function FlagsPage() {
  let {path} = useRouteMatch();
  return (
    <div>
      <Switch>
        <Route exact path={path}>
          <ListFlagsPage/>
        </Route>
        <Route path={`${path}/:id`}>
          <EditFlagPage/>
        </Route>
      </Switch>
    </div>
  )
}

export default FlagsPage;