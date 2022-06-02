import React, { useState, useEffect } from "react";
import DeckGL from "@deck.gl/react";
import { GeoJsonLayer } from "@deck.gl/layers";
import MapGL from "react-map-gl";
import Popover from "../../Components/Popover";
import { useQuery, gql } from "@apollo/client";
import MDSpinner from "react-md-spinner";
import { useParams } from "react-router-dom";

const REGION = gql`
  query Region($id: ID!) {
    region(id: $id) {
      id
      name
      code
      province
      geometry {
        type
        polygon {
          type
          coordinates
        }
        multipolygon {
          type
          coordinates
        }
      }
    }
  }
`;

// Set your mapbox access token here
const MAPBOX_ACCESS_TOKEN = process.env.REACT_APP_MAPBOX_ACCESS_TOKEN;
const MAPBOX_STYLE = process.env.REACT_APP_MAP_STYLE;

const Region = (props) => {
  const { id } = useParams();
  const { loading, error, data } = useQuery(REGION, { variables: { id } });
  const [showPopup, setShowPopup] = useState(false);
  const [name, setName] = useState("");
  const [popupX, setPopupX] = useState(0);
  const [popupY, setPopupY] = useState(0);
  const [layers, setLayers] = useState([]);
  const [expanded, setExpanded] = useState(false);
  const popoverClass = `popover ${expanded ? "expand" : "shrink"}`;
  const [viewport, setViewport] = useState({
    longitude: 47.52186,
    latitude: -18.91449,
    zoom: 11.97,
    bearing: 0,
    pitch: 30,
  });

  useEffect(() => {
    setShowPopup(false);
    setLayers([]);
  }, [id]);

  useEffect(() => {
    if (!loading && !error) {
      const { geometry, name } = data.region;
      const { type } = geometry;
      const [longitude, latitude] =
        type === "Polygon"
          ? geometry.polygon.coordinates[0][0]
          : geometry.multipolygon.coordinates[0][0][0];
      const location = {
        ...viewport,
        longitude,
        latitude,
        zoom: 6,
      };
      setName(name);
      setViewport(location);
      const geojson = {
        type: "FeatureCollection",
        features: [
          {
            type: "Feature",
            geometry:
              type === "Polygon" ? geometry.polygon : geometry.multipolygon,
          },
        ],
      };
      setLayers([
        new GeoJsonLayer({
          id: "geojson-layer",
          data: geojson,
          pickable: true,
          stroked: false,
          filled: true,
          extruded: false,
          lineWidthScale: 20,
          lineWidthMinPixels: 2,
          getElevation: 1,
          getFillColor: [82, 196, 26, 127],
          onHover: ({ x, y }) => {
            if (x > 0 && y > 0) {
              setPopupX(x - 40);
              setPopupY(y - 40);
              setShowPopup(true);
            }
          },
        }),
      ]);
    }
  }, [loading, error, data]);

  if (loading) {
    return (
      <div className="spinner">
        <MDSpinner />
      </div>
    );
  }

  return (
    <div>
      {showPopup && (
        <div
          className="ant-popover ant-popover-placement-top"
          style={{ position: "absolute", left: popupX, top: popupY }}
        >
          <div className="ant-popover-content">
            <div class="ant-popover-arrow"></div>
            <div className="ant-popover-inner" role="tooltip">
              <div>
                <div className="ant-popover-inner-content">
                  <div>{name}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
      <DeckGL
        initialViewState={viewport}
        controller
        layers={layers}
        onClick={() => setExpanded(false)}
      >
        <MapGL
          {...viewport}
          width="100vw"
          height="100vh"
          maxPitch={85}
          mapboxApiAccessToken={MAPBOX_ACCESS_TOKEN}
          mapStyle={MAPBOX_STYLE}
          // onViewportChange={value => setViewport(value)}
        ></MapGL>
      </DeckGL>
      <Popover {...{ popoverClass, setExpanded }} />
    </div>
  );
};

export default Region;
