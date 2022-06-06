/*var Shp = require('shp');
//var shpJson = Shp.readFileSync('../mada/shp/mdg_admbnda_adm1_BNGRC_OCHA_20181031.shp');
// or
Shp.readFile('../mada/shp/mdg_admbnda_adm1_BNGRC_OCHA_20181031', function(error, data){
   	console.log(JSON.stringify(data.features.filter(x => x.properties['ADM1_EN'].includes('Analan')).map(x => x.geometry.coordinates[0][0].length)));
})
*/
const fs = require("fs");
const shapefileToGeojson = require("shapefile-to-geojson");

async function main() {
  const files = [
    {
      shp: "../mada/shp/mdg_admbnda_adm0_BNGRC_OCHA_20181031.shp",
      dbf: "../mada/shp/mdg_admbnda_adm0_BNGRC_OCHA_20181031.dbf",
      geojson: "../mada/geojson/mdg_admbnda_adm0_BNGRC_OCHA_20181031.json",
    },
    {
      shp: "../mada/shp/mdg_admbnda_adm1_BNGRC_OCHA_20181031.shp",
      dbf: "../mada/shp/mdg_admbnda_adm1_BNGRC_OCHA_20181031.dbf",
      geojson: "../mada/geojson/mdg_admbnda_adm1_BNGRC_OCHA_20181031.json",
    },
    {
      shp: "../mada/shp/mdg_admbnda_adm2_BNGRC_OCHA_20181031.shp",
      dbf: "../mada/shp/mdg_admbnda_adm2_BNGRC_OCHA_20181031.dbf",
      geojson: "../mada/geojson/mdg_admbnda_adm2_BNGRC_OCHA_20181031.json",
    },
    {
      shp: "../mada/shp/mdg_admbnda_adm3_BNGRC_OCHA_20181031.shp",
      dbf: "../mada/shp/mdg_admbnda_adm3_BNGRC_OCHA_20181031.dbf",
      geojson: "../mada/geojson/mdg_admbnda_adm3_BNGRC_OCHA_20181031.json",
    },
    {
      shp: "../mada/shp/mdg_admbnda_adm4_BNGRC_OCHA_20181031.shp",
      dbf: "../mada/shp/mdg_admbnda_adm4_BNGRC_OCHA_20181031.dbf",
      geojson: "../mada/geojson/mdg_admbnda_adm4_BNGRC_OCHA_20181031.json",
    },
  ];
  for (const file of files) {
    const geoJSON = await shapefileToGeojson.parseFiles(file.shp, file.dbf);
    const features = geoJSON.features.map((x) => ({
      ...x,
      geometry: {
        ...x.geometry,
        coordinates: [x.geometry.coordinates],
      },
    }));
    geoJSON.features = features;
    /*console.log(
    JSON.stringify(
      geoJSON.features
        .filter((x) => x.properties["ADM1_EN"].includes("Sofia"))
        .map((x) => {
          x.geometry.coordinates;
          if (x.geometry.coordinates[1].length === 1) {
            return {
              ...x,
              geometry: {
                ...x.geometry,
                coordinates: [x.geometry.coordinates[1][0]],
              },
            }
          }
        })
    )
  );*/
    fs.writeFileSync(file.geojson, JSON.stringify(geoJSON));
  }
}

main();
