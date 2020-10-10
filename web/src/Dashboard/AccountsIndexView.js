import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { AccountIndexViewFields } from '../constants'

export default function StickyHeadTable(props) {
  const [loading, setLoading] = useState(true)
  const [users, setUsers] = useState([])
  useEffect(() => {
    const getAccounts = async () => {
      axios.get("http://localhost:3001/v1/accounts", { withCredentials: true }).then(response => {
        if (response.status === 200) {
          setUsers(response.data);
          setLoading(false);
        }
      });
    }

    if (props.isAuthed) {
      getAccounts()
    }
  }, [props])

  return loading ? (
    <div align="center">Loading...</div>
  ) : (
      <div className="container border col-10">
        <div className="row">
          <table class="table table-hover table-striped">
            <thead className="thead-dark">
              <tr>
                {
                  AccountIndexViewFields.map((column) => {
                    return (
                      <th scope="col">{column.label}</th>
                    );
                  })
                }
              </tr>
            </thead>
            <tbody>
              {
                users.map((user) => {
                  return (
                    <tr>
                      {
                        AccountIndexViewFields.map((column) => {
                          const value = user[column.id];
                          return (
                            column.id === 'FirstName' ? (
                              <td onclick="window.location='#';"><a className="d-block font-weight-bold" href={"/accounts/" + user.ID}>{column.formatedCell ? column.format(value) : value}</a></td>
                            ) : (
                                <td>{column.formatedCell ? column.format(value) : value}</td>
                              )
                          );
                        })
                      }
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

