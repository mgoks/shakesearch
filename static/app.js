const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
		// TODO Use a for-of loop or forEach function to process results and parse
		// JSON once marshalling to JSON is implemented in the back-end.
		for (let i = 0; i < results.length; i += 2) {
			const title = results[i];
			const text = results[i+1];
			// TODO Update table using (HTML) DOM API instead using innerHTML.
			const row = `
					<tr class="bg-white border-b dark:bg-gray-900 dark:border-gray-700">
						<td class="py-4 px-6">
							<h6 class="text-lg font-bold dark:text-blue">${title}</h6>
							<p>${text}</p>
						</td>
					</tr>`;
			rows.push(row);
		}
		if (rows.length == 0) {
			rows.push(`
					<tr class="bg-white border-b dark:bg-gray-900 dark:border-gray-700">
						<td class="py-4 px-6">
							<h6 class="text-lg font-bold dark:text-blue">
								Looks like Shakespeare did not use this word or phrase in any of
								his works.
							</h6>
						</td>
					</tr>`);
		}
    table.innerHTML = rows.join('\n');
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
