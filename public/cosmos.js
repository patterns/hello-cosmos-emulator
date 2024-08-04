// obtain json data used by the data table
async function getList() {
    try {
        const resp = await fetch(
            "/api/users",
            {
                method: "GET",
            },
        );

        if (!resp.ok) {
            return new Response('Fetch list fail', { status: 400 });
        }
        const data = await resp.json();
        return data;
    } catch(err) {
        return new Response('List JSON fail', { status: 400 });
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

