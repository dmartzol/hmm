import React, { useState, useEffect } from 'react'
import axios from 'axios'

export default function UserPage(props) {

    const sessionInitialStates = {
      AccountID: 0,
      LastActivityTime: "",
  };
  const [session, setSession] = useState(sessionInitialStates);
  useEffect(() => {
    const getSession = async () => {
      axios.get(
        "http://localhost:3001/v1/sessions",
        { withCredentials: true }
      ).then(response => {
        if (response.status === 200) {
          setSession(response.data);
          props.history.push('/accounts/' + session.AccountID);
        }
      }).catch(error => {
        props.history.push('/');
        return
      });
    }
    getSession()
  },[props, session])

  return (
    <div align="center">Loading...</div>
    )
}