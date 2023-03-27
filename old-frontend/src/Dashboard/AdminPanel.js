import React from 'react';
import Navigator from './Navigator';
import AuthorizationsIndexView from './AuthorizationsIndexView';
import AccountsIndexView from './AccountsIndexView';
import SingleAccountView from './SingleAccountView';
import MyAccount from './MyAccount';
import Header from './Header';
import Footer from './Footer';
import axios from "axios";
import { Route, Switch, withRouter } from 'react-router-dom';


function AdminPanel(props) {
  const initialStates = {
    loggedIn: false,
    session: {
      AccountID: 0,
      LastActivityTime: "",
    },
  };
  const [session, setSession] = React.useState(initialStates);

  React.useEffect(() => {
    const getSession = async () => {
      axios.get(
        "http://localhost:3001/v1/sessions",
        { withCredentials: true }
      ).then(response => {
        if (response.status === 200) {
          setSession({
            loggedIn: true,
            session: response.data,
          });
        }
      }).catch(error => {
        props.history.push('/');
        return
      });
    }
    getSession()
  }, [props])

  return (
    <div className="container-fluid">
      <div className="row">
        <div className="col-2 bg-dark min-vh-100">
          <Navigator />
        </div>
        <div className="col-10 bg-light p-0">
          <Header />
          <Switch>
            <Route exact path="/accounts" render={(props) => <AccountsIndexView {...props} session={session} isAuthed={session.loggedIn} />} />
            <Route exact path="/accounts/:id" render={(props) => <SingleAccountView {...props} session={session} isAuthed={session.loggedIn} />} />
            <Route exact path="/me" render={(props) => <MyAccount {...props} session={session} isAuthed={session.loggedIn} />} />
            <Route exact path="/authorizations" component={AuthorizationsIndexView} />
          </Switch>
          <Footer />
        </div>
      </div>
    </div>
  );
}

export default withRouter(AdminPanel);
