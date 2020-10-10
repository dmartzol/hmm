import React, { useState, useEffect } from 'react';
import axios from 'axios';


const columns = [
    { id: 'ID', label: 'ID#', minWidth: 170 },
    { id: 'FirstName', label: 'First Name', minWidth: 170 },
    { id: 'LastName', label: 'Last Name', minWidth: 100 },
    { id: 'Email', label: 'Email', minWidth: 170, align: 'center' },
    { id: 'Active', label: 'Active', formatedCell: true, minWidth: 170, align: 'center', format: (value) => value ? "True" : "False" },
];

export default function SingleAccount(props) {

    const [loading, setLoading] = useState(true)
    const [user, setUser] = useState({})
    useEffect(() => {
        const getMyAccount = async () => {
          axios.get(
            "http://localhost:3001/v1/accounts/" + props.match.params.id,
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

