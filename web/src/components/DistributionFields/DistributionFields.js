import React from 'react';
import { Chip, Grid, TextField, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/styles';
import PropTypes from 'prop-types';
import { isNumber } from 'lodash';
import { VariantTypes } from '../../helpers';
import { BooleanType } from '../../views/FlagForm/copy';

const useStyles = makeStyles(theme => ({
  mainGrid: {
    alignItems: 'center',
    alignContent: 'center',
    justify: 'center',
  },
  sumGrid: {
    marginTop: '10px',
  },
}));

const DistributionFields = props => {
  const classes = useStyles();
  const {
    distribution,
    variant,
    sum,
    onUpdateDistribution,
  } = props;

  return (
    <Grid container spacing={2} className={classes.mainGrid}>
      <Grid item sm={1} xs={false}/>
      <Grid item md={2} sm={3} xs={4}>
        <Chip
          label={typeof variant.value === VariantTypes.BOOLEAN ? BooleanType[variant.value] : variant.value}
          variant="outlined"
        />
      </Grid>
      <Grid item sm={3} xs={5}>
        <TextField
          id="outlined-basic"
          label="Percentage"
          name="percentage"
          onChange={onUpdateDistribution}
          variant="outlined"
          type="number"
          value={distribution.percentage}
          margin="dense"
        />
      </Grid>
      {isNumber(sum) && (
        <Grid item sm={2} xs={false} className={classes.sumGrid}>
          <Typography variant="subtitle2" gutterBottom>
            Total: {sum}%
          </Typography>
        </Grid>
      )}
    </Grid>
  );
};

DistributionFields.propTypes = {
  distribution: PropTypes.object.isRequired,
  variant: PropTypes.object.isRequired,
  sum: PropTypes.number,
  onUpdateDistribution: PropTypes.func.isRequired,
};

export default DistributionFields;
