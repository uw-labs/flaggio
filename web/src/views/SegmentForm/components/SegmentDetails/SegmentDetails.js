import React, { useState } from 'react';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/styles';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Divider,
  Grid,
  TextField,
  Tooltip,
} from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';
import { Link } from 'react-router-dom';
import { reject, set } from 'lodash';
import { DeleteSegmentDialog, RuleFields } from '../';
import { newConstraint, newRule } from '../../models';

const AllowMaxRules = 5;

const useStyles = makeStyles(theme => ({
  root: {},
  formControl: {
    fullWidth: true,
    display: 'flex',
    wrap: 'nowrap',
  },
  actionButton: {
    margin: theme.spacing(1),
  },
}));

const SegmentDetails = props => {
  const { className, segment: sgmnt, operations, onSaveSegment, onDeleteSegment, ...rest } = props;
  const [segment, setSegment] = useState(sgmnt);
  const [deleteSegmentDlgOpen, setDeleteSegmentDlgOpen] = React.useState(false);
  const [deletedItems, setDeletedItems] = React.useState([]);
  const classes = useStyles();

  // TODO: delete segments from flag rules
  const handleAddRule = () => setSegment({ ...segment, rules: [...segment.rules, newRule()] });
  const handleDelRule = rule => () => {
    setDeletedItems([...deletedItems, { type: 'rule', id: rule.id, segmentId: segment.id }]);
    setSegment({ ...segment, rules: reject(segment.rules, r => r === rule) });
  };
  const handleAddConstraint = rule => () => {
    setSegment({
      ...segment, rules: segment.rules.map(r => {
        if (r === rule) r = { ...r, constraints: [...r.constraints, newConstraint()] };
        return r;
      }),
    });
  };
  const handleDelConstraint = rule => constraint => () => {
    setSegment({
      ...segment, rules: segment.rules.map(r => {
        if (r === rule) r = { ...r, constraints: reject(r.constraints, c => c === constraint) };
        return r;
      }),
    });
  };
  const handleChange = (prefix = '', prefix2 = '') => event => {
    set(segment, `${prefix}__changed`, true);
    set(segment, `${prefix + prefix2}__changed`, true);
    setSegment(
      set({ ...segment }, `${prefix + prefix2}${event.target.name}`, event.target.value),
    );
  };
  const handleChange2Deep = (prefix = '') => (prefix2 = '') => handleChange(prefix, prefix2);
  const showAddRuleButton = segment.rules.length < AllowMaxRules;

  return (
    <Card
      {...rest}
      className={clsx(classes.root, className)}
    >
      <form autoComplete="off" noValidate>
        <CardHeader
          subheader="Defines a group of users that falls under a certain criteria"
          title="Segment"
        />
        <Divider/>
        <CardContent>
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Name"
                margin="dense"
                name="name"
                onChange={handleChange()}
                required
                value={segment.name}
                variant="outlined"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Description"
                margin="dense"
                name="description"
                onChange={handleChange()}
                value={segment.description}
                variant="outlined"
              />
            </Grid>
          </Grid>
        </CardContent>
        <CardHeader
          subheader="Based on a set of constraints, decide if a user belongs to this segment"
          title="Rules"
        />
        <Divider/>
        <CardContent>
          {
            segment.rules.map((rule, idx) => (
              <RuleFields
                key={rule.id}
                rule={rule}
                variants={segment.variants}
                operations={operations}
                onDeleteRule={handleDelRule(rule)}
                onAddConstraint={handleAddConstraint(rule)}
                onUpdateConstraint={handleChange2Deep(`rules[${idx}].`)}
                onDeleteConstraint={handleDelConstraint(rule)}
              />
            ))
          }
          {showAddRuleButton && (
            <Grid container>
              <Grid item xs={12}>
                <Button
                  color="inherit"
                  variant="outlined"
                  className={classes.actionButton}
                  onClick={handleAddRule}
                >
                  New rule
                </Button>
              </Grid>
            </Grid>
          )}
        </CardContent>

        <Divider/>
        <CardActions>
          {!segment.__new && (
            <>
              <DeleteSegmentDialog
                segment={segment}
                open={deleteSegmentDlgOpen}
                onConfirm={() => onDeleteSegment(segment.id)}
                onClose={() => setDeleteSegmentDlgOpen(false)}
              />
              <Tooltip title="Delete segment" placement="top">
                <Button
                  color="secondary"
                  onClick={() => setDeleteSegmentDlgOpen(true)}
                >
                  <DeleteIcon/>
                </Button>
              </Tooltip>
            </>
          )}
          <div style={{ flexGrow: 1 }}/>
          <Link to="/segments">
            <Button className={classes.actionButton}>Cancel</Button>
          </Link>
          <Button
            color="primary"
            variant="outlined"
            className={classes.actionButton}
            onClick={() => onSaveSegment(segment, deletedItems)}
          >
            Save
          </Button>
        </CardActions>
      </form>
    </Card>
  );
};

SegmentDetails.propTypes = {
  className: PropTypes.string,
  segment: PropTypes.object.isRequired,
  operations: PropTypes.arrayOf(PropTypes.string).isRequired,
  onSaveSegment: PropTypes.func.isRequired,
  onDeleteSegment: PropTypes.func.isRequired,
};

export default SegmentDetails;
