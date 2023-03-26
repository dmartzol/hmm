import React from "react";
import Footer from "../Footer/Footer";
import NavBar from "../NavBar/NavBar";

function PageLayout(props) {
  return (
    <>
      <NavBar />
      <main>{props.children}</main>
      <Footer />
    </>
  );
}

export default PageLayout;
