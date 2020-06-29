import React, { useEffect, useState } from 'react';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import { Link } from 'react-router-dom';
import axios from "axios";

const styles = {
    root: {
        flexGrow: 1,
    },
    title: {
        flexGrow: 1,
    },
    appbar: {
        backgroundColor: '#ED1C16',
    },
};

function Landing(props) {
    const { classes } = props;
    const initialUserState = {
        loggedIn: false,
        session: {
            AccountID: 0,
            LastActivityTime: "",
        },
    };
    const [session, setSession] = useState(initialUserState)

    useEffect(() => {
        const getSession = async () => {
            axios.get(
                "http://localhost:3001/v1/sessions",
                { withCredentials: true }
            ).then(response => {
                if (response.error) {
                    return
                }
                if (response.status === 200) {
                    setSession({
                        loggedIn: true,
                        session: response.data,
                    });
                }
            });
        }
        getSession()
    }, [])


    return (
        <React.Fragment>
            <div className={classes.root}>
                <AppBar position="static" className={classes.appbar}>
                    <Toolbar>
                        <Typography variant="h6" className={classes.title}>Hackerspace Membership Management</Typography>
                        {
                            session.loggedIn ?
                                <div>
                                    <Button component={Link} to={'/me'} color="inherit">My account</Button>
                                    <Button component={Link} to={'/logout'} color="inherit">Logout</Button>
                                </div> :
                                <div>
                                    <Button component={Link} to={'/login'} color="inherit">Login</Button>
                                    <Button component={Link} to={'/signup'} color="inherit">Sign Up</Button>
                                </div>
                        }
                    </Toolbar>
                </AppBar>
            </div>
        </React.Fragment >
    );
}

export default withStyles(styles)(Landing);