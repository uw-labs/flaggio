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
    display: "flex",
    flexGrow: 1,
  },
  deleteRule: {
    display: 'flex',
    justifyContent: 'flex-end',alignItems: 'flex-end'
  },
  sideButton: {
    minWidth: theme.spacing(0),
  },
}));

const RuleFields = ({ rule, variants, operations, onDeleteRule, onAddConstraint, onDeleteConstraint }) => {
  const classes = useStyles();
  return (
    <Paper className={classes.paper}>
      <Grid container>
        <Grid item xs={12}>
          {rule.constraints.map((constraint, idx) => (
            <ConstraintFields
              key={constraint.id}
              constraint={constraint}
              isLast={idx === rule.constraints.length - 1}
              operations={operations}
              onAddConstraint={onAddConstraint}
              onDeleteConstraint={onDeleteConstraint}
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
              // onChange={e => setDefaults({...defaults, defaultWhenOn: e.target.value})}
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
  onDeleteRule: PropTypes.func.isRequired,
  onAddConstraint: PropTypes.func.isRequired,
  onDeleteConstraint: PropTypes.func.isRequired,
};

export default RuleFields;
