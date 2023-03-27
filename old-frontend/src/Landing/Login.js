import React, { Component } from "react";
import { withRouter } from 'react-router-dom';
import axios from "axios";

class Login extends Component {
  constructor(props) {
    super(props);

    this.state = {
      email: "",
      password: "",
      session: {},
    };

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(event) {
    this.setState({
      [event.target.name]: event.target.value
    });
  }

  handleSubmit(event) {
    const { email, password, } = this.state;

    axios
      .post(
        "http://localhost:3001/v1/sessions",
        {
          email: email,
          password: password
        },
        { withCredentials: true }
      )
      .then(response => {
        if (response.status === 200) {
          this.props.history.push('/me');
        }
      })
      .catch(error => {
        console.log("login error", error);
      });
    event.preventDefault();
  }

  render() {
    return (
      <div className="container-fluid bg-light d-flex justify-content-center align-items-center min-vh-100">
          <div className="p-0 col-xs-12 col-sm-8 col-lg-4 col-xl-3">
            <form className="px-3 py-5 bg-white card border-dark" onSubmit={this.handleSubmit}>
              <input
                className="form-control mb-2 p-1"
                type="email"
                name="email"
                placeholder="email@example.com"
                value={this.state.email}
                onChange={this.handleChange}
                required
              />
              <input
                className="form-control p-1"
                type="password"
                name="password"
                placeholder="Password"
                value={this.state.password}
                onChange={this.handleChange}
                required
              />
              <div className="mt-3">
                <button className="btn btn-lg btn-primary btn-block" type="submit">Log in</button>
              </div>
            </form>
          </div>
      </div>


    );
  }
}

export default withRouter(Login);