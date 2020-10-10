

export const AccountIndexViewFields = [
    { id: 'ID', label: 'ID#', minWidth: 170 },
    { id: 'FirstName', label: 'First Name', minWidth: 170 },
    { id: 'LastName', label: 'Last Name', minWidth: 100 },
    { id: 'Email', label: 'Email', minWidth: 170, align: 'center' },
    { id: 'Active', label: 'Active', formatedCell: true, minWidth: 170, align: 'center', format: (value) => value ? "True" : "False" },
];

export const SingleAccountFields = [
    { id: 'ID', label: 'ID#', minWidth: 170, Editable: false },
    { id: 'FirstName', label: 'First Name', minWidth: 170, Editable: true },
    { id: 'LastName', label: 'Last Name', minWidth: 100, Editable: true },
    { id: 'Email', label: 'Email', minWidth: 170, align: 'center', Editable: false },
    { id: 'FailedLoginsCount', label: 'Failed Logins', minWidth: 170, align: 'center', Editable: false },
    { id: 'Active', label: 'Active', formatedCell: true, minWidth: 170, align: 'center', format: (value) => value ? "True" : "False", Editable: false },
];