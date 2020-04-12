import React from 'react';
import PropTypes from 'prop-types';
import { Button, FormControl, Grid, InputLabel, MenuItem, Paper, Select, Tooltip } from '@material-ui/core';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import { makeStyles } from '@material-ui/styles';
import ConstraintFields from '../../../../components/ConstraintFields';
import * as copy from '../../copy';
import { BooleanType } from '../../copy';
import { OperationTypes, VariantTypes } from '../../models';

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
  } = props;
  const classes = useStyles();
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
              value={rule.distributions[0].variant.id}
              name="distributions[0].variant.id"
              onChange={onUpdateRule}
              labelWidth={50}
            >
              {variants.map(variant => (
                <MenuItem key={variant.id} value={variant.id}>
                  {typeof variant.value === VariantTypes.BOOLEAN ? BooleanType[variant.value] : variant.value}
                </MenuItem>
              ))}
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
};

export default RuleFields;
