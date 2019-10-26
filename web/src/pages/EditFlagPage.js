import React from 'react';
import {Link, Redirect, useParams} from "react-router-dom";
import {reject} from "lodash";
import Content from "../theme/Content";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Divider,
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  TextField,
  Tooltip,
  Typography,
  withStyles
} from "@material-ui/core";
import {Delete as DeleteIcon, RemoveCircleOutline as RemoveIcon} from "@material-ui/icons";
import Grid from "@material-ui/core/Grid";
import {useMutation, useQuery} from "@apollo/react-hooks";
import {DELETE_FLAG_QUERY, FLAG_QUERY, FLAGS_QUERY} from "../Queries";
import {BooleanType, Operations, OperationTypes, VariantType, VariantTypes} from "./copy";

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
    margin: theme.spacing(2),
  },
  textField: {
    // marginLeft: theme.spacing(1),
    // marginRight: theme.spacing(1),
  },
  section1: {
    margin: theme.spacing(0, 0, 2, 0),
  },
  section2: {
    margin: theme.spacing(0),
  },
  section3: {
    margin: theme.spacing(1, 2, 1, 1),
  },
  footer: {
    marginTop: theme.spacing(1),
    paddingTop: theme.spacing(1),
    borderTop: 'solid #E4E4E4 1px',
  },
  footerDiv: {
    flexGrow: 1,
  },
  margin: {
    margin: theme.spacing(1),
  },
  formControl: {
    margin: 0,
    fullWidth: true,
    display: "flex",
    wrap: "nowrap",
  },
  paper: {
    margin: theme.spacing(0, 0, 2, 0),
    display: "flex",
    flexGrow: 1,
  },
});

const newVariant = (variant = {}) => ({
  id: variant.id || String(Math.random()),
  description: variant.description,
  value: variant.value,
  type: typeof variant.value,
  defaultWhenOn: false,
  defaultWhenOff: false,
});

const newRule = (rule = {}) => ({
  id: rule.id || String(Math.random()),
  constraints: rule.constraints || [newConstraint()],
  distributions: rule.distributions || [],
});

const newConstraint = (constraint = {}) => ({
  id: constraint.id || String(Math.random()),
  property: constraint.property || "",
  operation: constraint.operation || OperationTypes.ONE_OF,
  values: constraint.values || [""],
});

