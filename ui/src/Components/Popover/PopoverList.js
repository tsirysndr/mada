import React, { Component } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import _ from "lodash";

function withRouter(Child) {
  return (props) => {
    const location = useLocation();
    const navigate = useNavigate();
    return <Child {...props} navigate={navigate} location={location} />;
  };
}

class PopoverList extends Component {
  componentDidMount() {
    this.refs.iScroll.addEventListener("scroll", () => {
      if (
        this.refs.iScroll.scrollTop + this.refs.iScroll.clientHeight >=
        this.refs.iScroll.scrollHeight
      ) {
        this.props.handleUpdate();
      }
    });
  }

  render() {
    const { keyword, filter, loading, error, search, navigate } = this.props;
    const fokontany =
      keyword === ""
        ? this.props.fokontany
        : _.defaultTo(search.hits, [])
            .filter((x) => x.fields.type === "fokontany")
            .map((x) => ({ id: x.id, name: x.fields.fokontany }));
    const communes =
      keyword === ""
        ? this.props.communes
        : _.defaultTo(search.hits, [])
            .filter((x) => x.fields.type === "commune")
            .map((x) => ({ id: x.id, name: x.fields.commune }));
    const districts =
      keyword === ""
        ? this.props.districts
        : _.defaultTo(search.hits, [])
            .filter((x) => x.fields.type === "district")
            .map((x) => ({ id: x.id, name: x.fields.district }));
    const regions =
      keyword === ""
        ? this.props.regions
        : _.defaultTo(search.hits, [])
            .filter((x) => x.fields.type === "region")
            .map((x) => ({ id: x.id, name: x.fields.region }));
    return (
      <div ref="iScroll" style={{ overflow: "auto" }}>
        {filter === 1 && !loading && !error && (
          <ul className="popover-list">
            <li>
              {regions.map((item, index) => (
                <a href={`#/regions/${item.id}`} className="item" key={index}>
                  {item.name}
                </a>
              ))}
            </li>
          </ul>
        )}
        {filter === 2 && !loading && !error && (
          <ul className="popover-list">
            <li>
              {districts.map((item, index) => (
                <div
                  className="item"
                  key={index}
                  onClick={() => navigate(`/districts/${item.id}`)}
                >
                  <a href={`#/districts/${item.id}`}>{item.name}</a>
                  <div>{item.region}</div>
                </div>
              ))}
            </li>
          </ul>
        )}
        {filter === 3 && !loading && !error && (
          <ul className="popover-list">
            <li>
              {communes.map((item, index) => (
                <div
                  className="item"
                  key={index}
                  onClick={() => navigate(`/communes/${item.id}`)}
                >
                  <a href={`#/communes/${item.id}`}>{item.name}</a>
                  <div>
                    {item.district} &middot; {item.region}
                  </div>
                </div>
              ))}
            </li>
          </ul>
        )}
        {filter === 4 && !loading && !error && (
          <ul className="popover-list">
            <li>
              {fokontany.map((item, index) => (
                <div
                  className="item"
                  key={index}
                  onClick={() => navigate(`/fokontany/${item.id}`)}
                >
                  <a href={`#/fokontany/${item.id}`}>{item.name}</a>
                  <div>
                    {item.commune} &middot; {item.district} &middot;{" "}
                    {item.region}
                  </div>
                </div>
              ))}
            </li>
          </ul>
        )}
      </div>
    );
  }
}

export default withRouter(PopoverList);
