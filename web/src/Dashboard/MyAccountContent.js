import React, { useState, useEffect } from 'react'
import axios from 'axios'

export default function UserPage(props) {

  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState({})
  useEffect(() => {
    const getMyAccount = async () => {
      axios.get(
        "http://localhost:3001/v1/accounts/" + props.session.session.AccountID,
        { withCredentials: true }
      ).then(response => {
        if (response.status === 200) {
          setUser(
            response.data,
          );
          setLoading(false);
        }
      });
    }

    if (props.isAuthed) {
      getMyAccount()
    }
  }, [props])

  return loading ? (
    <div align="center">Loading...</div>
  ) : (
      <div className="container">
        <h1>{props.match.params.id}</h1>
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Last Name</th>
              <th>Date of birth</th>
              <th>Email</th>
              <th>Gender</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>{user.FirstName}</td>
              <td>{user.LastName}</td>
              <td>{user.DateOfBird}</td>
              <td><a href={user.Email}>{user.Email}</a></td>
              <td>{user.Gender}</td>
            </tr>
          </tbody>
        </table>
      </div>
    )
}