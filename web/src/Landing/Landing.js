import React, { useEffect, useState } from 'react';
import Button from '@material-ui/core/Button';
import { Link } from 'react-router-dom';
import axios from "axios";

function Landing() {
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

        <nav class="navbar navbar-expand-lg navbar-light bg-light">
            <a class="navbar-brand" href="/">Hmm!</a>
            <div class="collapse navbar-collapse" id="navbarSupportedContent">
                <ul class="navbar-nav mr-auto">
                    <li>
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
                    </li>
                </ul>
                <form class="form-inline my-2 my-lg-0">
                    <input class="form-control mr-sm-2" type="search" placeholder="Search" aria-label="Search" />
                    <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
                </form>
            </div>
        </nav>
    );
}

export default Landing;