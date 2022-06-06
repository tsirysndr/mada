import React, { useState, useEffect } from "react";
import { useQuery, gql } from "@apollo/client";
import { Tag } from "antd";
import MDSpinner from "react-md-spinner";
import PopoverList from "./PopoverList";
import _ from "lodash";

const ALL_FOKONTANY = gql`
  query AllFokontany($skip: Int, $size: Int) {
    allFokontany(skip: $skip, size: $size) {
      data {
        id
        name
        commune
        district
        region
      }
    }
  }
`;

const COMMUNES = gql`
  query Communes($skip: Int, $size: Int) {
    communes(skip: $skip, size: $size) {
      data {
        id
        name
        district
        region
      }
      after {
        id
      }
    }
  }
`;

const DISTRICTS = gql`
  query Districts($skip: Int, $size: Int) {
    districts(skip: $skip, size: $size) {
      data {
        id
        name
        region
      }
    }
  }
`;

const REGIONS = gql`
  query Regions($skip: Int, $size: Int) {
    regions(skip: $skip, size: $size) {
      data {
        id
        name
      }
    }
  }
`;

const SEARCH = gql`
  query Search($keyword: String!) {
    search(keyword: $keyword) {
      fokontany {
        id
        name
        commune
        district
        region
        country
        point
        coordinates
      }
      commune {
        id
        name
        district
        region
        country
        point
        coordinates
      }
      district {
        id
        name
        region
        country
        point
        coordinates
      }
      region {
        id
        name
        country
        point
        coordinates
      }
      hits {
        id
        score
        fields {
          fokontany
          commune
          district
          region
          country
          type
        }
      }
    }
    countRegions
    countDistricts
    countCommunes
    countFokontany
  }
`;

/*
const COUNT = gql`
  query {
    countRegions
    countDistricts
    countCommunes
    countFokontany
  }
`
*/

const Popover = (props) => {
  const [filter, setFilter] = useState(1);
  const [keyword, setKeyword] = useState("");
  const [fokontany, setFokontany] = useState([]);
  const [communes, setCommunes] = useState([]);
  const [districts, setDistricts] = useState([]);
  const [regions, setRegions] = useState([]);
  const [fokontanySkip, setFokontanySkip] = useState(0);
  const [communeSkip, setCommuneSkip] = useState(0);
  const [districtSkip, setDistrictSkip] = useState(0);
  const { loading, error, data } = useQuery(SEARCH, { variables: { keyword } });
  const search = _.get(data, "search", {});
  console.log("data", search, loading, error, keyword);
  const allFokontanyRes = useQuery(ALL_FOKONTANY, {
    variables: { skip: fokontanySkip, size: 100 },
  });
  const communesRes = useQuery(COMMUNES, {
    variables: { skip: communeSkip, size: 100 },
  });
  const districtsRes = useQuery(DISTRICTS, {
    variables: { skip: districtSkip, size: 100 },
  });
  const regionsRes = useQuery(REGIONS, { variables: { skip: 0, size: 100 } });

  const handleUpdate = () => {
    switch (filter) {
      case 1:
        break;
      case 2:
        setDistrictSkip(districtSkip + 100);
        break;
      case 3:
        setCommuneSkip(communeSkip + 100);
        break;
      case 4:
        setFokontanySkip(fokontanySkip + 100);
        break;
      default:
        break;
    }
  };

  useEffect(() => {
    if (allFokontanyRes.data) {
      const newFokontany = !fokontanySkip
        ? allFokontanyRes.data.allFokontany.data
        : [...fokontany, ...allFokontanyRes.data.allFokontany.data];
      setFokontany(newFokontany);
    }
  }, [allFokontanyRes.data]);

  useEffect(() => {
    if (communesRes.data) {
      const newCommunes = !communeSkip
        ? communesRes.data.communes.data
        : [...communes, ...communesRes.data.communes.data];
      setCommunes(newCommunes);
    }
  }, [communesRes.data]);

  useEffect(() => {
    if (districtsRes.data) {
      const newDistricts = !districtSkip
        ? districtsRes.data.districts.data
        : [...districts, ...districtsRes.data.districts.data];
      setDistricts(newDistricts);
    }
  }, [districtsRes.data]);

  useEffect(() => {
    if (regionsRes.data) {
      setRegions(regionsRes.data.regions.data);
    }
  }, [regionsRes.data]);

  return (
    <div id="search-popover" className={props.popoverClass}>
      <div
        style={{ display: "flex", height: 64 }}
        onClick={() => props.setExpanded(true)}
      >
        <input
          id="search"
          placeholder="Search ..."
          style={{ marginLeft: 50 }}
          autoComplete="off"
          onChange={(evt) => {
            if (evt.target.value.length > 2 || evt.target.value.length === 0) {
              setKeyword(evt.target.value);
            }
          }}
        />
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 14 14"
          className="search-icon"
          style={{ position: "absolute", top: 20, left: 15 }}
        >
          <g fill="none" fillRule="evenodd">
            <path d="M-5-5v24h24V-5z" />
            <path
              fill="#5C5C6E"
              fillRule="nonzero"
              d="M12.003 13.657L8.96 10.616c-.016-.016-.027-.035-.042-.052a5.756 5.756 0 1 1 1.645-1.645c.017.015.036.026.052.042l3.041 3.042a1.17 1.17 0 0 1-1.654 1.654zM9.516 5.756a3.76 3.76 0 1 0-7.52 0 3.76 3.76 0 0 0 7.52 0z"
            />
          </g>
        </svg>
      </div>
      {(loading || error) && (
        <div id="loader">
          <MDSpinner />
        </div>
      )}
      {!loading && !error && (
        <div style={{ padding: 10 }}>
          <Tag
            color="cyan"
            className={`tag ${filter !== 1 ? "inactive" : ""}`}
            onClick={() => setFilter(1)}
          >
            Regions (
            {keyword === ""
              ? data.countRegions
              : _.defaultTo(search.hits, []).filter(
                  (x) => x.fields.type === "region"
                ).length}
            )
          </Tag>
          <Tag
            color="cyan"
            className={`tag ${filter !== 2 ? "inactive" : ""}`}
            onClick={() => setFilter(2)}
          >
            Districts (
            {keyword === ""
              ? data.countDistricts
              : _.defaultTo(search.hits, []).filter(
                  (x) => x.fields.type === "district"
                ).length}
            )
          </Tag>
          <Tag
            color="cyan"
            className={`tag ${filter !== 3 ? "inactive" : ""}`}
            onClick={() => setFilter(3)}
          >
            Communes (
            {keyword === ""
              ? data.countCommunes
              : _.defaultTo(search.hits, []).filter(
                  (x) => x.fields.type === "commune"
                ).length}
            )
          </Tag>
          <Tag
            color="cyan"
            className={`tag ${filter !== 4 ? "inactive" : ""}`}
            onClick={() => setFilter(4)}
          >
            Fokontany (
            {keyword === ""
              ? data.countFokontany
              : _.defaultTo(search.hits, []).filter(
                  (x) => x.fields.type === "fokontany"
                ).length}
            )
          </Tag>
        </div>
      )}
      <PopoverList
        {...{
          filter,
          search,
          error,
          handleUpdate,
          loading,
          keyword,
          fokontany,
          communes,
          districts,
          regions,
        }}
      />
    </div>
  );
};

export default Popover;
