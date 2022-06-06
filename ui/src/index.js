import React from "react";
import ReactDOM from "react-dom/client";
import 'mapbox-gl/dist/mapbox-gl.css';
import "./index.css";
import reportWebVitals from "./reportWebVitals";
import { ApolloClient, ApolloProvider, InMemoryCache } from "@apollo/client";
import { HashRouter, Route, Routes } from "react-router-dom";
import Home from "./Containers/Home";
import Region from "./Containers/Region";
import District from "./Containers/District";
import Commune from "./Containers/Commune";
import Fokontany from "./Containers/Fokontany";

const client = new ApolloClient({
  // eslint-disable-next-line no-restricted-globals
  uri: location.origin + "/query",
  // uri: 'http://localhost:8010/query',
  cache: new InMemoryCache(),
});

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <ApolloProvider client={client}>
    <React.StrictMode>
      <HashRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/regions/:id" element={<Region />} />
          <Route path="/districts/:id" element={<District />} />
          <Route path="/communes/:id" element={<Commune />} />
          <Route path="/fokontany/:id" element={<Fokontany />} />
        </Routes>
      </HashRouter>
    </React.StrictMode>
  </ApolloProvider>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
