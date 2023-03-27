import React from 'react';

function Footer() {

  return (
    <footer class="footer mt-5 py-3">
      <div className="container">
        <div className="container-fluid text-center text-secondary">
          Hmm! © {new Date().getFullYear()}{'.'}
        </div>
      </div>
    </footer>
  );
}

export default Footer;