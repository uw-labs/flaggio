import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { Button, Divider, FormControl, Grid, InputLabel, MenuItem, Paper, Select, Tooltip } from '@material-ui/core';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import { makeStyles } from '@material-ui/styles';
import ConstraintFields from '../../../../components/ConstraintFields';
import DistributionFields from '../../../../components/DistributionFields';
import * as copy from '../../copy';
import { BooleanType } from '../../copy';
import { OperationTypes, PercentageRollout, VariantTypes } from '../../models';

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
    variants,
    segments,
    operations,
    onUpdateRule,
    onDeleteRule,
    onAddConstraint,
    onUpdateConstraint,
    onDeleteConstraint,
    onUpdateDistribution,
  } = props;
  const classes = useStyles();
  const sum = rule.distributions.reduce((total, d) => (total + Number(d.percentage)), 0);

  return (
    <Paper className={classes.paper}>
      <Grid container spacing={1}>
        <Grid item xs={12}>
          {rule.constraints.map((constraint, idx) => (
            <ConstraintFields
              key={constraint.id}
              constraint={constraint}
              segments={segments}
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
        <Grid item xs={11}>
          <FormControl
            className={classes.formControl}
            margin="dense"
            variant="outlined"
            required
          >
            <InputLabel>Return</InputLabel>
            <Select
              value={rule.returnVariant}
              name="returnVariant"
              onChange={onUpdateRule}
              labelWidth={50}
            >
              {variants.map(variant => (
                <MenuItem key={variant.id} value={variant.id}>
                  {typeof variant.value === VariantTypes.BOOLEAN ? BooleanType[variant.value] : variant.value}
                </MenuItem>
              ))}
              <MenuItem key="distribution" value={PercentageRollout}>
                a percentage rollout
              </MenuItem>
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={1} className={classes.deleteRule}>
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

        {rule.returnVariant === PercentageRollout && (
          <Grid item xs={12}>
            <Divider/>
            {rule.distributions.map((distribution, idx) => (
              <Fragment key={distribution.id}>
                <DistributionFields
                  key={distribution.id}
                  distribution={distribution}
                  sum={idx === 0 ? sum : null}
                  variant={variants.find(v => v.id === distribution.variant.id)}
                  onUpdateDistribution={onUpdateDistribution(`distributions[${idx}].`)}
                />
                <Divider variant="middle"/>
              </Fragment>
            ))}
          </Grid>
        )}
      </Grid>
    </Paper>
  );
};

RuleFields.propTypes = {
  rule: PropTypes.object.isRequired,
  variants: PropTypes.arrayOf(PropTypes.object).isRequired,
  operations: PropTypes.arrayOf(PropTypes.string).isRequired,
  segments: PropTypes.arrayOf(PropTypes.object).isRequired,
  onUpdateRule: PropTypes.func.isRequired,
  onDeleteRule: PropTypes.func.isRequired,
  onAddConstraint: PropTypes.func.isRequired,
  onUpdateConstraint: PropTypes.func.isRequired,
  onDeleteConstraint: PropTypes.func.isRequired,
  onUpdateDistribution: PropTypes.func.isRequired,
};

export default RuleFields;
