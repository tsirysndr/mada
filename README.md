# Mada CLI

Mada is a CLI that facilitates the search for regions, districts, communes and fokontany.

<img src="./preview.png" />


## 🚚 Installation

You need to install libspatialite before continuing
```bash
# For Ubuntu
sudo apt install libsqlite3-mod-spatialite
```

```bash
go install -v github.com/tsirysndr/mada@latest
```

## 🚀 Usage

```bash
~> mada search miaramasoandro

4 matches, showing 1 through 4, took 237.941µs
    1. 14f4faf96bea4762d84c869c049ec219862e5d49f6fbbbb114007f96f47ec9a5 (4.626249)
        name
                Miaramasoandro
        fokontany
                Miaramasoandro
        country
                Madagascar
        district
                Faratsiho
        province
                Antananarivo
        region
                Vakinankaratra
        type
                fokontany
        commune
                Miandrarivo
    2. e4597e4987161de9fbe39c2b4a8f20734b547e8cf078365e38046ce86f682b3a (4.626249)
        name
                Miaramasoandro
        fokontany
                Miaramasoandro
        region
                Itasy
        type
                fokontany
        commune
                Tamponala
        country
                Madagascar
        district
                Soavinandriana
        province
                Antananarivo
    3. ca8efe0e1942ff0b097e1e5a210392f2f1c5f971d46c250fbb2cbd11d8417375 (4.626249)
        fokontany
                Miaramasoandro
        name
                Miaramasoandro
        commune
                Ambato
        country
                Madagascar
        district
                Ambohidratrimo
        province
                Antananarivo
        region
                Analamanga
        type
                fokontany
    4. 614df925b36deff00561292d26c9ac8de1c10511fbeceb60a619ebb89e627aaf (4.361670)
        fokontany
                Miaramasoandro
        name
                Miaramasoandro
        district
                Antsirabe I
        province
                Antananarivo
        region
                Vakinankaratra
        type
                fokontany
        commune
                Mahazoarivo Avarabohitra
        country
                Madagascar

```
