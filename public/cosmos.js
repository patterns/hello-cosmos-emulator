// obtain json data used by the data table
async function getList() {
    try {
        const response = await fetch(
            "/api/users",
            {
                method: "GET",
            },
        );

        if (!response.ok) {
            throw new Error(`Fetch list fail, ${response.status}`);
        }
        const data = await response.json();
        return data;
    }
    catch(error) {
        console.log(error);
    }
}

// print html table with json list
getList().then(data => {
    let table = '<table style="border-collapse: collapse;">';
    table += `
        <thead>
          <tr>
            <th>Name</th>
            <th>Role</th>
            <th>Email</th>
            <th>Deactivate</th>
          </tr>
        </thead>
        <tbody>
    `;

    data.forEach(user => {
        table += `
          <tr>
            <td>${user.name}</td>
            <td>${user.role}</td>
            <td>${user.email}</td>
            <td>${user.deactivated}</td>
          </tr>
        `;
    });

    table += `
        </tbody>
      </table>
    `;

    const container = document.getElementById("container");
    container.innerHTML = table;
});

