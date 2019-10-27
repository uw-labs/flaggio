import React, { useEffect } from 'react';
import { makeStyles } from '@material-ui/styles';
import { Redirect, useParams } from 'react-router-dom';
import { Grid } from '@material-ui/core';
import { useMutation, useQuery } from '@apollo/react-hooks';
import { FlagDetails } from './components';
import { DELETE_FLAG_QUERY, FLAG_QUERY } from './queries';
import { newFlag } from './models';
import { FLAGS_QUERY } from '../FlagList/queries';
import { reject } from 'lodash';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
  },
}));

const FlagForm = () => {
  const { id } = useParams();
  const [toFlagsPage, setToFlagsPage] = React.useState(false);
  const { loading, error, data } = useQuery(FLAG_QUERY, { variables: { id } });
  const [deleteFlag] = useMutation(DELETE_FLAG_QUERY, {
    update(cache, { data: { deleteFlag: id } }) {
      const { flags } = cache.readQuery({ query: FLAGS_QUERY });
      cache.writeQuery({
        query: FLAGS_QUERY,
        data: { flags: reject(flags, { id }) },
      });
    },
  });
  useEffect(() => {
    const handleEsc = (event) => {
      if (event.key === 'Escape') setToFlagsPage(true);
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, []);
  const classes = useStyles();
  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>"Error while loading flag details :("</div>;
  const handleDeleteFlag = id => {
    deleteFlag({ variables: { id } })
      .then(() => setToFlagsPage(true));
  };

  return (
    <div className={classes.root}>
      {toFlagsPage && <Redirect to='/flags'/>}
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <FlagDetails
            flag={newFlag(data.flag)}
            operations={data.operations.enumValues.map(v => v.name)}
            onDeleteFlag={handleDeleteFlag}
          />
        </Grid>
      </Grid>
    </div>
  );
};

export default FlagForm;
