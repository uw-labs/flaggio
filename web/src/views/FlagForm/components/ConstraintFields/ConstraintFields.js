import React from 'react';
import PropTypes from 'prop-types';
import { includes } from 'lodash';
import { Button, Chip, FormControl, Grid, InputLabel, MenuItem, Select, TextField, Tooltip } from '@material-ui/core';
import { blue } from '@material-ui/core/colors'
import RemoveIcon from '@material-ui/icons/RemoveCircleOutline';
import AddIcon from '@material-ui/icons/AddCircleOutline';
import { BooleanType, Operations } from '../../copy';
import { makeStyles } from '@material-ui/styles';
import ChipInput from 'material-ui-chip-input'
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
  valuesLabel: {
    background: 'white',
    padding: theme.spacing(0, 1, 0, 1),
  },
}));

const chipRenderer = ({ value, text, isFocused, isDisabled, isReadOnly, handleClick, handleDelete, className }, key) => (
  <Chip
    key={key}
    className={className}
    style={{
      pointerEvents: isDisabled || isReadOnly ? 'none' : undefined,
      backgroundColor: isFocused ? blue[400] : undefined,
    }}
    size="small"
    onClick={handleClick}
    onDelete={handleDelete}
    label={typeof text === 'boolean' ? BooleanType[text] : text}
  />
);

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
  const {
    IS_IN_SEGMENT,
    ISNT_IN_SEGMENT,
    EXISTS,
    DOESNT_EXIST,
    GREATER,
    GREATER_OR_EQUAL,
    LOWER,
    LOWER_OR_EQUAL,
  } = OperationTypes;
  const classes = useStyles();

  return (
    <Grid container spacing={2}>
      <Grid item xs={3}>
        <TextField
          label="Property"
          value={constraint.property}
          name="property"
          onChange={onUpdateConstraint}
          fullWidth
          disabled={includes([IS_IN_SEGMENT, ISNT_IN_SEGMENT], constraint.operation)}
          required={!includes([IS_IN_SEGMENT, ISNT_IN_SEGMENT], constraint.operation)}
          variant="outlined"
          InputProps={{ labelWidth: 65 }}
        />
      </Grid>
      <Grid item xs={2}>
        <FormControl
          className={classes.formControl}
          variant="outlined"
          required
        >
          <InputLabel>Operation</InputLabel>
          <Select
            value={constraint.operation}
            name="operation"
            onChange={e => {
              onUpdateConstraint(e);
              onUpdateConstraint({ target: { name: 'values', value: [] } });
            }}
            labelWidth={70}
          >
            {operations.filter(op => !!Operations[op]).map(operation => (
              <MenuItem key={operation} value={operation}>{Operations[operation] || operation}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      <Grid item xs={6}>
        {includes([IS_IN_SEGMENT, ISNT_IN_SEGMENT], constraint.operation) ? (
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
              labelWidth={60}
            >
              {segments.map(segment => (
                <MenuItem key={segment.id} value={segment.id}>{segment.name}</MenuItem>
              ))}
            </Select>
          </FormControl>
        ) : includes([GREATER, GREATER_OR_EQUAL, LOWER, LOWER_OR_EQUAL], constraint.operation) ? (
          <TextField
            label="Value"
            value={constraint.values[0]}
            name="values[0]"
            type="number"
            onChange={onUpdateConstraint}
            fullWidth
            required={!includes([EXISTS, DOESNT_EXIST], constraint.operation)}
            variant="outlined"
            InputProps={{ labelWidth: 45 }}
          />
        ) : includes([EXISTS, DOESNT_EXIST], constraint.operation) ? (
          <div/>
        ) : (
          <ChipInput
            chipRenderer={chipRenderer}
            label="Values"
            name="values"
            blurBehavior="add"
            defaultValue={constraint.values}
            disabled={includes([EXISTS, DOESNT_EXIST], constraint.operation)}
            required={!includes([EXISTS, DOESNT_EXIST], constraint.operation)}
            fullWidth
            variant="outlined"
            InputLabelProps={{ className: classes.valuesLabel }}
            onChange={value => onUpdateConstraint({ target: { name: 'values', value } })}
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
