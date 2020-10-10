import React, { useState, useEffect } from 'react';
import axios from 'axios';


const columns = [
    { id: 'ID', label: 'ID#', minWidth: 170 },
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
                                {
                                    columns.map((column) => {
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
                                                columns.map((column) => {
                                                    const value = user[column.id];
                                                    return (
                                                        <td>{column.formatedCell ? column.format(value) : value}</td>
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

