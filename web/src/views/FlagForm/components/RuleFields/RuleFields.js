import React from 'react';
import { Button, FormControl, Grid, InputLabel, MenuItem, Paper, Select, Tooltip } from '@material-ui/core';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import ConstraintFields from '../ConstraintFields';
import { BooleanType } from '../../copy';
import { VariantTypes } from '../../models';
import { makeStyles } from '@material-ui/styles';
import PropTypes from 'prop-types';

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
      <Grid container>
        <Grid item xs={12}>
          {rule.constraints.map((constraint, idx) => (
            <ConstraintFields
              key={constraint.id}
              constraint={constraint}
              segments={segments}
              isLast={idx === rule.constraints.length - 1}
              operations={operations}
              onAddConstraint={onAddConstraint}
              onUpdateConstraint={onUpdateConstraint(`constraints[${idx}].`)}
              onDeleteConstraint={onDeleteConstraint(constraint)}
            />
          ))}
        </Grid>
        <Grid item xs={8}>
          <FormControl
            className={classes.formControl}
            margin="dense"
            variant="outlined"
          >
            <InputLabel>Return</InputLabel>
            <Select
              value={rule.distributions[0].variant.id}
              name="distributions[0].variant.id"
              onChange={onUpdateRule}
            >
              {variants.map(variant => (
                <MenuItem key={variant.id} value={variant.id}>
                  {typeof variant.value === VariantTypes.BOOLEAN ? BooleanType[variant.value] : variant.value}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={2}/>
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
