import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { AccountFields } from '../constants'

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
                    <table class="table table-hover">
                        <thead>
                            <tr>
                                {
                                    AccountFields.map((column) => {
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
                                                AccountFields.map((column) => {
                                                    const value = user[column.id];
                                                    return (
                                                        <td><a className="d-block" href={"/accounts/" + user.ID}>{column.formatedCell ? column.format(value) : value}</a></td>
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

