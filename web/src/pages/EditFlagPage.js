import React from 'react';
import {
  useParams
} from "react-router-dom";
import Header from "../theme/Header";
import Content from "../theme/Content";
import {withStyles} from "@material-ui/core";

const styles = theme => ({
  main: {
    flex: 1,
    padding: theme.spacing(6, 4),
    background: '#eaeff1',
  },
});

function EditFlagPage(props) {
  const {classes} = props;
  let { id } = useParams();
  return (
    <main className={classes.main}>
      <Content/>
    </main>
  )
}

export default withStyles(styles)(EditFlagPage);