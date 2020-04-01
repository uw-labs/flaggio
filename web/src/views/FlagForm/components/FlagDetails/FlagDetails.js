import React, { useState } from 'react';
import clsx from 'clsx';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/styles';
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Divider,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Select,
  Tab,
  Tabs,
  TextField,
  Tooltip,
} from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';
import { Link } from 'react-router-dom';
import { reject, set } from 'lodash';
import slugify from '@sindresorhus/slugify';
import { DeleteFlagDialog, RuleFields, VariantFields } from '../';
import { newConstraint, newRule, newVariant, VariantTypes } from '../../models';
import { BooleanType } from '../../copy';

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

const FlagDetails = props => {
  const { className, flag: flg, operations, segments, onSaveFlag, onDeleteFlag, ...rest } = props;
  const [flag, setFlag] = useState(flg);
  const [tab, setTab] = React.useState(0);
  const [deleteFlagDlgOpen, setDeleteFlagDlgOpen] = React.useState(false);
  const [deletedItems, setDeletedItems] = React.useState([]);
  const classes = useStyles();

  const handleAddVariant = () => setFlag({ ...flag, variants: [...flag.variants, newVariant()] });
  const handleDelVariant = variant => () => {
    setDeletedItems([...deletedItems, { type: 'variant', id: variant.id, flagId: flag.id }]);
    const newVariants = reject(flag.variants, v => v === variant);
    if (newVariants.length === 0) newVariants.push(newVariant());
    setFlag({ ...flag, variants: newVariants });
  };
  const handleAddRule = () => setFlag({ ...flag, rules: [...flag.rules, newRule()] });
  const handleDelRule = rule => () => {
    setDeletedItems([...deletedItems, { type: 'rule', id: rule.id, flagId: flag.id }]);
    setFlag({ ...flag, rules: reject(flag.rules, r => r === rule) });
  };
  const handleAddConstraint = rule => () => {
    setFlag({
      ...flag, rules: flag.rules.map(r => {
        if (r === rule) r = { ...r, constraints: [...r.constraints, newConstraint()] };
        return r;
      }),
    });
  };
  const handleDelConstraint = rule => constraint => () => {
    setFlag({
      ...flag, rules: flag.rules.map(r => {
        const newConstraints = reject(r.constraints, c => c === constraint);
        if (newConstraints.length === 0) newConstraints.push(newConstraint());
        if (r === rule) r = { ...r, constraints: newConstraints };
        return r;
      }),
    });
  };
  const handleChange = (prefix = '', prefix2 = '') => event => {
    set(flag, `${prefix}__changed`, true);
    set(flag, `${prefix + prefix2}__changed`, true);
    setFlag(
      set({ ...flag }, `${prefix + prefix2}${event.target.name}`, event.target.value),
    );
  };
  const handleChange2Deep = (prefix = '') => (prefix2 = '') => handleChange(prefix, prefix2);
  const getKeySlug = () => slugify(flag.name, { separator: '.' });
  const handleSetSlugifiedKey = () => {
    if (flag.key) return;
    handleChange()({ target: { name: 'key', value: getKeySlug() } });
  };

  return (
    <Card
      {...rest}
      className={clsx(classes.root, className)}
    >
      <form autoComplete="off" noValidate>
        <Tabs
          value={tab}
          onChange={(event, newValue) => setTab(newValue)}
          indicatorColor="primary"
          textColor="primary"
          variant="standard"
        >
          <Tab label="General"/>
          <Tab label="Rules"/>
        </Tabs>

        {/*********** GENERAL TAB ***********/}

        <Box role="tabpanel" value={tab} hidden={tab !== 0}>
          <CardHeader
            subheader="Identified by a key, a flag will return a value (variant) based on a set of rules"
            title="Flag"
          />
          <Divider/>
          <CardContent>
            <Grid container spacing={3}>
              <Grid item md={6} xs={12}>
                <TextField
                  fullWidth
                  label="Name"
                  margin="dense"
                  name="name"
                  onChange={handleChange()}
                  onBlur={handleSetSlugifiedKey}
                  required
                  autoFocus
                  value={flag.name}
                  variant="outlined"
                />
              </Grid>
              <Grid item md={6} xs={12}>
                <TextField
                  fullWidth
                  label="Key"
                  margin="dense"
                  name="key"
                  onChange={handleChange()}
                  required
                  value={flag.key || getKeySlug()}
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
                  value={flag.description}
                  variant="outlined"
                />
              </Grid>
            </Grid>
          </CardContent>
          <CardHeader
            subheader="Possible values this flag can return"
            title="Variants"
          />
          <Divider/>
          <CardContent>
            {
              flag.variants.map((variant, idx) => (
                <VariantFields
                  key={variant.id}
                  variant={variant}
                  isLast={idx === flag.variants.length - 1}
                  onAddVariant={handleAddVariant}
                  onUpdateVariant={handleChange(`variants[${idx}].`)}
                  onDeleteVariant={handleDelVariant(variant)}
                />
              ))
            }
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <FormControl
                  className={classes.formControl}
                  margin="dense"
                  variant="outlined"
                  required
                >
                  <InputLabel>If no rules are matched, return</InputLabel>
                  <Select
                    value={flag.defaultVariantWhenOn.id}
                    name="defaultVariantWhenOn.id"
                    onChange={handleChange()}
                    labelWidth={190}
                  >
                    {flag.variants.map(variant => (
                      <MenuItem key={variant.id} value={variant.id}>
                        {typeof variant.value === VariantTypes.BOOLEAN ? BooleanType[variant.value] : variant.value}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={6}>
                <FormControl
                  className={classes.formControl}
                  margin="dense"
                  variant="outlined"
                  required
                >
                  <InputLabel>If flag is disabled, return</InputLabel>
                  <Select
                    value={flag.defaultVariantWhenOff.id}
                    name="defaultVariantWhenOff.id"
                    onChange={handleChange()}
                    labelWidth={155}
                  >
                    {flag.variants.map(variant => (
                      <MenuItem key={variant.id} value={variant.id}>
                        {typeof variant.value === VariantTypes.BOOLEAN ? BooleanType[variant.value] : variant.value}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </CardContent>
        </Box>

        {/*********** RULES TAB ***********/}

        <Box role="tabpanel" value={tab} hidden={tab !== 1}>
          <CardHeader
            subheader="Based on a set of constraints, decide which value should be returned as result"
            title="Rules"
          />
          <Divider/>
          <CardContent>
            {
              flag.rules.map((rule, idx) => (
                <RuleFields
                  key={rule.id}
                  rule={rule}
                  variants={flag.variants}
                  operations={operations}
                  segments={segments}
                  onDeleteRule={handleDelRule(rule)}
                  onUpdateRule={handleChange(`rules[${idx}].`)}
                  onAddConstraint={handleAddConstraint(rule)}
                  onUpdateConstraint={handleChange2Deep(`rules[${idx}].`)}
                  onDeleteConstraint={handleDelConstraint(rule)}
                />
              ))
            }
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
          </CardContent>
        </Box>

        <Divider/>
        <CardActions>
          {!flag.__new && (
            <>
              <DeleteFlagDialog
                flag={flag}
                open={deleteFlagDlgOpen}
                onConfirm={() => onDeleteFlag(flag.id)}
                onClose={() => setDeleteFlagDlgOpen(false)}
              />
              <Tooltip title="Delete flag" placement="top">
                <Button
                  color="secondary"
                  onClick={() => setDeleteFlagDlgOpen(true)}
                >
                  <DeleteIcon/>
                </Button>
              </Tooltip>
            </>
          )}
          <div style={{ flexGrow: 1 }}/>
          <Link to="/flags">
            <Button className={classes.actionButton}>Cancel</Button>
          </Link>
          <Button
            color="primary"
            variant="outlined"
            className={classes.actionButton}
            onClick={() => onSaveFlag(flag, deletedItems)}
          >
            Save
          </Button>
        </CardActions>
      </form>
    </Card>
  );
};

FlagDetails.propTypes = {
  className: PropTypes.string,
  flag: PropTypes.object.isRequired,
  operations: PropTypes.arrayOf(PropTypes.string).isRequired,
  segments: PropTypes.arrayOf(PropTypes.object).isRequired,
  onSaveFlag: PropTypes.func.isRequired,
  onDeleteFlag: PropTypes.func.isRequired,
};

export default FlagDetails;
