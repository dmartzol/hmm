import React from 'react';
import Button from '@material-ui/core/Button';
import { Link } from 'react-router-dom';

function Header() {

  return (
    <nav class="navbar navbar-expand-lg navbar-light bg-light mb-5">
      <div class="navbar-collapse justify-content-end">
        <form class="form-inline my-2 my-lg-0">
          <input class="form-control mr-sm-2" type="search" placeholder="Search" aria-label="Search" />
          <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
        </form>
      </div>
    </nav>
  );
}

export default Header;
