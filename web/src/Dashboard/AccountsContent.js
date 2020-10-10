import React, { useState, useEffect } from 'react';
import axios from 'axios';


const columns = [
    { id: 'ID', label: 'ID', minWidth: 170 },
    { id: 'FirstName', label: 'First Name', minWidth: 170 },
    { id: 'LastName', label: 'Last Name', minWidth: 100 },
    { id: 'Email', label: 'Email', minWidth: 170, align: 'center' },
    { id: 'Active', label: 'Active', formatedCell: true, minWidth: 170, align: 'center', format: (value) => value ? "True" : "False" },
];

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
                                <th scope="col">ID#</th>
                                <th scope="col">First Name</th>
                                <th scope="col">Last Name</th>
                                <th scope="col">Email</th>
                                <th scope="col">Active</th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                users.map((user) => {
                                    return (
                                        <tr>
                                            <th scope="row">{user.ID}</th>
                                            <td>{user.FirstName}</td>
                                            <td>{user.LastName}</td>
                                            <td>{user.Email}</td>
                                            <td>{user.Active}</td>
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

