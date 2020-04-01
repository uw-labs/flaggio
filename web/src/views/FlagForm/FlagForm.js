import React, { useEffect } from 'react';
import { makeStyles } from '@material-ui/styles';
import { Redirect, useParams } from 'react-router-dom';
import { Grid } from '@material-ui/core';
import { useMutation, useQuery } from '@apollo/react-hooks';
import { FlagDetails } from './components';
import {
  CREATE_FLAG_QUERY,
  CREATE_FLAG_RULE_QUERY,
  CREATE_VARIANT_QUERY,
  DELETE_FLAG_QUERY,
  DELETE_FLAG_RULE_QUERY,
  DELETE_VARIANT_QUERY,
  FLAG_QUERY,
  OPERATIONS_SEGMENTS_QUERY,
  UPDATE_FLAG_QUERY,
  UPDATE_FLAG_RULE_QUERY,
  UPDATE_VARIANT_QUERY,
} from './queries';
import { formatFlag, formatRule, formatVariant, newFlag, VariantTypes } from './models';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
    [theme.breakpoints.down('sm')]: {
      padding: theme.spacing(2),
    },
    [theme.breakpoints.down('xs')]: {
      padding: theme.spacing(0),
    },
  },
}));

const FlagForm = (props) => {
  const { id } = useParams();
  const [toFlagsPage, setToFlagsPage] = React.useState(false);
  const { loading, error, data = {} } = useQuery(FLAG_QUERY, {
    variables: { id },
    fetchPolicy: 'cache-and-network',
    skip: id === undefined,
  });
  const { loading: loadingOps, error: errorOps, data: dataOps } = useQuery(OPERATIONS_SEGMENTS_QUERY);
  const [deleteFlag] = useMutation(DELETE_FLAG_QUERY);
  const [createFlag] = useMutation(CREATE_FLAG_QUERY);
  const [updateFlag] = useMutation(UPDATE_FLAG_QUERY);
  const [createVariant] = useMutation(CREATE_VARIANT_QUERY);
  const [updateVariant] = useMutation(UPDATE_VARIANT_QUERY);
  const [deleteVariant] = useMutation(DELETE_VARIANT_QUERY);
  const [createRule] = useMutation(CREATE_FLAG_RULE_QUERY);
  const [updateRule] = useMutation(UPDATE_FLAG_RULE_QUERY);
  const [deleteRule] = useMutation(DELETE_FLAG_RULE_QUERY);
  const handleSaveFlag = async (flag, deletedItems) => {
    const variantsRef = {};
    if (flag.__new) {
      const { key, name } = flag;
      await createFlag({
        variables: { input: { key, name } },
        update(cache, { data: { createFlag: createdFlag } }) {
          flag.id = createdFlag.id;
        },
      });
    }
    await Promise.all(
      flag.variants.map(variant => {
        const variables = {
          id: variant.id,
          flagId: flag.id,
          input: formatVariant(variant),
        };
        if (variant.__new) {
          return createVariant({
            variables,
            update(cache, { data: { createVariant: createdVariant } }) {
              variantsRef[variant.id] = createdVariant.id;
            },
          });
        }
        if (variant.__changed) {
          return updateVariant({ variables });
        }
        return null;
      }),
    );
    await Promise.all([
      ...flag.rules.map(rule => {
        const variables = {
          id: rule.id,
          flagId: flag.id,
          input: formatRule(rule, variantsRef),
        };
        if (rule.__new) {
          return createRule({ variables });
        }
        if (rule.__changed) {
          return updateRule({ variables });
        }
        return null;
      }),
      ...deletedItems.map(item => {
        switch (item.type) {
          case 'variant':
            return deleteVariant({ variables: item });
          case 'rule':
            return deleteRule({ variables: item });
          default:
            return null;
        }
      }),
    ]);
    if (flag.__changed) {
      await updateFlag({ variables: { id: flag.id, input: formatFlag(flag, variantsRef) } });
    }
    setToFlagsPage(true);
  };
  const handleDeleteFlag = id => {
    deleteFlag({ variables: { id } })
      .then(() => setToFlagsPage(true));
  };
  const { flagType = VariantTypes.BOOLEAN } = props.location;
  useEffect(() => {
    const handleEsc = (event) => {
      if (event.key === 'Escape') setToFlagsPage(true);
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, []);
  const classes = useStyles();
  if (loading || loadingOps) return <div>"Loading..."</div>;
  if (error || errorOps) return <div>"Error while loading flag details :("</div>;

  return (
    <div className={classes.root}>
      {toFlagsPage && <Redirect to='/flags'/>}
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <FlagDetails
            flag={newFlag(data.flag, flagType)}
            operations={dataOps.operations.enumValues.map(v => v.name)}
            segments={dataOps.segments}
            onSaveFlag={handleSaveFlag}
            onDeleteFlag={handleDeleteFlag}
          />
        </Grid>
      </Grid>
    </div>
  );
};

export default FlagForm;
