import React from 'react';


function Navigator() {

  return (
    <div className="container bg-dark min-vh-100">
      <div className="row">
        <ul className="list-group mt-5 p-0">
          <a activeClassName="active" className="list-group-item mt-0 py-2 btn-secondary" href="/me"><i className="material-icons mr-3">home</i>My Account</a>
          <hr />
          <a activeClassName="active" className="list-group-item mt-0 py-2 btn-secondary" href="/accounts"><i className="material-icons mr-3">people</i>Accounts</a>
          <a activeClassName="active" className="list-group-item mt-0 py-2 btn-secondary" href="/authorizations"><i className="material-icons mr-3">fingerprint</i>Authorizations</a>
          <a activeClassName="active" className="list-group-item mt-0 py-2 btn-secondary" href="/equipment"><i className="material-icons mr-3">build</i>Equipment</a>
          <a activeClassName="active" className="list-group-item mt-0 py-2 btn-secondary" href="/roles"><i className="material-icons mr-3">category</i>Roles</a>
          <hr />
          <a className="list-group-item mt-0 btn-secondary" href="/logout"><i className="material-icons mr-3">exit_to_app</i>Logout</a>
        </ul>
      </div>
    </div>
  );
}

export default Navigator;

