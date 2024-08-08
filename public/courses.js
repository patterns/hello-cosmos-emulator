// obtain json data used by the data table
async function getCourseList() {
    try {
        const resp = await fetch(
            "/api/courses",
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
getCourseList().then(data => {
    let table = '<table style="border-collapse: collapse;">';
    table += `
        <thead>
          <tr>
            <th>Title</th>
            <th>Description</th>
            <th>URL</th>
            <th>Published</th>
          </tr>
        </thead>
        <tbody>
    `;

    data.forEach(course => {
        table += `
          <tr>
            <td>${course.title}</td>
            <td>${course.description}</td>
            <td>${course.url}</td>
            <td>${course.published}</td>
          </tr>
        `;
    });

    table += `
        </tbody>
      </table>
    `;

    const container = document.getElementById("course-list");
    container.innerHTML = table;
});

