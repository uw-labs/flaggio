import React, { useEffect } from 'react';
import { makeStyles } from '@material-ui/styles';
import { Redirect, useParams } from 'react-router-dom';
import { Grid } from '@material-ui/core';
import { useMutation, useQuery } from '@apollo/react-hooks';
import { SegmentDetails } from './components';
import {
  CREATE_SEGMENT_QUERY,
  CREATE_SEGMENT_RULE_QUERY,
  DELETE_SEGMENT_QUERY,
  DELETE_SEGMENT_RULE_QUERY,
  OPERATIONS_QUERY,
  SEGMENT_QUERY,
  UPDATE_SEGMENT_QUERY,
  UPDATE_SEGMENT_RULE_QUERY,
} from './queries';
import { formatSegmentRule, formatSegment, newSegment } from '../../helpers';

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

const SegmentForm = () => {
  const { id } = useParams();
  const [toSegmentsPage, setToSegmentsPage] = React.useState(false);
  const { loading, error, data = {}, refetch } = useQuery(SEGMENT_QUERY, {
    variables: { id },
    fetchPolicy: 'cache-and-network',
    skip: id === undefined,
  });
  const { loading: loadingOps, error: errorOps, data: dataOps } = useQuery(OPERATIONS_QUERY);
  const [deleteSegment] = useMutation(DELETE_SEGMENT_QUERY);
  const [createSegment] = useMutation(CREATE_SEGMENT_QUERY);
  const [updateSegment] = useMutation(UPDATE_SEGMENT_QUERY);
  const [createRule] = useMutation(CREATE_SEGMENT_RULE_QUERY);
  const [updateRule] = useMutation(UPDATE_SEGMENT_RULE_QUERY);
  const [deleteRule] = useMutation(DELETE_SEGMENT_RULE_QUERY);
  const handleSaveSegment = async (segment, deletedItems) => {
    const variantsRef = {};
    if (segment.__new) {
      const { name, description } = segment;
      await createSegment({
        variables: { input: { name, description } },
        update(cache, { data: { createSegment: createdSegment } }) {
          segment.id = createdSegment.id;
        },
      });
    }
    await Promise.all([
      ...segment.rules.map(rule => {
        const variables = {
          id: rule.id,
          segmentId: segment.id,
          input: formatSegmentRule(rule, variantsRef),
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
        if (item.type === 'rule') {
          return deleteRule({ variables: item });
        }
        return null;
      }),
    ]);
    if (segment.__changed) {
      await updateSegment({ variables: { id: segment.id, input: formatSegment(segment, variantsRef) } });
    }
    if (!segment.__new) {
      await refetch();
    }
    setToSegmentsPage(true);
  };
  const handleDeleteSegment = id => {
    deleteSegment({ variables: { id } })
      .then(() => setToSegmentsPage(true));
  };
  useEffect(() => {
    const handleEsc = (event) => {
      if (event.key === 'Escape') setToSegmentsPage(true);
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, []);
  const classes = useStyles();
  if (loading || loadingOps) return <div>"Loading..."</div>;
  if (error || errorOps) return <div>"Error while loading segment details :("</div>;

  return (
    <div className={classes.root}>
      {toSegmentsPage && <Redirect to='/segments'/>}
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <SegmentDetails
            segment={newSegment(data.segment)}
            operations={dataOps.operations.enumValues.map(v => v.name)}
            onSaveSegment={handleSaveSegment}
            onDeleteSegment={handleDeleteSegment}
          />
        </Grid>
      </Grid>
    </div>
  );
};

export default SegmentForm;
