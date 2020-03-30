import React from 'react';
import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import Grow from '@material-ui/core/Grow';
import Paper from '@material-ui/core/Paper';
import Popper from '@material-ui/core/Popper';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';
import { makeStyles } from '@material-ui/styles';
import { Link } from 'react-router-dom';
import { VariantTypes } from '../../../FlagForm/models';

const options = [
  { name: 'Add boolean flag', to: { pathname: '/flags/new', flagType: VariantTypes.BOOLEAN } },
  { name: 'Add numeric flag', to: { pathname: '/flags/new', flagType: VariantTypes.NUMBER } },
  { name: 'Add string flag', to: { pathname: '/flags/new', flagType: VariantTypes.STRING } },
];

const useStyles = makeStyles(theme => ({
  popper: {
    zIndex: '9999',
  },
  mainButton: {
    color: `${theme.palette.primary.contrastText} !important`,
  },
}));

export default function AddFlagButton() {
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef(null);
  const classes = useStyles();

  const handleToggle = () => {
    setOpen(prevOpen => !prevOpen);
  };

  const handleClose = event => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }

    setOpen(false);
  };

  return (
    <>
      <ButtonGroup variant="contained" color="primary" ref={anchorRef} aria-label="split button">
        <Button
          className={classes.mainButton}
          color="primary"
          variant="contained"
          component={Link}
          to={options[0].to}
        >
          Add flag
        </Button>
        <Button
          color="primary"
          size="small"
          aria-controls={open ? 'split-button-menu' : undefined}
          aria-expanded={open ? 'true' : undefined}
          aria-label="add flag"
          aria-haspopup="menu"
          onClick={handleToggle}
        >
          <ArrowDropDownIcon/>
        </Button>
      </ButtonGroup>
      <Popper className={classes.popper} open={open} anchorEl={anchorRef.current} role={undefined} transition
              disablePortal>
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin: placement === 'bottom' ? 'center top' : 'center bottom',
            }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList id="split-button-menu">
                  {options.map((option, index) => (
                    <MenuItem
                      key={option.name}
                      component={Link}
                      to={option.to}
                    >
                      {option.name}
                    </MenuItem>
                  ))}
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
}