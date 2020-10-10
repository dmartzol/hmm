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
        <div className="row">
          <div className="col-sm">
            <div className="title mt-4 border-bottom border-dark">
              <h1>{user.FirstName} {user.LastName}</h1>
            </div>
          </div>
        </div>
        <div className="row mt-3">
          <div className="col-lg mt-2">
            <div>
              <h3 className="border-bottom border-dark m-0">Account details</h3>
            </div>
            <table className="table table-hover">
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
          <div className="col-lg mt-2">
            <div>
              <h3 className="border-bottom border-dark m-0">Roles</h3>
            </div>
            <table className="table table-hover">
              <tbody>
                {
                  user.Roles ? (
                    user.Roles.map((role) => {
                      return (
                        <tr><td>{role.Name}</td></tr>
                      );
                    })
                  ) : (<div></div>)
                }
              </tbody>
            </table>
          </div>
        </div>

      </div>
    )
}

