import React from 'react';
import {useParams} from "react-router-dom";
import Content from "../theme/Content";
import {Button, Checkbox, FormControlLabel, TextField, withStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import {useQuery} from "@apollo/react-hooks";
import {FLAG_QUERY} from "../Queries";

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
  footer: {
    marginTop: theme.spacing(1),
    paddingTop: theme.spacing(1),
    borderTop: 'solid #E4E4E4 1px',
  },
  footerDiv: {
    flexGrow: 1,
  }
});

function Variant() {

}

function EditFlagPage(props) {
  const {classes} = props;
  let {id} = useParams();
  const {loading, error, data} = useQuery(FLAG_QUERY, {variables: {id}});
  if (loading) return <div>"Loading..."</div>;
  if (error) return <div>"Error while loading flag details :("</div>;
  const {flag} = data;
  const setFlag = (field, value) => flag[field] = value;

  return (
    <Content>
      <form className={classes.container} noValidate autoComplete="off">
        <Grid container spacing={2} className={classes.section1}>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Name"
              className={classes.textField}
              value={flag.name}
              // onChange={e => setFlag('name', e.target.value)}
              fullWidth
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              label="Key"
              className={classes.textField}
              value={flag.key}
              // onChange={e => setFlag('key', e.target.value)}
              fullWidth
            />
          </Grid>
          <Grid item xs={12} sm={12}>
            <TextField
              label="Description"
              className={classes.textField}
              value={flag.description || ''}
              // onChange={e => setFlag('key', e.target.value)}
              fullWidth
            />
          </Grid>
        </Grid>

        {
          flag.variants.map(variant => (
            <Grid container spacing={2} className={classes.section2} key={variant.id}>
              <Grid item xs={6} sm={6}>
                <TextField
                  label="Value"
                  className={classes.textField}
                  value={variant.value}
                  // onChange={e => setFlag('name', e.target.value)}
                  fullWidth
                />
              </Grid>
              <Grid item xs={6} sm={6}>
                <TextField
                  label="Name"
                  className={classes.textField}
                  value={variant.name}
                  // onChange={e => setFlag('name', e.target.value)}
                  fullWidth
                />
              </Grid>
              <Grid item xs={4} sm={4}>
                <FormControlLabel
                  control={
                    <Checkbox
                      // checked={state.checkedB}
                      // onChange={handleChange('checkedB')}
                      color="primary"
                    />
                  }
                  label="Use when flag is on"
                />
              </Grid>
              <Grid item xs={4} sm={4}>
                <FormControlLabel
                  control={
                    <Checkbox
                      // checked={state.checkedB}
                      // onChange={handleChange('checkedB')}
                      color="primary"
                    />
                  }
                  label="Use when flag is off"
                />
              </Grid>
            </Grid>
          ))
        }
        <Grid container alignContent="space-around" direction="row-reverse" className={classes.footer}>
          <Grid item>
            <Button color="primary">
              Save
            </Button>
          </Grid>
          <Grid item>
            <Button color="secondary">
              Cancel
            </Button>
          </Grid>
        </Grid>
      </form>
    </Content>
  )
}

export default withStyles(styles)(EditFlagPage);