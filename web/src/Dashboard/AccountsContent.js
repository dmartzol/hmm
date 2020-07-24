import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import axios from 'axios';

const useStyles = makeStyles({
    table: {
        minWidth: 650,
    },
});

export default function SimpleTable(props) {
    const classes = useStyles();

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
            <TableContainer component={Paper}>
                <Table className={classes.table} aria-label="simple table">
                    <TableHead>
                        <TableRow>
                            <TableCell align="center">First Name</TableCell>
                            <TableCell align="center">Last Name</TableCell>
                            <TableCell align="center">Email</TableCell>
                            <TableCell align="center">Active</TableCell>
                            <TableCell align="center">Role</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {users.map((user) => (
                            <TableRow key={user.FirstName}>
                                <TableCell align="center" component="th" scope="user">{user.FirstName}</TableCell>
                                <TableCell align="center">{user.LastName}</TableCell>
                                <TableCell align="center">{user.Email}</TableCell>
                                <TableCell align="center">{user.Active}</TableCell>
                                <TableCell align="center">{user.Roles[0] ? user.Roles[0].Name : "-"}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        )

}