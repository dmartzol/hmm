import React from 'react';

function Footer() {

  return (
    <footer class="container-fluid mt-5">
      <div className="row">
        <div className="container-fluid text-center text-secondary">
          Hmm! © {new Date().getFullYear()}{'.'}
        </div>
      </div>
    </footer>
  );
}

export default Footer;
