import React, { useState, useEffect } from "react";
import DeckGL from "@deck.gl/react";
import { GeoJsonLayer } from "@deck.gl/layers";
import MapGL from "react-map-gl";
import Popover from "../../Components/Popover";
import { useQuery, gql } from "@apollo/client";
import MDSpinner from "react-md-spinner";
import { useParams } from "react-router-dom";

const DISTRICT = gql`
  query District($id: ID!) {
    district(id: $id) {
      id
      name
      coordinates
    }
  }
`;

// Set your mapbox access token here
const MAPBOX_ACCESS_TOKEN = process.env.REACT_APP_MAPBOX_ACCESS_TOKEN;
const MAPBOX_STYLE = process.env.REACT_APP_MAP_STYLE;

const District = (props) => {
  const { id } = useParams();
  const { loading, error, data } = useQuery(DISTRICT, { variables: { id } });
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
      const { coordinates, name } = data.district;
      const [longitude, latitude] = coordinates[0][0][0];
      const location = {
        ...viewport,
        longitude,
        latitude,
        zoom: 8,
      };
      setName(name);
      setViewport(location);
      const geojson = {
        type: "FeatureCollection",
        features: [
          {
            type: "Feature",
            geometry: {
              type: "MultiPolygon",
              coordinates: coordinates[0].map(x => [x]),
            },
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
          getFillColor: [47, 84, 235, 127],
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
            <div className="ant-popover-arrow"></div>
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
          onViewportChange={(value) => setViewport(value)}
        />
      </DeckGL>
      <Popover {...{ popoverClass, setExpanded }} />
    </div>
  );
};

export default District;
