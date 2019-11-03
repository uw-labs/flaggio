import React from 'react';
import PropTypes from 'prop-types';
import { includes } from 'lodash';
import { Button, FormControl, Grid, InputLabel, MenuItem, Select, TextField, Tooltip } from '@material-ui/core';
import RemoveIcon from '@material-ui/icons/RemoveCircleOutline';
import AddIcon from '@material-ui/icons/AddCircleOutline';
import { Operations } from '../../copy';
import { makeStyles } from '@material-ui/styles';
import { OperationTypes } from '../../models';

const useStyles = makeStyles(theme => ({
  formControl: {
    fullWidth: true,
    display: 'flex',
    wrap: 'nowrap',
  },
  sideButton: {
    minWidth: theme.spacing(0),
  },
  actionButtons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
}));

const ConstraintFields = props => {
  const {
    constraint,
    isLast,
    operations,
    segments,
    onAddConstraint,
    onDeleteConstraint,
    onUpdateConstraint,
  } = props;
  const classes = useStyles();
  return (
    <Grid container spacing={1}>
      <Grid item xs={4}>
        <TextField
          label="Property"
          value={constraint.property}
          margin="dense"
          name="property"
          onChange={onUpdateConstraint}
          fullWidth
          disabled={includes([OperationTypes.IS_IN_SEGMENT, OperationTypes.ISNT_IN_SEGMENT], constraint.operation)}
          required={!includes([OperationTypes.IS_IN_SEGMENT, OperationTypes.ISNT_IN_SEGMENT], constraint.operation)}
          variant="outlined"
          InputProps={{labelWidth:"65"}}
        />
      </Grid>
      <Grid item xs={4}>
        <FormControl
          className={classes.formControl}
          margin="dense"
          variant="outlined"
          required
        >
          <InputLabel>Operation</InputLabel>
          <Select
            value={constraint.operation}
            name="operation"
            onChange={onUpdateConstraint}
            labelWidth="70"
          >
            {operations.map(operation => (
              <MenuItem key={operation} value={operation}>{Operations[operation] || operation}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      <Grid item xs={3}>
        {includes([OperationTypes.IS_IN_SEGMENT, OperationTypes.ISNT_IN_SEGMENT], constraint.operation) ? (
          <FormControl
            className={classes.formControl}
            margin="dense"
            variant="outlined"
          >
            <InputLabel>Segment</InputLabel>
            <Select
              value={constraint.values[0]}
              name="values[0]"
              required
              onChange={onUpdateConstraint}
              labelWidth="60"
            >
              {segments.map(segment => (
                <MenuItem key={segment.id} value={segment.id}>{segment.name}</MenuItem>
              ))}
            </Select>
          </FormControl>
          ) : (
          <TextField
            label="Value"
            value={constraint.values[0]}
            margin="dense"
            name="values[0]"
            onChange={onUpdateConstraint}
            fullWidth
            disabled={includes([OperationTypes.EXISTS, OperationTypes.DOESNT_EXIST], constraint.operation)}
            required={!includes([OperationTypes.EXISTS, OperationTypes.DOESNT_EXIST], constraint.operation)}
            variant="outlined"
            InputProps={{labelWidth:"45"}}
          />
        )}
      </Grid>
      <Grid item xs={1} className={classes.actionButtons}>
        <Tooltip title="Delete constraint" placement="top">
          <Button
            size="small"
            color="secondary"
            className={classes.sideButton}
            onClick={onDeleteConstraint}
          >
            <RemoveIcon/>
          </Button>
        </Tooltip>
        {
          isLast ? (
            <Tooltip title="New constraint" placement="top">
              <Button
                size="small"
                color="primary"
                className={classes.sideButton}
                onClick={onAddConstraint}
              >
                <AddIcon/>
              </Button>
            </Tooltip>
          ) : (
            <Button
              size="small"
              color="primary"
              className={classes.sideButton}
              disabled
            >
              <AddIcon style={{ visibility: 'hidden' }}/>
            </Button>
          )
        }
      </Grid>
    </Grid>
  );
};

ConstraintFields.propTypes = {
  constraint: PropTypes.object.isRequired,
  segments: PropTypes.arrayOf(PropTypes.object).isRequired,
  isLast: PropTypes.bool.isRequired,
  onAddConstraint: PropTypes.func.isRequired,
  onUpdateConstraint: PropTypes.func.isRequired,
  onDeleteConstraint: PropTypes.func.isRequired,
};

export default ConstraintFields;
