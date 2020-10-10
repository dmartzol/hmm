import React, { useState, useEffect } from 'react'
import axios from 'axios'

export default function UserPage(props) {

    const initialStates = {
    loggedIn: false,
    session: {
      AccountID: 0,
      LastActivityTime: "",
    },
  };
  const [session, setSession] = useState(initialStates);
  useEffect(() => {
    const getSession = async () => {
      axios.get(
        "http://localhost:3001/v1/sessions",
        { withCredentials: true }
      ).then(response => {
        if (response.status === 200) {
          setSession({
            loggedIn: true,
            session: response.data,
          });
          props.history.push('/accounts/' + response.data.AccountID);
        }
      }).catch(error => {
        props.history.push('/');
        return
      });
    }
    getSession()
  },[props])

  return (
    <div align="center">Loading...</div>
    )
}