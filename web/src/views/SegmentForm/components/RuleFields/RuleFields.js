import React from 'react';
import PropTypes from 'prop-types';
import { Button, Grid, Paper, Tooltip } from '@material-ui/core';
import { makeStyles } from '@material-ui/styles';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import ConstraintFields from '../../../../components/ConstraintFields';
import * as copy from '../../copy';
import { OperationTypes } from '../../../../helpers';

const AllowMaxConstraints = 5;

const useStyles = makeStyles(theme => ({
  formControl: {
    fullWidth: true,
    display: 'flex',
    wrap: 'nowrap',
  },
  paper: {
    margin: theme.spacing(0, 0, 2, 0),
    padding: theme.spacing(2),
    display: 'flex',
    flexGrow: 1,
  },
  deleteRule: {
    display: 'flex',
    justifyContent: 'flex-end',
    alignItems: 'flex-end',
  },
  sideButton: {
    minWidth: theme.spacing(0),
  },
}));

const RuleFields = props => {
  const {
    rule,
    operations,
    onDeleteRule,
    onAddConstraint,
    onUpdateConstraint,
    onDeleteConstraint,
  } = props;
  const classes = useStyles();
  return (
    <Paper className={classes.paper}>
      <Grid container>
        <Grid item xs={12}>
          {rule.constraints.map((constraint, idx) => (
            <ConstraintFields
              key={constraint.id}
              constraint={constraint}
              showAddButton={idx === rule.constraints.length - 1 && rule.constraints.length < AllowMaxConstraints}
              first={idx === 0}
              operations={operations}
              copy={copy}
              operationTypes={OperationTypes}
              onAddConstraint={onAddConstraint}
              onUpdateConstraint={onUpdateConstraint(`constraints[${idx}].`)}
              onDeleteConstraint={onDeleteConstraint(constraint)}
            />
          ))}
        </Grid>
        <Grid item xs={10}/>
        <Grid item xs={2} className={classes.deleteRule}>
          <Tooltip title="Delete rule" placement="top">
            <Button
              size="small"
              color="secondary"
              className={classes.sideButton}
              onClick={onDeleteRule}
            >
              <DeleteOutlineIcon/>
            </Button>
          </Tooltip>
        </Grid>
      </Grid>
    </Paper>
  );
};

RuleFields.propTypes = {
  rule: PropTypes.object.isRequired,
  operations: PropTypes.arrayOf(PropTypes.string).isRequired,
  onDeleteRule: PropTypes.func.isRequired,
  onAddConstraint: PropTypes.func.isRequired,
  onUpdateConstraint: PropTypes.func.isRequired,
  onDeleteConstraint: PropTypes.func.isRequired,
};

export default RuleFields;
