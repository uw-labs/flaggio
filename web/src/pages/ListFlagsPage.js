import React from 'react';
import Content from "../theme/Content";
import {withStyles} from "@material-ui/core";
import FlagsListTable from "../components/FlagsListTable";


const styles = theme => ({});

function ListFlagsPage(props) {
  return (
    <Content>
      <FlagsListTable/>
    </Content>
  );
}

export default withStyles(styles)(ListFlagsPage);