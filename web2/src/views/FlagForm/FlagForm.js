import React from 'react';
import { makeStyles } from '@material-ui/styles';
import { useParams } from 'react-router-dom';
import { Grid } from '@material-ui/core';
import { useQuery } from '@apollo/react-hooks';
import { FlagDetails } from './components';
import { FLAG_QUERY } from './queries';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
  },
}));

const FlagForm = () => {
  let {id} = useParams();
  const {loading, error, data} = useQuery(FLAG_QUERY, {variables: {id}});
  const classes = useStyles();
  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>"Error while loading flag details :("</div>;

  return (
    <div className={classes.root}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <FlagDetails flag={data.flag} operations={data.operations}/>
        </Grid>
      </Grid>
    </div>
  );
};

export default FlagForm;
