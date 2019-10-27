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
  Grid,
  Tab,
  Tabs,
  TextField,
  Tooltip,
} from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';
import { Link } from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  root: {},
  actionButton: {
    margin: theme.spacing(1),
  },
}));

const FlagDetails = props => {
  const {className, flag: flg, ...rest} = props;
  const [flag, setFlag] = useState(flg);
  const [tab, setTab] = React.useState(0);
  const classes = useStyles();

  const handleChange = event => {
    setFlag({
      ...flag,
      [event.target.name]: event.target.value,
    });
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
          aria-label="full width tabs example"
        >
          <Tab label="General"/>
          <Tab label="Rules"/>
          <Tab label="Evaluation" disabled/>
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
                  onChange={handleChange}
                  required
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
                  onChange={handleChange}
                  required
                  value={flag.key}
                  variant="outlined"
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Description"
                  margin="dense"
                  name="description"
                  onChange={handleChange}
                  value={flag.description}
                  variant="outlined"
                />
              </Grid>
            </Grid>
          </CardContent>
          <CardHeader
            subheader="Possible values that this flag can return"
            title="Variants"
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
                  onChange={handleChange}
                  required
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
                  onChange={handleChange}
                  required
                  value={flag.key}
                  variant="outlined"
                />
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
            <Grid container spacing={3}>
              <Grid item md={6} xs={12}>
                <TextField
                  fullWidth
                  label="Name"
                  margin="dense"
                  name="name"
                  onChange={handleChange}
                  required
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
                  onChange={handleChange}
                  required
                  value={flag.key}
                  variant="outlined"
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Description"
                  margin="dense"
                  name="description"
                  onChange={handleChange}
                  value={flag.description}
                  variant="outlined"
                />
              </Grid>
            </Grid>
          </CardContent>
        </Box>

        <Divider/>
        <CardActions>
          <Tooltip title="Delete flag" placement="top">
            <Button color="secondary">
              <DeleteIcon/>
            </Button>
          </Tooltip>
          <div style={{flexGrow: 1}}/>
          <Link to="/flags">
            <Button className={classes.actionButton}>Cancel</Button>
          </Link>
          <Button color="primary" variant="outlined" className={classes.actionButton}>
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
};

export default FlagDetails;