function VariantField({classes, variant: vrnt, handleDeleteVariant, handleUpdateVariant}) {
  const [variant, setVariant] = React.useState(vrnt);
  const setVariantField = (field, value) => {
    handleUpdateVariant(field, value);
    setVariant({...variant, [field]: value});
  };
  return (
    <Grid container spacing={2} className={classes.section2} key={variant.id}>
      <Grid item xs={3}>
        <FormControl className={classes.formControl}>
          <InputLabel>Variant type</InputLabel>
          <Select
            value={variant.type}
            onChange={e => setVariantField('type', e.target.value)}
          >
            {Object.keys(VariantTypes).map(type => (
              <MenuItem key={type} value={VariantTypes[type]}>{VariantType[type]}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      <Grid item xs={4}>
        {
          variant.type === VariantTypes.BOOLEAN ?
            (
              <FormControl className={classes.formControl}>
                <InputLabel>Variant value</InputLabel>
                <Select
                  value={variant.value}
                  onChange={e => setVariantField('value', e.target.value)}
                >
                  {[true, false].map(val => (
                    <MenuItem key={val} value={val}>{BooleanType[val]}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            ) :
            (
              <TextField
                label="Value"
                className={classes.textField}
                value={variant.value}
                type={variant.type === VariantTypes.NUMBER ? "number" : "text"}
                onChange={e => setVariantField('value', e.target.value)}
                fullWidth
              />
            )
        }
      </Grid>
      <Grid item xs={5} style={{display: 'flex'}}>
        <TextField
          label="Name"
          className={classes.textField}
          value={variant.description || ''}
          onChange={e => setVariantField('description', e.target.value)}
          fullWidth
        />
        <Tooltip title="Delete variant" placement="top">
          <Button size="small" color="secondary" style={{minWidth: 0}} onClick={() => handleDeleteVariant(vrnt)}>
            <RemoveIcon/>
          </Button>
        </Tooltip>
      </Grid>
    </Grid>
  )
}

function ConstraintField({classes, constraint: cnstrnt, operations, handleDeleteConstraint, handleUpdateConstraint}) {
  const [constraint, setConstraint] = React.useState(cnstrnt);
  const setConstraintField = (field, value) => {
    // handleUpdateConstraint(field, value);
    setConstraint({...constraint, [field]: value});
  };
  return (
    <Grid container spacing={2} className={classes.section2} key={constraint.id}>
      <Grid item xs={4}>
        <TextField
          label="Property"
          className={classes.textField}
          value={constraint.property}
          onChange={e => setConstraintField('property', e.target.value)}
          fullWidth
        />
      </Grid>
      <Grid item xs={4}>
        <FormControl className={classes.formControl}>
          <InputLabel>Operation</InputLabel>
          <Select
            value={constraint.operation}
            onChange={e => setConstraintField('operation', e.target.value)}
          >
            {operations.map(operation => (
              <MenuItem key={operation} value={operation.name}>{Operations[operation.name] || operation.name}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      <Grid item xs={4} style={{display: 'flex'}}>
        <TextField
          label="Value"
          className={classes.textField}
          value={constraint.values[0]}
          onChange={e => setConstraintField('values', [e.target.value])}
          fullWidth
        />
        <Tooltip title="Delete constraint" placement="top">
          <Button size="small" color="secondary" style={{minWidth: 0}} onClick={() => handleDeleteConstraint(cnstrnt)}>
            <RemoveIcon/>
          </Button>
        </Tooltip>
      </Grid>
    </Grid>
  );
}

function FlagForm({classes, flag: flg, operations, handleDeleteFlag}) {
  const [flag, setFlag] = React.useState(flg);
  const [deleteFlagDlgOpen, setDeleteFlagDlgOpen] = React.useState(false);
  const [defaults, setDefaults] = React.useState(
    flag.variants.reduce((dflts, variant) => {
      if (variant.defaultWhenOn) return {...dflts, defaultWhenOn: variant};
      if (variant.defaultWhenOff) return {...dflts, defaultWhenOff: variant};
      return dflts;
    }, {})
  );
  const setFlagField = (field, value) => setFlag({...flag, [field]: value});
  const handleClickOpen = () => setDeleteFlagDlgOpen(true);
  const handleClose = () => setDeleteFlagDlgOpen(false);
  const confirmDeleteFlag = () => {
    handleDeleteFlag(flag.id);
    setDeleteFlagDlgOpen(false);
  };
  const addVariant = () => {
    // TODO: make the variant type smarter (check existing variants)
    setFlag({...flag, variants: [...flag.variants, newVariant()]});
  };
  const updateVariant = variant => (field, value) => variant[field] = value;
  const deleteVariant = ({id}) => {
    setFlag({...flag, variants: reject(flag.variants, {id})});
  };
  const addConstraint = ruleId => () => {
    setFlag({
      ...flag, rules: flag.rules.map(rule => {
        if (rule.id === ruleId) rule = {...rule, constraints: [...rule.constraints, newConstraint()]};
        return rule;
      })
    });
  };
  const deleteConstraint = ruleId => ({id}) => {
    setFlag({
      ...flag, rules: flag.rules.map(rule => {
        if (rule.id === ruleId) rule = {...rule, constraints: reject(rule.constraints, {id})};
        return rule;
      })
    });
  };
  const addRule = () => {
    setFlag({...flag, rules: [...flag.rules, newRule()]});
  };
  const deleteRule = ({id}) => {
    setFlag({...flag, rules: reject(flag.rules, {id})});
  };

  return (
    <form className={classes.container} noValidate autoComplete="off">
      <Grid container spacing={2} className={classes.section1}>
        <Grid item xs={12} sm={6}>
          <TextField
            label="Name"
            className={classes.textField}
            value={flag.name}
            onChange={e => setFlagField('name', e.target.value)}
            fullWidth
          />
        </Grid>
        <Grid item xs={12} sm={6}>
          <TextField
            label="Key"
            className={classes.textField}
            value={flag.key}
            onChange={e => setFlagField('key', e.target.value)}
            fullWidth
          />
        </Grid>
        <Grid item xs={12} sm={12}>
          <TextField
            label="Description"
            className={classes.textField}
            value={flag.description || ''}
            onChange={e => setFlagField('description', e.target.value)}
            fullWidth
          />
        </Grid>
      </Grid>

      <Grid container spacing={2} className={classes.section1}>
        <Grid item xs={12}>
          <Typography variant="h6" color="textSecondary">Variants</Typography>
          <Divider light/>
        </Grid>
      </Grid>

      <Grid container className={classes.section1}>
        {
          flag.variants.map(variant => (
            <VariantField key={variant.id} classes={classes} variant={newVariant(variant)}
                          handleUpdateVariant={updateVariant(variant)} handleDeleteVariant={deleteVariant}/>
          ))
        }
        <Button variant="outlined" size="small" color="primary" onClick={addVariant} className={classes.margin}>
          New Variant
        </Button>
      </Grid>

      <Grid container spacing={2} className={classes.section1}>
        <Grid item xs={6}>
          <FormControl className={classes.formControl}>
            <InputLabel>Default value when flag is enabled</InputLabel>
            <Select
              value={defaults.defaultWhenOn}
              onChange={e => setDefaults({...defaults, defaultWhenOn: e.target.value})}
            >
              {flag.variants.map(variant => (
                <MenuItem key={variant.id} value={variant}>
                  {typeof variant.value === VariantTypes.BOOLEAN ? BooleanType[variant.value] : variant.value}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={6}>
          <FormControl className={classes.formControl}>
            <InputLabel>Default value when flag is disabled</InputLabel>
            <Select
              value={defaults.defaultWhenOff}
              onChange={e => setDefaults({...defaults, defaultWhenOff: e.target.value})}
            >
              {flag.variants.map(variant => (
                <MenuItem key={variant.id} value={variant}>
                  {typeof variant.value === VariantTypes.BOOLEAN ? BooleanType[variant.value] : variant.value}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
      </Grid>

      <Grid container spacing={2} className={classes.section1}>
        <Grid item xs={12}>
          <Typography variant="h6" color="textSecondary">Rules</Typography>
          <Divider light/>
        </Grid>
      </Grid>

      <Grid container className={classes.section1}>
        <Grid item xs={12}>
          {
            flag.rules.map(rule => (
              <Paper key={rule.id} className={classes.paper}>
                <Grid container className={classes.section3}>
                  <Grid item xs={12}>
                    {rule.constraints.map(constraint => (
                      <ConstraintField key={constraint.id} classes={classes} constraint={constraint}
                                       operations={operations} handleDeleteConstraint={deleteConstraint(rule.id)}/>
                    ))}
                  </Grid>
                  <Button variant="outlined" size="small" color="primary" onClick={addConstraint(rule.id)}
                          className={classes.margin}>
                    New Constraint
                  </Button>
                  <Button variant="outlined" size="small" color="secondary" onClick={() => deleteRule(rule)}
                          className={classes.margin}>
                    Delete Rule
                  </Button>
                </Grid>
              </Paper>
            ))
          }
        </Grid>
        <Button variant="outlined" size="small" color="primary" onClick={addRule} className={classes.margin}>
          New Rule
        </Button>
      </Grid>

      <Grid container alignContent="space-around" direction="row-reverse" className={classes.footer}>
        <Grid item>
          <Button color="primary" onClick={() => console.log(flag.variants)}>
            Save
          </Button>
        </Grid>
        <Grid item>
          <Link to="/flags">
            <Button>
              Cancel
            </Button>
          </Link>
        </Grid>
        <Grid item style={{flexGrow: 1}}>
          <DeleteFlagDialog open={deleteFlagDlgOpen} onConfirm={confirmDeleteFlag} handleClose={handleClose}
                            flag={flag}/>
          <Tooltip title="Delete flag" placement="top">
            <Button color="secondary" onClick={handleClickOpen}>
              <DeleteIcon/>
            </Button>
          </Tooltip>
        </Grid>
      </Grid>
    </form>
  )
}

function DeleteFlagDialog({open, flag, onConfirm, handleClose}) {
  return (
    <Dialog
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">Delete flag?</DialogTitle>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          Are you sure you want to delete flag "{flag.name}"?
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
          No, keep it
        </Button>
        <Button onClick={onConfirm} color="secondary" autoFocus>
          Yes, delete it
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function EditFlagPage({classes}) {
  const [toFlagsPage, setToFlagsPage] = React.useState(false);
  let {id} = useParams();
  const {loading, error, data} = useQuery(FLAG_QUERY, {variables: {id}});
  const [deleteFlag] = useMutation(DELETE_FLAG_QUERY, {
    update(cache, {data: {deleteFlag: id}}) {
      const {flags} = cache.readQuery({query: FLAGS_QUERY});
      cache.writeQuery({
        query: FLAGS_QUERY,
        data: {flags: reject(flags, {id})},
      });
    }
  });
  if (toFlagsPage === true) {
    return <Redirect to='/flags'/>
  }
  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>"Error while loading flag details :("</div>;
  const handleDeleteFlag = (id) => {
    deleteFlag({variables: {id}}).then(() => setToFlagsPage(true));
  };

  return (
    <Content>
      <FlagForm classes={classes}
                flag={data.flag}
                operations={data.operations.enumValues}
                handleDeleteFlag={handleDeleteFlag}/>
    </Content>
  )
}

export default withStyles(styles)(EditFlagPage);