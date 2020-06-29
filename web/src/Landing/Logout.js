import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import axios from "axios";

class Logout extends Component {

  componentDidMount(event) {
    axios.delete(
      "http://localhost:3001/v1/sessions",
      { withCredentials: true }
    ).then(response => {
      if (response.status === 200) {
        this.props.history.push('/');
      }
    })
  }

  render() {
    return (
      <div></div>
    );
  }
}

export default withRouter(Logout);