import React, { Component } from "react";
import Landing from "./Landing/Landing";
import { BrowserRouter, Switch, Route } from 'react-router-dom';
import Login from './Landing/Login';
import Logout from './Landing/Logout';
import Paperbase from './Dashboard/Paperbase';

export default class App extends Component {

  render() {
    return (
      <BrowserRouter>
      <Switch>
        <Route exact path={"/"} render={props => (<Landing />)} />
        <Route exact path={"/login"} render={props => (<Login />)} />
        <Route exact path={"/logout"} render={props => (<Logout />)} />
        <Route exact path={"/signup"} render={props => (<Logout />)} />
        <Route exact path={"/me"} render={props => (<Paperbase />)} />
        <Route exact path={"/accounts"} render={props => (<Paperbase />)} />
        <Route exact path={"/authorizations"} render={props => (<Paperbase />)} />
        <Route exact path={"/equipment"} render={props => (<Paperbase />)} />
        <Route exact path={"/roles"} render={props => (<Paperbase />)} />
      </Switch>
  </BrowserRouter>
    );
  }
}