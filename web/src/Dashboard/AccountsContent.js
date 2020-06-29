import React, { useState, useEffect } from 'react'
import axios from 'axios'

export default function UserPage(props) {
  const initialState = {
    user: {},
    loading: true,
  }

  const [state, setState] = useState(initialState)
  useEffect(() => {
    const getAccounts = async () => {
      axios.get(
        "http://localhost:3001/v1/accounts",
        { withCredentials: true}
      ).then(response => {
        if (response.status === 200) {
          console.log(response);
          setState({
            loading: false,
            user: response.data,
          });
        }
      });
    }

    getAccounts()
  }, [initialState])

  return initialState.loading ? (
    <div align="center">Loading...</div>
  ) : (
      <div className="container">
        <h1>{props.match.params.id}</h1>
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Location</th>
              <th>Website</th>
              <th>Followers</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>{state.user.name}</td>
              <td>{state.user.location}</td>
              <td>
                <a href={state.user.blog}>{state.user.blog}</a>
              </td>
              <td>{state.user.followers}</td>
            </tr>
          </tbody>
        </table>
      </div>
    )
}