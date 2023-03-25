import Landing from "./components/Landing/Landing";
import { BrowserRouter, Switch, Route } from "react-router-dom";
import "./App.css";

function App() {
  return (
    <BrowserRouter>
      <Switch>
        <Route exact path={"/"} render={(props) => <Landing />} />
      </Switch>
    </BrowserRouter>
  );
}

export default App;
