import React from 'react';
import PropTypes from 'prop-types';
import { includes } from 'lodash';
import {
  Button,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Select,
  TextField,
  Tooltip,
  Typography,
} from '@material-ui/core';
import RemoveIcon from '@material-ui/icons/RemoveCircleOutline';
import AddIcon from '@material-ui/icons/AddCircleOutline';
import { makeStyles } from '@material-ui/styles';
import ChipInput from 'material-ui-chip-input'
import { VariantTypes } from '../../helpers';
import { VariantType } from '../../views/FlagForm/copy';

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
  propertyGrid: {
    justifyContent: 'space-between',
    display: 'flex',
  },
  constraintLogicLabel: {
    minWidth: theme.spacing(5),
    textAlign: 'center',
    paddingRight: theme.spacing(1),
    marginTop: theme.spacing(2.5),
  },
}));

const ConstraintFields = props => {
  const classes = useStyles();
  const {
    constraint,
    showAddButton,
    first,
    operations,
    segments,
    copy: { BooleanType, Operations },
    operationTypes,
    onAddConstraint,
    onDeleteConstraint,
    onUpdateConstraint,
  } = props;
  const {
    ONE_OF,
    NOT_ONE_OF,
    IS_IN_SEGMENT,
    ISNT_IN_SEGMENT,
    EXISTS,
    DOESNT_EXIST,
    GREATER,
    GREATER_OR_EQUAL,
    LOWER,
    LOWER_OR_EQUAL,
  } = operationTypes;
  const disabledPropertyField = includes([IS_IN_SEGMENT, ISNT_IN_SEGMENT], constraint.operation);
  const showTypeField = includes([ONE_OF, NOT_ONE_OF], constraint.operation);
  const showSegmentInput = includes([IS_IN_SEGMENT, ISNT_IN_SEGMENT], constraint.operation);
  const showNumberInput = constraint.type === VariantTypes.NUMBER;
  const showBooleanInput = constraint.type === VariantTypes.BOOLEAN;
  const showNoInput = includes([EXISTS, DOESNT_EXIST], constraint.operation);

  return (
    <Grid container spacing={2}>
      <Grid item xs={3} className={classes.propertyGrid}>
        <Typography
          className={classes.constraintLogicLabel}
          variant="caption"
          display="block"
        >
          {first ? 'IF' : 'AND'}
        </Typography>
        <TextField
          label="Property"
          value={constraint.property}
          name="property"
          onChange={onUpdateConstraint}
          fullWidth
          disabled={disabledPropertyField}
          required={!disabledPropertyField}
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
              if (includes([IS_IN_SEGMENT, ISNT_IN_SEGMENT], e.target.value)) {
                onUpdateConstraint({ target: { name: 'property', value: '' } });
              }
              if (includes([GREATER, GREATER_OR_EQUAL, LOWER, LOWER_OR_EQUAL], e.target.value)) {
                onUpdateConstraint({ target: { name: 'type', value: VariantTypes.NUMBER } });
              } else {
                onUpdateConstraint({ target: { name: 'type', value: VariantTypes.STRING } });
              }
            }}
            labelWidth={70}
          >
            {operations.filter(op => !!Operations[op]).map(operation => (
              <MenuItem key={operation} value={operation}>{Operations[operation] || operation}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      {showTypeField && (
        <Grid item xs={1}>
          <FormControl
            className={classes.formControl}
            variant="outlined"
            required
          >
            <InputLabel>Type</InputLabel>
            <Select
              value={constraint.type}
              name="type"
              required
              onChange={e => {
                onUpdateConstraint(e);
                // reset constraint value
                if (e.target.value === VariantTypes.BOOLEAN) {
                  return onUpdateConstraint({ target: { name: 'values', value: [false] } });
                } else {
                  return onUpdateConstraint({ target: { name: 'values', value: [] } });
                }
              }}
              labelWidth={30}
            >
              {Object.keys(VariantTypes).map(type => (
                <MenuItem key={type} value={VariantTypes[type]}>{VariantType[type]}</MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
      )}
      <Grid item xs={showTypeField ? 5 : 6}>
        {showSegmentInput ? (
          <FormControl
            className={classes.formControl}
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
        ) : showNumberInput ? (
          <TextField
            label="Value"
            value={constraint.values[0] || ''}
            name="values[0]"
            type="number"
            onChange={onUpdateConstraint}
            fullWidth
            required={true}
            variant="outlined"
            InputProps={{ labelWidth: 45 }}
          />
        ) : showBooleanInput ? (
          <FormControl
            className={classes.formControl}
            variant="outlined"
          >
            <InputLabel>Value</InputLabel>
            <Select
              value={constraint.values[0] || false}
              name="values[0]"
              onChange={onUpdateConstraint}
              labelWidth={45}
            >
              {[true, false].map(val => (
                <MenuItem key={val} value={val}>{BooleanType[val]}</MenuItem>
              ))}
            </Select>
          </FormControl>
        ) : showNoInput ? (
          <div/>
        ) : (
          <ChipInput
            label="Values"
            name="values"
            blurBehavior="add"
            value={constraint.values}
            required={true}
            fullWidth
            variant="outlined"
            InputLabelProps={{ className: classes.valuesLabel }}
            onAdd={chip => onUpdateConstraint({ target: { name: 'values', value: [...constraint.values, chip] } })}
            onDelete={(chip, index) => onUpdateConstraint({
              target: {
                name: 'values',
                value: constraint.values.filter((v, idx) => idx !== index),
              },
            })}
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
          showAddButton ? (
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
  segments: PropTypes.arrayOf(PropTypes.object),
  showAddButton: PropTypes.bool.isRequired,
  first: PropTypes.bool.isRequired,
  copy: PropTypes.shape({
    BooleanType: PropTypes.object.isRequired,
    Operations: PropTypes.object.isRequired,
  }).isRequired,
  operationTypes: PropTypes.object.isRequired,
  onAddConstraint: PropTypes.func.isRequired,
  onUpdateConstraint: PropTypes.func.isRequired,
  onDeleteConstraint: PropTypes.func.isRequired,
};

export default ConstraintFields;
