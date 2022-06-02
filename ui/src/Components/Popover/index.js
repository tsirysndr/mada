import React, { useState, useEffect } from 'react'
import { useQuery, gql } from '@apollo/client'
import { Tag } from 'antd'
import MDSpinner from 'react-md-spinner'
import PopoverList from './PopoverList'

const ALL_FOKONTANY = gql`
  query AllFokontany($after: ID, $size: Int) {
    allFokontany(after: $after, size: $size) {
      data {
        id
        name
        code
        province
        commune
        district
        region
      }
      after {
        id
      }
    }
  }
`

const COMMUNES = gql`
  query Communes($after: ID, $size: Int) {
      communes(after: $after, size: $size) {
        data {
        id
        name
        province
        code
        district
        region
      }
      after {
        id
      }
    } 
  }
`

const DISTRICTS = gql`
  query Districts($after: ID, $size: Int) {
    districts(after: $after, size: $size) {
      data {
        id
        name
        code
        region
      }
      after {
        id
      }
    }
  }
`

const REGIONS = gql`
  query Regions($after: ID, $size: Int) {
    regions(after: $after, size: $size) {
      data {
        id
        name
        code
        province
      }
      after {
        id
      }
    }
  }
`

const SEARCH = gql`
  query Search($keyword: String!) {
    search(keyword: $keyword) {
      regions {
        id
        name
        code
      }
      fokontany {
        id
        name
        commune
        district
        region
        code
      }
      districts {
        id
        name
        region
        code
      }
      communes {
        id
        name
        district
        region
        code
      }
    }
    countRegions
    countDistricts
    countCommunes
    countFokontany
  }
`

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
  const [filter, setFilter] = useState(1)
  const [keyword, setKeyword] = useState('')
  const [fokontany, setFokontany] = useState([])
  const [communes, setCommunes] = useState([])
  const [districts, setDistricts] = useState([])
  const [regions, setRegions] = useState([])
  const [fokontanyAfter, setFokontanyAfter] = useState(null)
  const [communeAfter, setCommuneAfter] = useState(null)
  const [districtAfter, setDistrictAfter] = useState(null)
  const { loading, error, data } = useQuery(SEARCH, { variables: { keyword } })
  // const countRes = useQuery(COUNT)
  const allFokontanyRes = useQuery(ALL_FOKONTANY, { variables: { after: fokontanyAfter, size: 100 } })
  const communesRes = useQuery(COMMUNES, { variables: { after: communeAfter, size: 100 } })
  const districtsRes = useQuery(DISTRICTS, { variables: { after: districtAfter, size: 10 } })
  const regionsRes = useQuery(REGIONS, { variables: { after: null, size: 10 } })

  const handleUpdate = () => {
    switch (filter) {
      case 1:
        break
      case 2:
        setDistrictAfter(districts[districts.length - 1].id)
        break
      case 3:
        setCommuneAfter(communes[communes.length - 1].id)
        break
      case 4:
        setFokontanyAfter(fokontany[fokontany.length - 1].id)
        break
      default:
        break
    }
  }

  useEffect(() => {
    if (allFokontanyRes.data) {
      const newFokontany = !fokontanyAfter ? allFokontanyRes.data.allFokontany.data : [...fokontany, ...allFokontanyRes.data.allFokontany.data]
      setFokontany(newFokontany)
    }
  }, [allFokontanyRes.data])

  useEffect(() => {
    if (communesRes.data) {
      const newCommunes = !communeAfter ? communesRes.data.communes.data : [...communes, ...communesRes.data.communes.data]
      setCommunes(newCommunes)
    }
  }, [communesRes.data])

  useEffect(() => {
    if (districtsRes.data) {
      const newDistricts = !districtAfter ? districtsRes.data.districts.data : [...districts, ...districtsRes.data.districts.data]
      setDistricts(newDistricts)
    }
  }, [districtsRes.data])

  useEffect(() => {
    if (regionsRes.data) {
      setRegions(regionsRes.data.regions.data)
    }
  }, [regionsRes.data])

  return (
    <div id='search-popover' className={props.popoverClass}>
      <div style={{ display: 'flex', height: 64 }} onClick={() => props.setExpanded(true)}>
        <input
          id='search'
          placeholder='Search ...'
          style={{ marginLeft: 50 }}
          autoComplete='off'
          onChange={evt => {
            if (evt.target.value.length > 2 || evt.target.value.length === 0) {
              setKeyword(evt.target.value)
            }
          }}
        />
        <svg
          xmlns='http://www.w3.org/2000/svg' width='20' height='20' viewBox='0 0 14 14' className='search-icon'
          style={{ position: 'absolute', top: 20, left: 15 }}
        >
          <g fill='none' fillRule='evenodd'>
            <path d='M-5-5v24h24V-5z' />
            <path fill='#5C5C6E' fillRule='nonzero' d='M12.003 13.657L8.96 10.616c-.016-.016-.027-.035-.042-.052a5.756 5.756 0 1 1 1.645-1.645c.017.015.036.026.052.042l3.041 3.042a1.17 1.17 0 0 1-1.654 1.654zM9.516 5.756a3.76 3.76 0 1 0-7.52 0 3.76 3.76 0 0 0 7.52 0z' />
          </g>
        </svg>
      </div>
      {
        (loading || error) && (
          <div id='loader'>
            <MDSpinner />
          </div>
        )
      }
      {
        !loading && !error && (
          <div style={{ padding: 10 }}>
            <Tag
              color='cyan'
              className={`tag ${filter !== 1 ? 'inactive' : ''}`}
              onClick={() => setFilter(1)}
            >
              Regions ({keyword === '' ? data.countRegions : data.search.regions.length})
            </Tag>
            <Tag
              color='cyan'
              className={`tag ${filter !== 2 ? 'inactive' : ''}`}
              onClick={() => setFilter(2)}
            >
              Districts ({keyword === '' ? data.countDistricts : data.search.districts.length})
            </Tag>
            <Tag
              color='cyan'
              className={`tag ${filter !== 3 ? 'inactive' : ''}`}
              onClick={() => setFilter(3)}
            >
              Communes ({keyword === '' ? data.countCommunes : data.search.communes.length})
            </Tag>
            <Tag
              color='cyan'
              className={`tag ${filter !== 4 ? 'inactive' : ''}`}
              onClick={() => setFilter(4)}
            >
              Fokontany ({keyword === '' ? data.countFokontany : data.search.fokontany.length})
            </Tag>
          </div>
        )
      }
      <PopoverList {...{
        filter,
        data,
        error,
        handleUpdate,
        loading,
        keyword,
        fokontany,
        communes,
        districts,
        regions
      }}
      />
    </div>
  )
}

export default Popover