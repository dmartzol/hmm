import React, { Component } from "react";
import { Fragment } from "react";
import { withRouter } from 'react-router-dom';

class Signup extends Component {

    render() {
        return (
            <Fragment>
                <div className="container">
                    <div class="py-5 text-center">
                        <img class="d-block mx-auto mb-4" src="#" alt="" width="72" height="72" />
                        <h2>Sign up form</h2>
                        <p class="lead">Basic info now and you will be able to configure your membership later.</p>
                    </div>
                    <div className="row">
                        <div className="col-md-12">
                            <h4 class="mb-3">Info</h4>
                            <form action="" className="needs-validation" novalidate>
                                <div class="row">
                                    <div class="col-md-4 mb-3">
                                        <label for="firstName">First name</label>
                                        <input type="text" class="form-control" id="firstName" placeholder="" value="" required />
                                        <div class="invalid-feedback">Valid first name is required.</div>
                                    </div>
                                    <div class="col-md-4 mb-3">
                                        <label for="firstName">Middle name</label>
                                        <input type="text" class="form-control" id="firstName" placeholder="" value="" required />
                                        <div class="invalid-feedback">Valid first name is required.</div>
                                    </div>
                                    <div class="col-md-4 mb-3">
                                        <label for="lastName">Last name</label>
                                        <input type="text" class="form-control" id="lastName" placeholder="" value="" required />
                                        <div class="invalid-feedback">Valid last name is required.</div>
                                    </div>
                                </div>

                                <div className="row">
                                    <div class="col-md-4 mb-3">
                                        <label for="email">Email</label>
                                        <input type="email" class="form-control" id="email" placeholder="you@example.com" />
                                        <div class="invalid-feedback">Please enter a valid email address for shipping updates.</div>
                                    </div>
                                    <div class="col-md-4 mb-3">
                                        <label for="email">Password</label>
                                        <input type="password" class="form-control" id="password" placeholder="Password" />
                                        <div class="invalid-feedback">Please enter a valid email address for shipping updates.</div>
                                    </div>
                                    <div class="col-md-4 mb-3">
                                        <label for="email">Repeat Password</label>
                                        <input type="password" class="form-control" id="password" placeholder="Repeat Password" />
                                        <div class="invalid-feedback">Please enter a valid email address for shipping updates.</div>
                                    </div>
                                </div>

                                <div class="mb-3">
                                    <label for="address">Address</label>
                                    <input type="text" class="form-control" id="address" placeholder="1234 Main St" required />
                                    <div class="invalid-feedback">Please enter your shipping address.</div>
                                </div>
                                <div class="mb-3">
                                    <label for="address2">Address 2 <span class="text-muted">(Optional)</span></label>
                                    <input type="text" class="form-control" id="address2" placeholder="Apartment or suite" />
                                </div>

                                <div class="row">
                                    <div class="col-md-5 mb-3">
                                        <label for="country">Country</label>
                                        <select class="custom-select d-block w-100" id="country" required>
                                            <option value="">Choose...</option>
                                            <option>United States</option>
                                        </select>
                                        <div class="invalid-feedback">Please select a valid country.</div>
                                    </div>
                                    <div class="col-md-4 mb-3">
                                        <label for="state">State</label>
                                        <select class="custom-select d-block w-100" id="state" required>
                                            <option value="">Choose...</option>
                                            <option>California</option>
                                            <option>California</option>
                                            <option>California</option>
                                        </select>
                                        <div class="invalid-feedback">Please provide a valid state.</div>
                                    </div>
                                    <div class="col-md-3 mb-3">
                                        <label for="zip">Zip</label>
                                        <input type="text" class="form-control" id="zip" placeholder="" required />
                                        <div class="invalid-feedback">Zip code required.</div>
                                    </div>
                                </div>
                                <hr class="mb-4" />
                                <button class="btn btn-primary btn-lg btn-block" type="submit">Submit</button>
                            </form>
                        </div>
                    </div>
                </div>
            </Fragment>
        );
    }
}

export default withRouter(Signup);