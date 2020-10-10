import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { SingleAccountFields } from '../constants'

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
        <div className="row d-block">
          <div className="title my-4 border-bottom border-dark">
            <h1>{user.FirstName} {user.LastName}</h1>
          </div>
          <table className="table-sm table-hover">
            <tbody>
              {
                SingleAccountFields.map((accountField) => {
                  const value = user[accountField.id];
                  return (
                    <tr>
                      <td>{accountField.label}:</td>
                      <td>{accountField.formatedCell ? accountField.format(value) : value}</td>
                      <td><button className="btn btn-sm m-0" disabled={!accountField.Editable}>Edit</button></td>
                    </tr>
                  );
                })
              }
            </tbody>
          </table>
        </div>

      </div>
    )
}

